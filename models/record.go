package models

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mushcatshiro/gostatictracker/common"
)

const GBookmarklet = "bookmarklet"
const GQuick = "quick"
const GTodo = "todo"
const GImage = "image"
const GReminder = "Reminder"
const GBlogpost = "blogpost"

var ReservedGroupName = []string{GBookmarklet, GQuick, GTodo, GImage, GReminder, GBlogpost}

type Record struct {
	ID          int64           `json:"id"`
	Start       *time.Time      `json:"start"`
	End         *time.Time      `json:"end"`
	ActualStart *time.Time      `json:"actualStart"`
	ActualEnd   *time.Time      `json:"actualEnd"`
	InsertTime  *time.Time      `json:"insertTime"`
	Group       string          `json:"group"`
	DefaultMode string          `json:"mode"`   // default render mode, todo etc
	Repeat      string          `json:"repeat"` // cron like syntax
	AllDay      bool            `json:"allDay"`
	Title       string          `json:"title"`
	URL         string          `json:"url"`
	Description string          `json:"description"`
	PID         int64           `json:"pid"`
	Priority    common.Priority `json:"priority"`
	Metadata    string          `json:"metadata"`
	Status      common.Status   `json:"status"`
}

func RecordToDataMap(r Record) map[string]string {
	retMap := make(map[string]string)
	retMap["id"] = strconv.FormatInt(r.ID, 10)
	if r.Start != nil {
		retMap["start"] = r.Start.Format(common.TimeLayout)
	}
	if r.End != nil {
		retMap["end"] = r.End.Format(common.TimeLayout)
	}
	if r.ActualStart != nil {
		retMap["actualstart"] = r.ActualStart.Format(common.TimeLayout)
	}
	if r.ActualEnd != nil {
		retMap["actualend"] = r.ActualEnd.Format(common.TimeLayout)
	}
	retMap["end"] = r.InsertTime.Format(common.TimeLayout)
	retMap["group"] = r.Group
	retMap["title"] = r.Title
	retMap["url"] = r.URL
	retMap["description"] = r.Description
	retMap["priority"] = strconv.FormatInt(int64(r.Priority), 10)
	retMap["metadata"] = r.Metadata
	retMap["status"] = strconv.FormatInt(int64(r.Status), 10)
	return retMap
}

// marhsall Event to JSON
func (r *Record) MarshalJSON() ([]byte, error) {
	type Alias Record// Create an alias to avoid recursion
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON unmarshals JSON data into an Event struct
func (r *Record) UnmarshalJSON(data []byte) error {
	type Alias Record// Create an alias to avoid recursion
	aux := &struct {
		Start       string `json:"start"`
		End         string `json:"end"`
		ActualStart string `json:"actualStart"`
		ActualEnd   string `json:"actualEnd"`
		InsertTime  string `json:"insertTime"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Start != "" {
		parsedStart, err := time.Parse(common.TimeLayout, aux.Start)
		if err != nil {
			return fmt.Errorf("invalid start format: %q (expected MM-DD-YYYY HH:MM): %w", aux.Start, err)
		}
		r.Start = &parsedStart
	} else {
		r.Start = nil
	}

	if aux.End != "" {
		parsedEnd, err := time.Parse(common.TimeLayout, aux.End)
		if err != nil {
			return fmt.Errorf("invalid end format: %q (expected MM-DD-YYYY HH:MM): %w", aux.End, err)
		}
		r.End = &parsedEnd
	} else {
		r.Start = nil
	}

	if aux.ActualStart != "" {
		parsedActualStart, err := time.Parse(common.TimeLayout, aux.ActualStart)
		if err != nil {
			return fmt.Errorf("invalid actualStart format: %q (expected MM-DD-YYYY HH:MM): %w", aux.ActualStart, err)
		}
		r.ActualStart = &parsedActualStart
	} else {
		r.Start = nil
	}

	if aux.ActualEnd != "" {
		parsedActualEnd, err := time.Parse(common.TimeLayout, aux.ActualEnd)
		if err != nil {
			return fmt.Errorf("invalid actualEnd format: %q (expected MM-DD-YYYY HH:MM): %w", aux.ActualEnd, err)
		}
		r.ActualEnd = &parsedActualEnd
	} else {
		r.Start = nil
	}

	now := time.Now()
	r.InsertTime = &now

	return nil
}

func (r *Record) Validate() bool {
	// validate repeat format
	// validate metadata format (tags)
	r.Group = strings.ToLower(r.Group)
	if slices.Contains(ReservedGroupName, r.Group) {
		return false
	}
	return true
}
