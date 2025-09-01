package dbop

import (
	"database/sql"
	"fmt"
	"log"
)

type GanttEvent struct {
	ID          int64  `json:"id"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Group       string `json:"group"`
	AllDay      bool   `json:"allDay"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

func GetGanttGroupEvents(db *sql.DB, groupName string, dateOnly bool) ([]GanttEvent, error) {
	// TODO consider adding ordering by "end" DESC
	var events []GanttEvent
	var datetimeFormat string
	if dateOnly {
		datetimeFormat = "MM-DD-YYYY"
	} else {
		datetimeFormat = "MM-DD-YYYY HH24:MI"
	}
	query := fmt.Sprintf(`SELECT
	  "id",
		TO_CHAR(start, '%s'),
		TO_CHAR("end", '%s'),
		"group",
		allDay,
		title,
		url,
		description
	FROM events
	WHERE "group" = $1
	ORDER BY "id";`, datetimeFormat, datetimeFormat)
	rows, err := db.Query(query, groupName)
	if err != nil {
		return events, fmt.Errorf("failed to get events: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var e GanttEvent
		if err := rows.Scan(&e.ID, &e.Start, &e.End, &e.Group, &e.AllDay, &e.Title, &e.URL, &e.Description); err != nil {
			return events, fmt.Errorf("failed to scan event: %v", err)
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return events, fmt.Errorf("error reading events: %v", err)
	}
	return events, nil
}

func GetUniqueGroups(db *sql.DB) ([]string, error) {
	var groups []string
	query := `SELECT DISTINCT "group" FROM events;`
	rows, err := db.Query(query)
	if err != nil {
		// log.Fatalf("Failed to get unique groups: %v", err)
		return groups, err
	}
	defer rows.Close()

	for rows.Next() {
		var group string
		if err := rows.Scan(&group); err != nil {
			// log.Fatalf("Failed to scan group: %v", err)
			return groups, err
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Error reading groups: %v", err)
	}
	return groups, nil
}
