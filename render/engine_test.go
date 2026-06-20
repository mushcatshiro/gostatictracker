package render

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/mushcatshiro/gostatictracker/assets"
	"github.com/mushcatshiro/gostatictracker/mock"
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
	re, err := NewRenderEngine(assets.TemplateFS)
	if err != nil {
		t.Errorf("unable to create render engine: %v", err)
	}

	entries, err := assets.TemplateFS.ReadDir("templates")
	assert.NoError(t, err)
	totalTmpls := len(entries)
	actlTmpls := len(re.tmpls)
	assert.True(
		t,
		totalTmpls >= actlTmpls,
		fmt.Sprintf("should have %d html files, got %d", totalTmpls, actlTmpls),
	)

	goldenDir := filepath.Join(".testfiles", "golden")
	err = os.MkdirAll(goldenDir, 0755)
	assert.NoError(t, err, "failed to create golden directory")

	// doesnt include form
	for tmplName, meta := range mock.TemplateMetaMap {
		t.Logf("processing %s form", tmplName)
		rc := models.RenderMeta{
			TemplateName: tmplName,
			Data:         meta,
		}
		var buf bytes.Buffer
		err := re.Render(&buf, rc)
		assert.NoError(t, err)

		goldenPath := filepath.Join(goldenDir, tmplName[:len(tmplName)-5]+".golden.html")
		if *updateFlag {
			err = os.WriteFile(goldenPath, buf.Bytes(), 0644)
			assert.NoError(t, err, "failed to update golden file")
			continue
		}
		expected, err := os.ReadFile(goldenPath)
		assert.NoError(t, err, "golden file missing for %s. Run tests with -update flag to generate it. Error: %v", tmplName, err)

		if !bytes.Equal(expected, buf.Bytes()) {
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

func TestRenderForm(t *testing.T) {
	goldenDir := filepath.Join(".testfiles", "golden")
	err := os.MkdirAll(goldenDir, 0755)
	assert.NoError(t, err, "failed to create golden directory")
	re, err := NewRenderEngine(assets.TemplateFS)
	if err != nil {
		t.Errorf("unable to create render engine: %v", err)
	}
	for formType, meta := range mock.FormTemplateMetaMap {
		t.Logf("processing %s form", formType)
		rc := models.RenderMeta{
			TemplateName: "form.html",
			Data:         meta,
		}
		var buf bytes.Buffer
		err := re.Render(&buf, rc)
		assert.NoError(t, err)

		goldenPath := filepath.Join(goldenDir, formType+".golden.html")
		if *updateFlag {
			err = os.WriteFile(goldenPath, buf.Bytes(), 0644)
			assert.NoError(t, err, "failed to update golden file")
			continue
		}
		expected, err := os.ReadFile(goldenPath)
		assert.NoError(
			t, err, "golden file missing for %s. Run tests with -update flag to generate it. Error: %v", formType, err)

		if !bytes.Equal(expected, buf.Bytes()) {
			assert.Equal(t, string(expected), buf.String(), "render output does not match golden file for template: "+formType)
		}
	}

}
