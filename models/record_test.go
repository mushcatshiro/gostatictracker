package models

import (
	"testing"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/stretchr/testify/assert"
)

func TestMarshalEventJSON(t *testing.T) {
	r := Record{
		ID:          1,
		Start:       common.ParseStringDate("01-02-2023 10:00", false, false),
		End:         common.ParseStringDate("01-02-2023 11:00", false, false),
		ActualEnd:   common.ParseStringDate("01-02-2023 11:00", false, false),
		Group:       "Group A",
		AllDay:      false,
		Title:       "Event A",
		URL:         "http://example.com/a",
		Description: "Description A",
	}

	data, err := r.MarshalJSON()
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
	assert.Contains(t, string(data), `"start":"01-02-2023 10:00"`)
	assert.Contains(t, string(data), `"title":"Event A"`)
}

func TestUnmarshalEventJSON(t *testing.T) {
	jsonData := `{
		"id": 1,
		"start": "01-02-2023 10:00",
		"end": "01-02-2023 11:00",
		"actualEnd": "01-02-2023 11:00",
		"group": "Group A",
		"allDay": false,
		"title": "Event A",
		"url": "http://example.com/a",
		"description": "Description A"
	}`

	var r Record
	err := r.UnmarshalJSON([]byte(jsonData))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), r.ID)
	assert.Equal(t, "01-02-2023 10:00", r.Start)
	assert.Equal(t, "01-02-2023 11:00", r.End)
	assert.Equal(t, "Group A", r.Group)
	assert.Equal(t, false, r.AllDay)
	assert.Equal(t, "Event A", r.Title)
	assert.Equal(t, "http://example.com/a", r.URL)
	assert.Equal(t, "Description A", r.Description)
}

func TestUnmarshalEventJSONInvalidStart(t *testing.T) {
	jsonData := `{
		"id": 1,
		"start": "invalid-date",
		"title": "Event A"
	}`

	var r Record
	err := r.UnmarshalJSON([]byte(jsonData))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid start format")
}
