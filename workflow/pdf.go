package workflow

import (
	"errors"
	"io"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/mushcatshiro/gostatictracker/render"
)

// allow authenticated user to upload pdf - convert to html - update record
// to database (private)
func ConvertPdfToHtml(
	w io.Writer, auth bool, db dbop.DB, re render.RenderEngine,
) error {
	if !auth {
		return errors.New("Unauthorized")
	}
	return nil
}

func ListPapers(
	w io.Writer, auth bool, f models.Record, db dbop.DB, re render.RenderEngine,
) error {
	if !auth {
		return errors.New("Unauthorized")
	}
	return nil
}

func DisplayPaper() {}
