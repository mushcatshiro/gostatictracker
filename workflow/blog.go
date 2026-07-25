package workflow

import (
	"database/sql"
	"errors"
	"fmt"
	"io"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/mushcatshiro/gostatictracker/render"
)

func ListBlogPosts(
	w io.Writer, auth bool, db dbop.DB, re render.RenderEngine,
) error {
	f := models.Record{Group: models.GBlogpost}
	if auth {
		f.Status = common.NOTSTARTED
	} else {
		f.Status= common.COMPLETED
	}
	bs, err := db.ReadBlogRecords(f)
	if err != nil {
		return err
	}
	rm := models.RenderMeta{TemplateName: "bloglist.html", Data: bs}
	err = re.Render(w, rm)
	if err != nil {
		return err
	}
	return nil
}

func DisplayBlogPost(
	w io.Writer, auth bool, f models.Record, db dbop.DB, re render.RenderEngine,
) error {
	f.Group = models.GBlogpost
	if auth {
		f.Status = common.NOTSTARTED
	}
	b, err := db.ReadRecord(f)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no blog found with title: %s or id: %d", f.Title, f.ID)
		}
		return err
	}
	rm := models.RenderMeta{TemplateName: "blog.html", Data: b}
	err = re.Render(w, rm)
	if err != nil {
		return err
	}
	return nil
}

func EditBlogPost(
	w io.Writer, auth bool, f models.Record, db dbop.DB, re render.RenderEngine,
) error {
	if !auth {
		return errors.New("unauthorized")
	}
	// check if BlogManager has it
	// if not render empty editor
	rm := models.RenderMeta{TemplateName: "editor.html"}
	return re.Render(w, rm)
}

func PreviewBlogPost() {}

func CreateBlogPostRecord(
	w io.Writer, auth bool, r models.Record, db dbop.DB, re render.RenderEngine,
) error {
	// reads from form (web) and persist Record to db
	// render editor + Record (immutable)
	id, err := db.InsertRecord(r)
	if err != nil {
		return err
	}
	rm := models.RenderMeta{TemplateName: "editor.html", Data: id}
	return re.Render(w, rm)
}

func SaveBlogPost(
	w io.Writer, auth bool, r models.Record, db dbop.DB, re render.RenderEngine,
) (string, error) {
	return fmt.Sprintf("/blog/?id=%d", r.ID), nil
}
