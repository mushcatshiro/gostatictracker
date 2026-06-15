package dbop

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/stretchr/testify/assert"
)

func TestSqlMockUpdateEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Define the event to update
	event := models.Event{
		ID:          1,
		Start:       common.ParseStringDate("10-01-2023 10:00", false, false),
		End:         common.ParseStringDate("10-01-2023 11:00", false, false),
		ActualEnd:   common.ParseStringDate("10-01-2023 11:00", false, false),
		Group:       "Group A",
		AllDay:      false,
		Title:       "Updated Event",
		URL:         "http://example.com/updated",
		Description: "Updated Description",
	}

	// Expect the update query
	mock.ExpectExec(`UPDATE events
	SET start = \$1, "end" = \$2, actualStart = \$3, actualEnd = \$4, insertTime = \$5,
	"group" = \$6, allDay = \$7, title = \$8, url = \$9, description = \$10, pid = \$11,
	priority = \$12, metadata = \$13, status = \$14
	WHERE id = \$15;`).
		WithArgs(
			event.Start,
			event.End,
			event.ActualStart,
			event.ActualEnd,
			event.InsertTime,
			event.Group,
			event.AllDay,
			event.Title,
			event.URL,
			event.Description,
			event.PID,
			event.Priority,
			event.Metadata,
			event.Status,
			event.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	err = UpdateEvent(db, event)
	assert.NoError(t, err)
}

func TestSqlMockUpdateEventDoesNotExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Define the event to update
	event := models.Event{
		ID:          999, // Non-existent ID
		Start:       common.ParseStringDate("10-01-2023 10:00", false, false),
		End:         common.ParseStringDate("10-01-2023 11:00", false, false),
		ActualEnd:   common.ParseStringDate("10-01-2023 11:00", false, false),
		Group:       "Group A",
		AllDay:      false,
		Title:       "Updated Event",
		URL:         "http://example.com/updated",
		Description: "Updated Description",
	}

	// Expect the update query to return no rows affected
	mock.ExpectExec(`UPDATE events
	SET start = \$1, "end" = \$2, actualStart = \$3, actualEnd = \$4, insertTime = \$5,
	"group" = \$6, allDay = \$7, title = \$8, url = \$9, description = \$10, pid = \$11,
	priority = \$12, metadata = \$13, status = \$14
	WHERE id = \$15;`).
		WithArgs(event.Start, event.End, event.ActualEnd, event.Group, event.AllDay, event.Title, event.URL, event.Description, event.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = UpdateEvent(db, event)
	assert.Error(t, err)
}

func TestEventOperations(t *testing.T) {
	tx := SetupTestTx(t)

	e := models.Event{
		Title: "test insert",
	}
	id, err := InsertEvent(tx, e)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, id, fmt.Sprintf("expected non zero value, got: %d", id))

	e.ID = id
	e.Title = "test update"
	err = UpdateEvent(tx, e)
	assert.NoError(t, err)

	ue, err := ReadEventById(tx, id)
	assert.NoError(t, err)
	assert.Equal(
		t,
		e.Title,
		ue.Title,
		fmt.Sprintf("expected %s, got: %s", e.Title, ue.Title),
	)

	dneE := models.Event{
		ID: 9999,
		Title: "DNE",
	}
	err = UpdateEvent(tx, dneE)
	assert.ErrorContains(t, err, "no event found with ID")

}
