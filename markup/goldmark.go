package markup

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

type RenderContext struct {
	FilePath string
}

type Converter interface {
	Convert(content []byte, ctx RenderContext) ([]byte, error)
}

type GoldmarkConverter struct {
	engine goldmark.Markdown
}

func NewGoldmarkConverter() *GoldmarkConverter {
	gm := goldmark.New(
		goldmark.WithExtensions(
			// 1. Hugo's Passthrough for MathJax
			passthrough.New(passthrough.Config{
				InlineDelimiters: []passthrough.Delimiters{
					{Open: "$", Close: "$"},
				},
				BlockDelimiters: []passthrough.Delimiters{
					{Open: "$$", Close: "$$"},
				},
			}),
		),
		goldmark.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(NewHugoHTMLRenderer(), 100),
			),
		),
	)
	return &GoldmarkConverter{engine: gm}
}

func (c *GoldmarkConverter) Convert(content []byte, ctx RenderContext) ([]byte, error) {
	var buf bytes.Buffer
	pc := parser.NewContext()
	pc.Set(pageContextKey, ctx)
	if err := c.engine.Convert(content, &buf); err != nil {
		return nil, fmt.Errorf("failed to convert markdown: %w", err)
	}
	return buf.Bytes(), nil
}

func Render(content []byte, ctx RenderContext) ([]byte, error) {
	c := NewGoldmarkConverter()
	return c.Convert(content, ctx)
}

// HugoHTMLRenderer handles both standard code blocks (Chroma)
// and our custom readFile shortcode
type HugoHTMLRenderer struct{}

func NewHugoHTMLRenderer() renderer.NodeRenderer {
	return &HugoHTMLRenderer{}
}

func (r *HugoHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindText, r.renderTextWithShortcodes)
}

func (r *HugoHTMLRenderer) renderCodeBlock(
	w util.BufWriter, source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkSkipChildren, nil
	}

	n := node.(*ast.FencedCodeBlock)
	lang := string(n.Language(source))

	if lang == "mermaid" {
		w.WriteString(`<pre><code class="language-mermaid">`)
		for i := 0; i < n.Lines().Len(); i++ {
			// FIX: Assign to a variable first to make it addressable
			line := n.Lines().At(i)
			w.Write(line.Value(source))
		}
		w.WriteString(`</code></pre>`)
		return ast.WalkSkipChildren, nil
	}

	var b strings.Builder
	for i := 0; i < n.Lines().Len(); i++ {
		line := n.Lines().At(i)
		b.Write(line.Value(source))
	}

	htmlContent, _ := highlight(b.String(), lang, "friendly")
	w.WriteString(htmlContent)
	return ast.WalkSkipChildren, nil
}

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

	fileName := match[1]
	lineRange := match[2]
	lang := match[3]

	// Access Page from Context
	// Note: In a real implementation, you'd pass the context through RegisterFuncs
	// For brevity here, we assume the converter provides the path via the Page struct

	// Implementation of reading and highlighting
	snippet, err := readSnippet(fileName, lineRange) // Needs implementation logic
	if err != nil {
		w.WriteString(fmt.Sprintf("<!-- error: %v -->", err))
		return ast.WalkContinue, nil
	}

	highlighted, _ := highlight(snippet, lang, "friendly")
	w.WriteString(highlighted)

	return ast.WalkContinue, nil
}

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
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	startLine := 1
	endLine := -1 // -1 acts as an indicator to read until the End Of File (EOF)

	// Parse the line range string
	parts := strings.Split(rangeStr, ":")
	if len(parts) > 0 && parts[0] != "" {
		if s, err := strconv.Atoi(parts[0]); err == nil {
			startLine = s
		}
	}
	if len(parts) > 1 {
		if parts[1] != "" {
			if e, err := strconv.Atoi(parts[1]); err == nil {
				endLine = e
			}
		}
	} else if len(parts) == 1 && parts[0] != "" {
		// If only a single number is provided without a colon (e.g., "10")
		endLine = startLine
	}

	var builder strings.Builder
	scanner := bufio.NewScanner(file)
	currentLine := 1

	for scanner.Scan() {
		// Start appending once we reach the start line
		if currentLine >= startLine {
			builder.WriteString(scanner.Text())
			builder.WriteString("\n")
		}

		// Stop reading early if we have hit the end line
		if endLine != -1 && currentLine >= endLine {
			break
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning file: %w", err)
	}

	// Remove the trailing newline for cleaner code block formatting
	return strings.TrimSuffix(builder.String(), "\n"), nil

}
