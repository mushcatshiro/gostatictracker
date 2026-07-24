package dbop

import (
	"fmt"
	"testing"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/mock"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/stretchr/testify/assert"
)

func TestUpdateRecord(t *testing.T) {
	db := SetupTestTx(t)

	record := models.Record{
		Start:       common.MustParseStringDate("2023-10-01T10:00", ""),
		End:         common.MustParseStringDate("2023-10-01T11:00", ""),
		ActualEnd:   common.MustParseStringDate("2023-10-01T11:00", ""),
		Group:       "Group A",
		AllDay:      false,
		Title:       "record",
		URL:         "http://example.com/updated",
		Description: "Updated Description",
	}

	id, err := db.InsertRecord(record)
	assert.NoError(t, err)

	record.ID = id
	record.Title = "updated record"

	err = db.UpdateRecord(record)
	assert.NoError(t, err)
}

func TestUpdateRecordDoesNotExist(t *testing.T) {
	db := SetupTestTx(t)

	record := models.Record{
		ID:          999,
		Start:       common.MustParseStringDate("2023-10-01T10:00", ""),
		End:         common.MustParseStringDate("2023-10-01T11:00", ""),
		ActualEnd:   common.MustParseStringDate("2023-10-01T11:00", ""),
		Group:       "Group A",
		AllDay:      false,
		Title:       "Updated Event",
		URL:         "http://example.com/updated",
		Description: "Updated Description",
	}

	err := db.UpdateRecord(record)
	assert.Error(t, err)
}

func TestInsertNullTimeRecord(t *testing.T) {
	db := SetupTestTx(t)

	r := models.Record{
		Title: "test insert",
	}
	id, err := db.InsertRecord(r)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, id, fmt.Sprintf("expected non zero value, got: %d", id))

	// test nullable
	rr, err := db.ReadRecord(models.Record{ID: id})
	assert.NoError(t, err)
	t.Logf("%+v", rr)
	assert.Nil(t, rr.Start)
	assert.Nil(t, rr.End)
}

func TestReadRecord(t *testing.T) {
	db := SetupTestTx(t)

	_, err := db.ReadRecord(models.Record{ID: 1})
	assert.NoError(t, err)

	rs, err := db.ReadRecords(models.Record{Group: "day view example"})
	assert.NoError(t, err)
	assert.Equal(t, len(rs), len(mock.DayViewMockData))
}
