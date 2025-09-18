package dbop

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mushcatshiro/gostatictracker/models"
)

func getOneMonth(db *sql.DB, d time.Time) (models.MonthGroup, error) {
	var mg models.MonthGroup
	query := `
	WITH month_range AS (
		SELECT tsrange(
			make_timestamp($1, $2, 1, 0, 0, 0),
			make_timestamp($1, $2, 1, 0, 0, 0) + interval '1 month',
			'[)'
		) AS period
	)
	SELECT
		title,
		"group",
		CASE
			WHEN "actualstart" IS NOT NULL AND "actualend" IS NOT NULL THEN "actualstart"
			WHEN "start" IS NOT NULL OR "end" IS NOT NULL THEN "start"
			ELSE "inserttime"
		END as s,
		CASE
			WHEN "actualstart" IS NOT NULL AND "actualend" IS NOT NULL THEN "actualend"
			WHEN "start" IS NOT NULL OR "end" IS NOT NULL THEN "end"
			ELSE "inserttime"
		END as e
	FROM events, month_range
	WHERE
		(CASE
			WHEN "actualstart" IS NOT NULL AND "actualend" IS NOT NULL THEN
				tsrange("actualstart", "actualend", '[]') && month_range.period
			WHEN "start" IS NOT NULL OR "end" IS NOT NULL THEN
				tsrange("start", "end", '[]') && month_range.period
			ELSE
				"inserttime" >= LOWER(month_range.period) AND "inserttime" < upper(month_range.period)
		END)
		AND
			(CASE
				WHEN "actualstart" IS NOT NULL AND "actualend" IS NOT NULL THEN
					("actualend" - "actualstart") <= interval '7 days'
				WHEN ("actualstart" IS NULL OR "actualend" IS NULL) AND ("start" IS NOT NULL AND "end" IS NOT NULL) THEN
					("end" - "start") <= interval '7 days'
				ELSE
					TRUE
			END);`
	rows, err := db.Query(query, d.Year(), d.Month())
	if err != nil {
		return mg, err
	}
	defer rows.Close()
	var ctr int
	for rows.Next() {
		var e models.CalendarEvent
		if err := rows.Scan(&e.Title, &e.Group, &e.Start, &e.End); err != nil {
			return mg, err
		}
		ctr += 1
		mg.Events = append(mg.Events, e)
	}
	if ctr < 1 {
		return mg, fmt.Errorf("No rows is returned")
	}
	return mg, nil
}

func getCalendarRenderRange(month, year int) ([]time.Time, error) {
	/*
		probably some of the checks need to move up to cli/UI level;
	*/
	if month > 12 {
		return nil, fmt.Errorf("Month must not be greater than 12")
	}
	if year < 0 {
		return nil, fmt.Errorf("Year must be greater or equal to 0")
	}
	if month != 0 && year == 0 {
		return nil, fmt.Errorf("Does not support unspecified year with specific month given")
	}
	var ret []time.Time
	var s, e time.Time
	if month == 0 && year != 0 {
		s = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		e = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	} else if month != 0 && year != 0 {
		s = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		e = time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
	} else {
		// default: render today +/- 6 months
		tdy := time.Now()
		s = tdy.AddDate(0, -6, 0)
		e = tdy.AddDate(0, 6, 0)
	}
	for d := s; !d.After(e); d = d.AddDate(0, 1, 0) {
		ret = append(ret, d)
	}
	return ret, nil
}

func GetCalendarMonthGroups(conn *sql.DB, month, year int) ([]models.MonthGroup, error) {
	var monthGroups []models.MonthGroup
	fmt.Printf("input %v, %v", month, year)

	renderRange, err := getCalendarRenderRange(month, year)
	if err != nil {
		return monthGroups, err
	}

	for _, d := range renderRange {
		fmt.Printf("processing %v, %v", d.Month(), d.Year())
		mg, err := getOneMonth(conn, d)
		if err != nil {
			return monthGroups, err
		}
		monthGroups = append(monthGroups, mg)
	}
	return monthGroups, nil
}
