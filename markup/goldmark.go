package markup

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gohugoio/hugo-goldmark-extensions/passthrough"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

var pageContextKey = parser.NewContextKey()

type goldmarkConverter struct {
	engine goldmark.Markdown
}

func NewGoldmarkConverter() *goldmarkConverter {
	gm := goldmark.New(
		goldmark.WithExtensions(
			// 1. Hugo's Passthrough for MathJax
			passthrough.New(passthrough.Config{
				InlineDelimiters: []passthrough.Delimiters{
					{ Open: "$", Close:"$"},
				},
				BlockDelimiters: []passthrough.Delimiters{
					{ Open: "$$", Close: "$$"},
				},
			}),
		),
		goldmark.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(NewHugoHTMLRenderer(), 100),
			),
		),
	)
	return &goldmarkConverter{engine: gm}
}

// HugoHTMLRenderer handles both standard code blocks (Chroma) and our custom readFile shortcode
type HugoHTMLRenderer struct{}

func NewHugoHTMLRenderer() renderer.NodeRenderer {
	return &HugoHTMLRenderer{}
}

func (r *HugoHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// Handle standard code fences
	reg.Register(ast.KindFencedCodeBlock, r.renderCodeBlock)
	// Handle our shortcode-like text
	reg.Register(ast.KindText, r.renderTextWithShortcodes)
}

// renderCodeBlock manually uses Chroma (Hugo style)
func (r *HugoHTMLRenderer) renderCodeBlock(
	w util.BufWriter, source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.FencedCodeBlock)
		lang := string(n.Language(source))

		var b strings.Builder
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			b.Write(line.Value(source))
		}

		// Use 'friendly' style from your hugo.toml
		htmlContent, _ := highlight(b.String(), lang, "friendly")
		w.WriteString(htmlContent)
	}
	return ast.WalkSkipChildren, nil
}

// renderTextWithShortcodes looks for {{< readFile ... >}} pattern
var readFileRegex = regexp.MustCompile(`\{\{<\s*readFile\s+file="([^"]+)"\s+lines="([^"]+)"\s+lang="([^"]+)"\s*>\}\}`)

func (r *HugoHTMLRenderer) renderTextWithShortcodes(
	w util.BufWriter, source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Text)
	segment := n.Segment
	content := string(segment.Value(source))

	match := readFileRegex.FindStringSubmatch(content)
	if match == nil {
		w.Write(segment.Value(source))
		return ast.WalkContinue, nil
	}

	// Extract params
	fileName := match[1]
	lineRange := match[2]
	lang := match[3]

	// Access Page from Context
	// Note: In a real implementation, you'd pass the context through RegisterFuncs
	// For brevity here, we assume the converter provides the path via the Page struct

	// Implementation of reading and highlighting
	snippet, err := readSnippet(fileName, lineRange) // Needs implementation logic
	if err != nil {
		w.WriteString(fmt.Sprintf("", err))
		return ast.WalkContinue, nil
	}

	highlighted, _ := highlight(snippet, lang, "friendly")
	w.WriteString(highlighted)

	return ast.WalkContinue, nil
}

// highlight is a helper to use Chroma manually
func highlight(code, lang, styleName string) (string, error) {
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	style := styles.Get(styleName)
	formatter := html.New(html.WithLineNumbers(true), html.TabWidth(4))

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	return buf.String(), err
}

func readSnippet(filePath string, rangeStr string) (string, error) {
	// Logic to open file and grab lines "3:22"
	// Use bufio.Scanner and a counter
	return "/* code snippet logic */", nil
}
