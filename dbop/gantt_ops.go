package dbop

import (
	"database/sql"
	"fmt"

	"github.com/mushcatshiro/gostatictracker/models"
)

func GetGanttGroupEvents(db *sql.DB, groupName string, dateOnly bool) ([]models.GanttEvent, error) {
	// TODO consider adding ordering by "end" DESC
	var events []models.GanttEvent
	query := `SELECT
	  "id",
		start,
		"end",
		"group",
		allDay,
		title,
		url,
		description
	FROM events
	WHERE "group" = $1
	ORDER BY "id";`
	rows, err := db.Query(query, groupName)
	if err != nil {
		return events, fmt.Errorf("failed to get events: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var e models.GanttEvent
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
		// log.Fatalf("Error reading groups: %v", err)
		return groups, err
	}
	return groups, nil
}
