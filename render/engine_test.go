package render

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/mushcatshiro/gostatictracker/assets"
	"github.com/mushcatshiro/gostatictracker/models"
)

var updateFlag = flag.Bool("update", false, "update golden files")

// supports two mode, manual approval and auto review
// trigger manual approval with `go test ./render -update`,
// *.golden.html will be created
// auto review works with no update flag. it compares the
// the generated to *.golden.html
func TestFullRenderEngine(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test render engine")
	}
	// read testdata/output from os Environ
	templates, err := template.ParseFS(assets.TemplateFS, "templates/*.html")
	assert.NoError(t, err, "failed to parse templates")

	re := NewRenderEngine(templates)

	entries, err := assets.TemplateFS.ReadDir("templates")
	assert.NoError(t, err)
	totalTmpls := len(entries)
	actlTmpls := len(templates.Templates())
	assert.True(
		t,
		totalTmpls >= actlTmpls,
		fmt.Sprintf("should have %d html files, got %d", totalTmpls, actlTmpls),
	)

	goldenDir := filepath.Join("testdata", "golden")
	err = os.MkdirAll(goldenDir, 0755)
	assert.NoError(t, err, "failed to create golden directory")

	testOutDir := filepath.Join("testdata", "output")
	err = os.MkdirAll(testOutDir, 0755)
	assert.NoError(t, err)

	for tmplName, meta := range models.TemplateMetaMap {
		rc := models.RenderMeta{
			TemplateName: tmplName,
			Data:         meta,
		}
		var buf bytes.Buffer
		err := re.Render(&buf, rc)
		assert.NoError(t, err)

		goldenPath := filepath.Join(goldenDir, tmplName+".golden.html")
		if *updateFlag {
			err = os.WriteFile(goldenPath, buf.Bytes(), 0644)
			assert.NoError(t, err, "failed to update golden file")
			continue
		}
		expected, err := os.ReadFile(goldenPath)
		assert.NoError(t, err, "golden file missing for %s. Run tests with -update flag to generate it. Error: %v", tmplName, err)
		/*
		outPath := filepath.Join(testOutDir, tmplName+".html")
		err = os.WriteFile(outPath, buf.Bytes(), 0644)
		assert.NoError(t, err, "failed to write visual check file")
		*/
		if !bytes.Equal(expected, buf.Bytes()) {
			// Using testify's assert.Equal on strings provides a beautiful, readable diff in the terminal when they mismatch
			assert.Equal(t, string(expected), buf.String(), "render output does not match golden file for template: "+tmplName)
		}
	}

	// template do not exists
	dneRc := models.RenderMeta{
		TemplateName: "dne",
	}
	var buf bytes.Buffer
	err = re.Render(&buf, dneRc)
	assert.EqualError(t, err, "template dne is not found")
}
