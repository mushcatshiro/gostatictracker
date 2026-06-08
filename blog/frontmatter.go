package blog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/mushcatshiro/gostatictracker/models"
)

func ExtractFrontMatter(r io.Reader) (models.FrontMatter, error) {
	scanner := bufio.NewScanner(r)
	var buffer bytes.Buffer
	var inFrontMatter bool
	var delimiterCount int
	var linecount int
	var fm models.FrontMatter

	const MaxFrontMatterLines = 30

	// This is where a label might be used if logic was more nested
PARSE:
	for scanner.Scan() {
		line := scanner.Text()
		linecount++

		if linecount > MaxFrontMatterLines && inFrontMatter {
			return fm, fmt.Errorf("front matter exceeded %d lines; missing closing '---'?", MaxFrontMatterLines)
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
		return fm, err
	}

	if delimiterCount == 1 {
		return fm, fmt.Errorf("found opening '---' but reached EOF without closing it")
	}

	// Now parse the collected YAML string into a map
	// fmt.Printf("%s\n", buffer.String())
	if buffer.Len() > 0 {
		err := toml.Unmarshal(buffer.Bytes(), &fm)
		if err != nil {
			return fm, fmt.Errorf("toml unmarshal error: %v", err)
		}
	} else {
		return fm, fmt.Errorf("no font matter found")
	}

	fm.ContentSidx = linecount

	return fm, nil
}
