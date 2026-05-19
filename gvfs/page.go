package gvfs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type Page struct {
	Path         string
	FullURL      string
	SideRepoPath string
	Content      []byte
	Meta         PageMeta
}

type PageMeta struct {
	Title            string   `toml:"title"`
	HasMermaid       bool     `toml:"mermaid"`
	HasMathJax       bool     `toml:"math"`
	IsDraft          bool     `toml:"draft"`
	IsPrivate        bool     `toml:"private"`
	Tags             []string `toml:"tags"`
	LastModifiedDate string   `toml:"lastmodified"`
	CreateDate       string   `toml:"date"`
	ContentSidx      int
}

func ExtractFrontMatter(r io.Reader) (PageMeta, error) {
	scanner := bufio.NewScanner(r)
	var buffer bytes.Buffer
	var inFrontMatter bool
	var delimiterCount int
	var linecount int

	const MaxFrontMatterLines = 30

	// This is where a label might be used if logic was more nested
PARSE:
	for scanner.Scan() {
		line := scanner.Text()
		linecount++

		if linecount > MaxFrontMatterLines && inFrontMatter {
			return PageMeta{}, fmt.Errorf("front matter exceeded %d lines; missing closing '---'?", MaxFrontMatterLines)
		}

		// Check for YAML delimiter '---'
		if strings.TrimSpace(line) == "+++" {
			delimiterCount++

			if delimiterCount == 1 {
				inFrontMatter = true
				continue // Skip the first delimiter line
			}

			if delimiterCount == 2 {
				// We've reached the end of the front matter
				break PARSE
			}
		}

		if inFrontMatter {
			buffer.WriteString(line + "\n")
		} else if linecount > 1 && !inFrontMatter {
			// 2. Safety Check: If the first line wasn't '---',
			// this file doesn't have front matter. Stop immediately.
			break PARSE
		}
	}

	if err := scanner.Err(); err != nil {
		return PageMeta{}, err
	}

	if delimiterCount == 1 {
		return PageMeta{}, fmt.Errorf("found opening '---' but reached EOF without closing it")
	}

	// Now parse the collected YAML string into a map
	var pageMeta PageMeta
	// fmt.Printf("%s\n", buffer.String())
	if buffer.Len() > 0 {
		err := toml.Unmarshal(buffer.Bytes(), &pageMeta)
		if err != nil {
			return PageMeta{}, fmt.Errorf("toml unmarshal error: %v", err)
		}
	} else {
		return PageMeta{}, fmt.Errorf("no font matter found")
	}

	pageMeta.ContentSidx = linecount

	return pageMeta, nil
}
