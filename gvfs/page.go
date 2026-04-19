package gvfs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Page struct {
	Path         string
	FullURL      string
	SideRepoPath string
	Meta         PageMeta
}

type PageMeta struct {
	Title            string    `yaml:"title"`
	HasMermaid       bool      `yaml:"mermaid"`
	HasMathJax       bool      `yaml:"mathjax"`
	IsDraft          bool      `yaml:"draft"`
	IsPrivate        bool      `yaml:"private"`
	Tags             []string  `yaml:"tags"`
	LastModifiedDate time.Time `yaml:"lastmodified"`
	CreateDate       time.Time `yaml:"date"`
}

func extractFrontMatter(r io.Reader) (PageMeta, error) {
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
		if strings.TrimSpace(line) == "---" {
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
	if buffer.Len() > 0 {
		err := yaml.Unmarshal(buffer.Bytes(), &pageMeta)
		if err != nil {
			return PageMeta{}, fmt.Errorf("yaml unmarshal error: %w", err)
		}
	} else {
		return PageMeta{}, fmt.Errorf("no font matter found")
	}

	return pageMeta, nil
}
