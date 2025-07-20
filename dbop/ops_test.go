package dbop

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Define the event to update
	event := Event{
		ID:          1,
		Start:       "2023-10-01 10:00:00",
		End:         "2023-10-01 11:00:00",
		ActualEnd:   "2023-10-01 11:00:00",
		Group:       "Group A",
		AllDay:      false,
		Title:       "Updated Event",
		URL:         "http://example.com/updated",
		Description: "Updated Description",
	}

	// Expect the update query
	mock.ExpectExec(`
	UPDATE events SET start = \$1, end = \$2, actualEnd = \$3, "group" = \$4, allDay = \$5, title = \$6, url = \$7, description = \$8 WHERE id = \$9;
	`).
		WithArgs(event.Start, event.End, event.ActualEnd, event.Group, event.AllDay, event.Title, event.URL, event.Description, event.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = UpdateEvent(db, event)
	assert.NoError(t, err)
}

func TestUpdateEventDoesNotExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Define the event to update
	event := Event{
		ID:          999, // Non-existent ID
		Start:       "2023-10-01 10:00:00",
		End:         "2023-10-01 11:00:00",
		ActualEnd:   "2023-10-01 11:00:00",
		Group:       "Group A",
		AllDay:      false,
		Title:       "Updated Event",
		URL:         "http://example.com/updated",
		Description: "Updated Description",
	}

	// Expect the update query to return no rows affected
	mock.ExpectExec(`
	UPDATE events SET start = \$1, end = \$2, actualEnd = \$3, "group" = \$4, allDay = \$5, title = \$6, url = \$7, description = \$8 WHERE id = \$9;
	`).
		WithArgs(event.Start, event.End, event.ActualEnd, event.Group, event.AllDay, event.Title, event.URL, event.Description, event.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = UpdateEvent(db, event)
	assert.Error(t, err)
}
