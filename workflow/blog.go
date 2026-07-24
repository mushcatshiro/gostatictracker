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
	f := models.Record{Group: "blog"}
	if auth {
		f.Status = common.NOTSTARTED
	}
	bs, err := db.ReadBlogRecords(f)
	if err != nil {
		return err
	}
	rm := models.RenderMeta{TemplateName: "", Data: bs}
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
	rm := models.RenderMeta{TemplateName: "", Data: b}
	err = re.Render(w, rm)
	if err != nil {
		return err
	}
	return nil
}

func EditBlogPost() {}

func PreviewBlogPost() {}

func CreateBlogPost() {}

func SaveBlogPost(
	w io.Writer, auth bool, r models.Record, db dbop.DB, re render.RenderEngine,
) (string, error) {
	var redirectTo string
	id, err := db.InsertRecord(r)
	if err != nil {
		return redirectTo, err
	}
	return fmt.Sprintf("/blog/?id=%d", id), nil
}
