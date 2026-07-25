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

func CreateNewRecord(w io.Writer, auth bool, r models.Record, db dbop.DB) {
	// check group in group table
	// if new group -> craete new group -> get gid
	// using gid -> create new Record
	// if insert fail
	// redirect to home + logout
	// if pid != 0
	// redirect to top level parent's group
	// else
	// redirect to current group
}

func ReadRecordsForSpecificGroup(
	w io.Writer, auth bool, groupName string, db dbop.DB, re render.RenderEngine,
) error {
	// if ggroup not found
	// redirect to home + logout
	// else
	// render
	return re.Render(w, models.RenderMeta{})
}

// support api & page
func readUniqueGroups() {
	// read groups and their corresponding records -> render to table + url
}

// support api version to help in PID
func readRecordsByMetadata() {}

func CopyFromRecord() {}

func CopyFromRecordRecursively() {}

func UpdateRecord() {}

func UpdatePriority() {}

// query by ID then compare the changes
func UpdateStatus() {}

func UpdateRelationship() {}

// handle the conversion of k=v,k=,
func UpdateMetadata() {}
