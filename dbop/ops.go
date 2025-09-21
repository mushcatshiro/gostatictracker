package dbop

import (
	"database/sql"
	"fmt"

	"github.com/mushcatshiro/gostatictracker/models"
)

func InitDB(db *sql.DB, truncate bool) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		start TIMESTAMP,
		"end" TIMESTAMP,
		actualStart TIMESTAMP,
		actualEnd TIMESTAMP,
		insertTime TIMESTAMP NOT NULL,
		"group" TEXT DEFAULT 'default' NOT NULL,
		allDay BOOLEAN DEFAULT FALSE,
		title TEXT NOT NULL,
		url TEXT,
		description TEXT,
		pid INT,
		priority INT,
		metadata TEXT,
		status INT
	);`
	if truncate {
		_, err := db.Exec("TRUNCATE TABLE table_name RESTART IDENTITY CASCADE")
		if err != nil {
			return err
		}
	}
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func readEventById(db *sql.DB, id int64) (models.Event, error) {
	var e models.Event
	readQuery := `SELECT * FROM events WHERE id = $1`
	err := db.QueryRow(readQuery, id).Scan(&e)
	if err != nil {
		return models.Event{}, err
	}
	return e, nil
}

func InsertEvent(db *sql.DB, event models.Event) (int64, error) {
	insertQuery := `
	INSERT INTO events (
		start, "end", actualStart, actualEnd, insertTime, "group", allDay, title,
		url, description, pid, priority, metadata, status
	)
	VALUES (
	  $1, $2, $3, $4, $5,
		$6, $7, $8, $9, $10, $11, $12, $13, $14
	)
	RETURNING id;`

	var id int64
	err := db.QueryRow(insertQuery,
		event.Start,
		event.End,
		event.ActualStart,
		event.ActualEnd,
		event.InsertTime,
		event.Group,
		event.AllDay,
		event.Title,
		event.URL,
		event.Description,
		event.PID,
		event.Priority,
		event.Metadata,
		event.Status,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateEvent(db *sql.DB, event models.Event) error {
	updateQuery := `
	UPDATE events
	SET (
		start = TO_TIMESTAMP($1, 'MM-DD-YYYY HH24:MI'),
		"end" = TO_TIMESTAMP($2, 'MM-DD-YYYY HH24:MI'),
		actualStart = TO_TIMESTAMP($3, 'MM-DD-YYYY HH24:MI'),
		actualEnd = TO_TIMESTAMP($4, 'MM-DD-YYYY HH24:MI'),
		insertTime = TO_TIMESTAMP($5, 'MM-DD-YYYY HH24:MI'),
		"group" = $6,
		allDay = $7,
		title = $8,
		url = $9,
		description = $10
		pid = $11,
		priority = $12,
		metadata = $13,
		status = $14
	)
	WHERE id = $15;`

	result, err := db.Exec(updateQuery,
		event.Start,
		event.End,
		event.ActualStart,
		event.ActualEnd,
		event.InsertTime,
		event.Group,
		event.AllDay,
		event.Title,
		event.URL,
		event.Description,
		event.PID,
		event.Priority,
		event.Metadata,
		event.Status,
		event.ID,
	)

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
