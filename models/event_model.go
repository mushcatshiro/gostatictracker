package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mushcatshiro/gostatictracker/common"
)

// struct event

type Event struct {
	ID          int64      `json:"id"`
	Start       *time.Time `json:"start"`
	End         *time.Time `json:"end"`
	ActualStart *time.Time `json:"actualStart"`
	ActualEnd   *time.Time `json:"actualEnd"`
	InsertTime  *time.Time `json:"insertTime"`
	Group       string     `json:"group"`
	AllDay      bool       `json:"allDay"`
	Title       string     `json:"title"`
	URL         string     `json:"url"`
	Description string     `json:"description"`
	PID         int64      `json:"pid"`
	Priority    int8       `json:"priority"`
	Metadata    string     `json:"metadata"`
	Status      int8       `json:"status"`
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
