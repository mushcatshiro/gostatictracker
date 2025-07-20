package dbop

import (
	"database/sql"
	"fmt"
	"log"
)

func InitDB(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		start TIMESTAMP NOT NULL,
		"end" TIMESTAMP,
		actualEnd TIMESTAMP,
		"group" TEXT DEFAULT 'default' NOT NULL,
		allDay BOOLEAN DEFAULT FALSE,
		title TEXT NOT NULL,
		url TEXT,
		description TEXT
	);`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func InsertEvent(db *sql.DB, event Event) (int64, error) {
	insertQuery := `
	INSERT INTO events (start, "end", actualEnd, "group", allDay, title, url, description)
	VALUES (
	  TO_TIMESTAMP($1, 'MM-DD-YYYY HH24:MI'),
		TO_TIMESTAMP($2, 'MM-DD-YYYY HH24:MI'),
		TO_TIMESTAMP($3, 'MM-DD-YYYY HH24:MI'),
		$4,
		$5,
		$6,
		$7,
		$8
	)
	RETURNING id;`

	var id int64
	err := db.QueryRow(insertQuery,
		event.Start,
		event.End,
		event.ActualEnd,
		event.Group,
		event.AllDay,
		event.Title,
		event.URL,
		event.Description).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateEvent(db *sql.DB, event Event) error {
	updateQuery := `
	UPDATE events
	SET start = $1, "end" = $2, "group" = $3, allDay = $4, title = $5, url = $6, description = $7
	WHERE id = $8;`

	result, err := db.Exec(updateQuery,
		event.Start,
		event.End,
		event.Group,
		event.AllDay,
		event.Title,
		event.URL,
		event.Description,
		event.ID)

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no event found with ID %d", event.ID)
	}

	return nil
}

func GetGanttGroupEvents(db *sql.DB, groupName string, dateOnly bool) ([]Event, error) {
	// TODO consider adding ordering by "end" DESC
	var events []Event
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
		actualEnd,
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
		var e Event
		if err := rows.Scan(&e.ID, &e.Start, &e.End, &e.ActualEnd, &e.Group, &e.AllDay, &e.Title, &e.URL, &e.Description); err != nil {
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
