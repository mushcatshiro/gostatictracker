package dbop

import (
	"database/sql"
	"fmt"
)

type ListEntry struct {
	ID         int64
	InsertTime string
	Group      string
	Title      string
	Priority   int8
	Status     int8
}

func GetListGroupEntries(db *sql.DB, groupName string) ([]ListEntry, error) {
	// TODO consider adding ordering by "end" DESC
	var events []ListEntry
	query := fmt.Sprintf(`SELECT
	  "id",
		TO_CHAR(insertTime, '%s'),
		"group",
		title,
		priority,
		status
	FROM events
	WHERE "group" = $1
	ORDER BY insertTime;`, TimeLayout)
	rows, err := db.Query(query, groupName)
	if err != nil {
		return events, fmt.Errorf("failed to get events: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var e ListEntry
		if err := rows.Scan(&e.ID, &e.InsertTime, &e.Group, &e.Title, &e.Priority, &e.Status); err != nil {
			return events, fmt.Errorf("failed to scan event: %v", err)
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return events, fmt.Errorf("error reading events: %v", err)
	}
	return events, nil
}
