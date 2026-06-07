package dbop

import (
	"database/sql"
	"fmt"
	"log"

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

func GetGanttRenderMetadata(db *sql.DB, groupName string) (models.GanttRenderMetadata, error) {
	g := models.GanttRenderMetadata{
		IsDayView:            true,
		GroupName:            groupName,
		RowTextInRectPadding: 4,
		RectToTextMargin:     2,
		TextYOffset:          6.3,
		RowRectMargin:        2,
		RowRectHeight:        10,
		HeadersOffset:        (24 + 2) * 3,
		BaseHeaderRectWidth:  30,
		HeaderRectHeight:     24,
		HeaderRectMargin:     2,
		HeaderTextYOffset:    2,
		Divisor:              24,
	}
	/*
		scans table twice however does not assume earliest time is in start and vice
		versa; also prior to pg16 alias is needed for subquery.
		`SELECT MIN(d) AS startDate, MAX(d) AS endDate
		FROM (SELECT "start" d FROM events WHERE "group" = $1 UNION ALL
		SELECT "end" d FROM events WHERE "group" = $1)`
	*/
	query := `SELECT
		MIN("start") as startDate,
		MAX("end") as endDate
	FROM events
	WHERE "group" = $1`
	row := db.QueryRow(query, groupName)
	if err := row.Scan(&g.GroupStartTime, &g.GroupEndTime); err != nil {
		return g, fmt.Errorf("\nfailed to scan start/end time and full duration:\n%w", err)
	}

	query = `SELECT MIN(("end"::date - start::date)+1) AS minTaskDuration FROM events WHERE "group" = $1`
	var minTaskDuration sql.NullInt64
	err := db.QueryRow(query, groupName).Scan(&minTaskDuration)
	if err != nil {
		if err == sql.ErrNoRows {
			return g, fmt.Errorf("\nno events found for group %s:\n%w", groupName, err)
		}
		log.Printf("Failed to scan minimum task duration: %v; using default day view", err)
	} else if minTaskDuration.Valid && minTaskDuration.Int64 >= 7 {
		g.IsDayView = false
		g.Divisor = 7 * 24
	}
	return g, nil
}
