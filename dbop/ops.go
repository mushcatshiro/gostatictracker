package dbop

import (
	"fmt"
	"strings"

	"github.com/mushcatshiro/gostatictracker/models"
)

func InitDB(db DBTX, truncate, drop bool) error {
	createEventTableQuery := `CREATE TABLE IF NOT EXISTS events (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		start TIMESTAMP,
		"end" TIMESTAMP,
		actualStart TIMESTAMP,
		actualEnd TIMESTAMP,
		insertTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"group" TEXT DEFAULT 'default' NOT NULL,
		allDay BOOLEAN DEFAULT FALSE,
		title TEXT NOT NULL,
		url TEXT,
		description TEXT,
		pid BIGINT,
		priority INT,
		metadata TEXT,
		status INT
	);`
	createUserTableQuery := `CREATE TABLE IF NOT EXISTS users (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		google_id TEXT NOT NULL UNIQUE,
		email TEXT,
		ipaddress TEXT,
		role TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	// deletes all rows while keeping table structure intact
	if truncate {
		_, err := db.Exec("TRUNCATE TABLE events RESTART IDENTITY CASCADE")
		if err != nil {
			return err
		}
		_, err = db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
		if err != nil {
			return err
		}
	}
	// removes the entire table definition with data, indexes and constraints
	if drop {
		_, err := db.Exec("DROP TABLE IF EXISTS events CASCADE")
		if err != nil {
			return err
		}
		_, err = db.Exec("DROP TABLE IF EXISTS users CASCADE")
		if err != nil {
			return err
		}
	}
	_, err := db.Exec(createEventTableQuery)
	if err != nil {
		return err
	}
	_, err = db.Exec(createUserTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func ReadEventById(db DBTX, id int64) (models.Event, error) {
	var e models.Event
	readQuery := `SELECT id, start, "end", actualStart, actualEnd, insertTime,
		"group", allDay, title, url, description, pid, priority, metadata, status
	FROM events WHERE id = $1`
	err := db.QueryRow(readQuery, id).Scan(
		&e.ID,
		&e.Start,
		&e.End,
		&e.ActualStart,
		&e.ActualEnd,
		&e.InsertTime,
		&e.Group,
		&e.AllDay,
		&e.Title,
		&e.URL,
		&e.Description,
		&e.PID,
		&e.Priority,
		&e.Metadata,
		&e.Status,
	)
	if err != nil {
		return models.Event{}, err
	}
	return e, nil
}

func ReadFilteredEvents(db DBTX, filterCols models.FilterCols) ([]models.Event, error) {
	var events []models.Event

	query := `SELECT id, start, "end", actualStart, actualEnd, insertTime,
		"group", allDay, title, url, description, pid, priority, metadata, status
	FROM events`
	var conditions []string
	var args []any
	if filterCols.Group != "" {
		conditions = append(conditions, fmt.Sprintf(`"group" = $%d`, len(conditions)+1))
		args = append(args, filterCols.Group)
	}
	if filterCols.Mode != "" {
		conditions = append(conditions, fmt.Sprintf(`"mode" = $%d`, len(conditions)+1))
		args = append(args, filterCols.Mode)
	}
	if filterCols.Status != nil {
		conditions = append(conditions, fmt.Sprintf(`"status" = $%d`, len(conditions)+1))
		args = append(args, &filterCols.Status)
	}
	if len(conditions) > 0 {
		query = query + " WHERE " + strings.Join(conditions, " AND ")
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return events, err
	}
	defer rows.Close()
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(
			&e.ID,
			&e.Start,
			&e.End,
			&e.ActualStart,
			&e.ActualEnd,
			&e.InsertTime,
			&e.Group,
			&e.AllDay,
			&e.Title,
			&e.URL,
			&e.Description,
			&e.PID,
			&e.Priority,
			&e.Metadata,
			&e.Status,
		); err != nil {
			return events, err
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return events, err
	}
	return events, nil
}

func GetUniqueGroups(db DBTX) ([]string, error) {
	var groups []string
	rows, err := db.Query(`SELECT DISTINCT "group" FROM events;`)
	if err != nil {
		return groups, err
	}
	defer rows.Close()

	for rows.Next() {
		var group string
		if err := rows.Scan(&group); err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		return groups, err
	}
	return groups, nil
}

func InsertEvent(db DBTX, event models.Event) (int64, error) {
	insertQuery := `
	INSERT INTO events (
		start, "end", actualStart, actualEnd, "group", allDay, title, url,
		description, pid, priority, metadata, status
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING id;`

	var id int64
	err := db.QueryRow(insertQuery,
		event.Start,
		event.End,
		event.ActualStart,
		event.ActualEnd,
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

func UpdateEvent(db DBTX, event models.Event) error {
	updateQuery := `UPDATE events
	SET start = $1, "end" = $2, actualStart = $3, actualEnd = $4, insertTime = $5,
		"group" = $6, allDay = $7, title = $8, url = $9, description = $10, pid = $11,
		priority = $12, metadata = $13, status = $14
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
