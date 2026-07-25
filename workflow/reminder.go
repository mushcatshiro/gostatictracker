package workflow

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/mushcatshiro/gostatictracker/render"
	"github.com/mushcatshiro/gostatictracker/reminder"
)

func CreateNewReminder(
	w io.Writer, r models.Record, db dbop.DB, re render.RenderEngine,
) error {
	sched, err := reminder.ParseInput(r.Repeat)
	if err != nil {
		return err
	}
	r.Start = sched.Next()
	r.Status = common.INPROGRESS
	_, err = db.InsertRecord(r)
	if err != nil {
		return err
	}
	return ReadRecentReminders(w, 14, db, re)
}

// recent is an input range, the default is End < + 14, status = INPROGRESS
func ReadRecentReminders(
	w io.Writer, ahead int, db dbop.DB, re render.RenderEngine,
) error {
	// by looking at `Record.Start`
	if ahead == 0 {
		ahead = 14
	}
	today := time.Now()
	// format
	start := today.AddDate(0, 0, ahead)
	rs, err := db.ReadRecords(models.Record{
		Group: models.GReminder, Start: &start, Status: common.INPROGRESS,
	})
	if err != nil {
		return err
	}
	// render
	rm := models.RenderMeta{TemplateName: "list-html", Data: rs}
	return re.Render(w, rm)
}

// only update the repeat field
func UpdateReminderSettings(r models.Record, db dbop.DB) error {
	_, err := reminder.ParseInput(r.Repeat)
	if err != nil {
		return err
	}
	return db.UpdateRecord(r)
}

// only moves one step ahead. e.g. weekly reminder is missed twice
// user had to clear it twice; this is intentionally designed as such
// add to metadata missed=,ontime
func CompleteCurrentTrigger(id int64, db dbop.DB) error {
	r, err := db.ReadRecord(models.Record{ID: id})
	if err != nil {
		return err
	}
	sched, err := reminder.ParseInput(r.Repeat)
	if err != nil {
		return err
	}
	r.Start = sched.Next()
	return db.UpdateRecord(r)  // TODO: create specific update sql
}

func TerminateTrigger(id int64, db dbop.DB) error {
	r, err := db.ReadRecord(models.Record{ID: id})
	if err != nil {
		return err
	}
	r.Status = common.COMPLETED
	return db.UpdateRecord(r)  // TODO: create specific update sql
}
