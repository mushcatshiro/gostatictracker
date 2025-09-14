package dbop

import (
	"database/sql"
	"fmt"
	"time"
)

type CalendarEvent struct {
	Title string
	Group string
	Start string
	End   string
}

type MonthGroup struct {
	FirstDayOfMonth time.Time
	Events   []CalendarEvent
}

func getOneMonth(db *sql.DB, d time.Time) (MonthGroup, error) {
	var mg MonthGroup
	query := `
	WITH month_range AS (
		SELECT tsrange(
			make_timestamp(:$1, :$2, 1, 0, 0, 0),
			make_timestamp(:$1, :$2, 1, 0, 0, 0) + interval '1 month',
			'[)'
		) AS period
	)
	SELECT * FROM events, month_range
	WHERE
		(CASE
			WHEN "actualStart" IS NOT NULL AND "actualEnd" IS NOT NULL THEN
				tsrange("actualStart", "actualEnd", '[]') && month_range.period
			WHEN "start" IS NOT NULL OR "END" IS NOT NULL THEN
				tsrange("start", "end", '[]') && month_range.period
			ELSE
				"insertTime" >= LOWER(month_range.period) AND "insertTime" < upper(month_range.period)
		END)
		AND
			(CASE
				WHEN "actualStart" IS NOT NULL AND "actualEnd" IS NOT NULL THEN
					("actualEnd" - "actualStart") <= interval '7 days'
				WHEN ("actualStart" IS NULL OR "actualEND" IS NULL) AND ("start" IS NOT NULL AND "end" IS NOT NULL) THEN
					("end" - "start") <= interval '7 days'
				ELSE
					TRUE
			END);`
	rows, err := db.Query(query, d.Year(), d.Month())
	if err != nil {
		return mg, err
	}
	defer rows.Close()
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

func GetCalendarMonthGroups(conn *sql.DB, month, year int) ([]MonthGroup, error) {
	var monthGroups []MonthGroup

	renderRange, err := getCalendarRenderRange(month, year)
	if err != nil {
		return monthGroups, err
	}

	for _, d := range renderRange {
		mg, err := getOneMonth(conn, d)
		if err != nil {
			return monthGroups, err
		}
		monthGroups = append(monthGroups, mg)
	}
	return monthGroups, nil
}
