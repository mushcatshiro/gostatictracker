package dbop

import (
	"database/sql"
	"fmt"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/models"
)

func GetListGroupEntries(db *sql.DB, groupName string) ([]models.ListEntry, error) {
	// TODO consider adding ordering by "end" DESC
	var events []models.ListEntry
	query := fmt.Sprintf(`SELECT
	  "id",
		TO_CHAR(insertTime, '%s'),
		"group",
		title,
		priority,
		status
	FROM events
	WHERE "group" = $1
	ORDER BY insertTime;`, common.TimeLayout)
	rows, err := db.Query(query, groupName)
	if err != nil {
		return events, fmt.Errorf("failed to get events: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var e models.ListEntry
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

func UpdateStatus(db *sql.DB, id int64, status common.Status) error {

	e, err := ReadEventById(db, id)
	if err != nil {
		return err
	}
	e.Status = status
	/*
		if status == COMPLETED {
			ae := time.Now()
			e.ActualEnd = &ae
		}
	*/
	id, err = InsertEvent(db, e)
	if err != nil {
		return err
	}
	return nil
}
