package dbop

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mushcatshiro/gostatictracker/models"
)

func (db *DB) GetGanttRenderMetadataV1(groupName string) (models.GanttRenderMetadata, error) {
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
	var gst, get AgnosticTime
	query := `SELECT
		MIN("start") as startDate,
		MAX("end") as endDate
	FROM records
	WHERE "group" = $1`
	row := db.QueryRow(query, groupName)
	if err := row.Scan(&gst, &get); err != nil {
		return g, fmt.Errorf("\nfailed to scan start/end time and full duration:\n%w", err)
	}
	g.GroupStartTime = &gst.Time
	g.GroupEndTime = &get.Time

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

func (db *DB) GetGanttRenderMetadata(groupName string) (models.GanttRenderMetadata, error) {
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
	records, err := db.ReadRecords(models.Record{Group: groupName})
	if err != nil {
		return g, err
	}
	if len(records) == 0 {
		return g, fmt.Errorf("no groups found for %s", groupName)
	}
	var gst, get *time.Time
	minDuration := -1
	for _, r := range records {
		if r.Start == nil || r.End == nil {
			continue
		}
		if gst == nil || r.Start.Before(*gst) {
			gst = r.Start
		}
		if get == nil || r.End.After(*get) {
			get = r.End
		}
		startDay := r.Start.Truncate(24 * time.Hour)
		endDay := r.End.Truncate(24 * time.Hour)
		days := int(endDay.Sub(startDay).Hours()/24) + 1
		if minDuration == -1 || days < minDuration {
			minDuration = days
		}
	}
	g.GroupStartTime = gst
	g.GroupEndTime = get
	if minDuration >= 7 {
		g.IsDayView = false
		g.Divisor = 7 * 24
	}
	return g, nil
}
