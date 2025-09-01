package dbop

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const TimeLayout = "01-02-2006 15:04" // MM-DD-YYYY hh:mm

const (
	STARTED = iota
	INPROGRESS
	COMPLETED
)

const (
	DONOW = iota
	DOLATER
	DELEGATE
	ELIMINATE
)

// struct event

type Event struct {
	ID          int64  `json:"id"`
	Start       string `json:"start"`
	End         string `json:"end"`
	ActualStart string `json:"actualStart"`
	ActualEnd   string `json:"actualEnd"`
	InsertTime  string `json:"insertTime"`
	Group       string `json:"group"`
	AllDay      bool   `json:"allDay"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	PID         int64  `json:"pid"`
	Priority    int8   `json:"priority"`
	Metadata    string `json:"metadata"`
	Status      int8   `json:"status"`
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
		End         string `json:"end,omitempty"`
		ActualStart string `json:"actualStart,omitempty"`
		ActualEnd   string `json:"actualEnd,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Start != "" {
		if _, err := time.Parse(TimeLayout, aux.Start); err != nil {
			return fmt.Errorf("invalid start format: %q (expected MM-DD-YYYY HH:MM): %w", aux.Start, err)
		}
		e.Start = aux.Start
	}

	if aux.End != "" {
		if _, err := time.Parse(TimeLayout, aux.End); err != nil {
			return fmt.Errorf("invalid end format: %q (expected MM-DD-YYYY HH:MM): %w", aux.End, err)
		}
		e.End = aux.End
	}

	if aux.ActualStart != "" {
		if _, err := time.Parse(TimeLayout, aux.ActualStart); err != nil {
			return fmt.Errorf("invalid actualStart format: %q (expected MM-DD-YYYY HH:MM): %w", aux.ActualStart, err)
		}
		e.ActualStart = aux.ActualStart
	}

	if aux.ActualEnd != "" {
		if _, err := time.Parse(TimeLayout, aux.ActualEnd); err != nil {
			return fmt.Errorf("invalid actualEnd format: %q (expected MM-DD-YYYY HH:MM): %w", aux.ActualEnd, err)
		}
		e.ActualEnd = aux.ActualEnd
	}

	return nil
}

// marshall list of Events to JSON
func MarshalEvents(events []Event) ([]byte, error) {
	type Alias Event // Create an alias to avoid recursion
	aliasEvents := make([]*struct {
		Start       string `json:"start"`
		End         string `json:"end,omitempty"`
		ActualStart string `json:"actualStart,omitempty"`
		ActualEnd   string `json:"actualEnd,omitempty"`
		*Alias
	}, len(events))

	for i, e := range events {
		aliasEvents[i] = &struct {
			Start       string `json:"start"`
			End         string `json:"end,omitempty"`
			ActualStart string `json:"actualStart,omitempty"`
			ActualEnd   string `json:"actualEnd,omitempty"`
			*Alias
		}{
			Start:       e.Start,
			End:         e.End,
			ActualStart: e.ActualStart,
			ActualEnd:   e.ActualEnd,
			Alias:       (*Alias)(&e),
		}
	}

	return json.Marshal(aliasEvents)
}

func UnmarshalEvents(data []byte) ([]Event, error) {
	type Alias Event // Create an alias to avoid recursion
	var aliasEvents []struct {
		Start       string `json:"start"`
		End         string `json:"end,omitempty"`
		ActualStart string `json:"actualStart,omitempty"`
		ActualEnd   string `json:"actualEnd,omitempty"`
		*Alias
	}

	if err := json.Unmarshal(data, &aliasEvents); err != nil {
		return nil, err
	}

	events := make([]Event, len(aliasEvents))
	for i, ae := range aliasEvents {
		e := Event(*ae.Alias)
		e.Start = ae.Start
		e.End = ae.End
		e.ActualStart = ae.ActualStart
		e.ActualEnd = ae.ActualEnd
		events[i] = e
	}

	return events, nil
}

func ConnectDB(username string, password string, dbhost, dbname string) (*sql.DB, error) {
	if username == "" || dbname == "" || password == "" {
		return nil, fmt.Errorf("username, password, dbname must be provided")
	}
	if dbhost == "" {
		dbhost = "localhost" // Default to localhost if no host is provided
	}

	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", username, password, dbhost, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Check if the connection is established
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
