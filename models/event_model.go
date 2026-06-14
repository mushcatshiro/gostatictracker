package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/mushcatshiro/gostatictracker/common"
)

type Event struct {
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

func (e *Event) ToGanttEvent() GanttEvent {
	ge := GanttEvent{
		ID:          e.ID,
		Start:       e.Start,
		End:         e.End,
		Group:       e.Group,
		AllDay:      e.AllDay,
		Title:       e.Title,
		URL:         e.URL,
		Description: e.Description,
	}
	return ge
}

func (e *Event) ToBookmarklet() Bookmarklet {
	return Bookmarklet{
		Title:       e.Title,
		Description: e.Description,
		URL:         e.URL,
		InsertTime:  e.InsertTime,
		Group:       e.Group,
	}
}

func (e *Event) ToDataMap() map[string]string {
	retMap := make(map[string]string)
	retMap["id"] = strconv.FormatInt(e.ID, 10)
	if e.Start != nil {
		retMap["start"] = e.Start.Format(common.TimeLayout)
	}
	if e.End != nil {
		retMap["end"] = e.End.Format(common.TimeLayout)
	}
	if e.ActualStart != nil {
		retMap["actualstart"] = e.ActualStart.Format(common.TimeLayout)
	}
	if e.ActualEnd != nil {
		retMap["actualend"] = e.ActualEnd.Format(common.TimeLayout)
	}
	retMap["end"] = e.InsertTime.Format(common.TimeLayout)
	retMap["group"] = e.Group
	retMap["title"] = e.Title
	retMap["url"] = e.URL
	retMap["description"] = e.Description
	retMap["priority"] = strconv.FormatInt(int64(e.Priority), 10)
	retMap["metadata"] = e.Metadata
	retMap["status"] = strconv.FormatInt(int64(e.Status), 10)
	return retMap
}

// marhsall Event to JSON
func (e *Event) MarshalJSON() ([]byte, error) {
	type Alias Event // Create an alias to avoid recursion
	return json.Marshal((*Alias)(e))
}

// UnmarshalJSON unmarshals JSON data into an Event struct
func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event // Create an alias to avoid recursion
	aux := &struct {
		Start       string `json:"start"`
		End         string `json:"end"`
		ActualStart string `json:"actualStart"`
		ActualEnd   string `json:"actualEnd"`
		InsertTime  string `json:"insertTime"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Start != "" {
		parsedStart, err := time.Parse(common.TimeLayout, aux.Start)
		if err != nil {
			return fmt.Errorf("invalid start format: %q (expected MM-DD-YYYY HH:MM): %w", aux.Start, err)
		}
		e.Start = &parsedStart
	} else {
		e.Start = nil
	}

	if aux.End != "" {
		parsedEnd, err := time.Parse(common.TimeLayout, aux.End)
		if err != nil {
			return fmt.Errorf("invalid end format: %q (expected MM-DD-YYYY HH:MM): %w", aux.End, err)
		}
		e.End = &parsedEnd
	} else {
		e.Start = nil
	}

	if aux.ActualStart != "" {
		parsedActualStart, err := time.Parse(common.TimeLayout, aux.ActualStart)
		if err != nil {
			return fmt.Errorf("invalid actualStart format: %q (expected MM-DD-YYYY HH:MM): %w", aux.ActualStart, err)
		}
		e.ActualStart = &parsedActualStart
	} else {
		e.Start = nil
	}

	if aux.ActualEnd != "" {
		parsedActualEnd, err := time.Parse(common.TimeLayout, aux.ActualEnd)
		if err != nil {
			return fmt.Errorf("invalid actualEnd format: %q (expected MM-DD-YYYY HH:MM): %w", aux.ActualEnd, err)
		}
		e.ActualEnd = &parsedActualEnd
	} else {
		e.Start = nil
	}

	now := time.Now()
	e.InsertTime = &now

	return nil
}

func (e *Event) Validate() error {
	// validate repeat format
	// validate metadata format (tags)
	return nil
}
