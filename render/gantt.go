package render

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mushcatshiro/gostatictracker/dbop"
)

const TimeLayout = "01-02-2006" // MM-DD-YYYY

type eventGanttHeader struct {
	RectX     int
	RectWidth int
	TextX     int
	TextY     int
	TextVal   string
	DataDate  string
}

type eventGanttRow struct {
	RectX       int
	RectY       int
	RectWidth   int
	DataDetails string
	LineX1      int
	LineY1      int
	LineX2      int
	LineY2      int
	TextX       int
	TextY       float32
	TextVal     string
}

type eventGanttBase struct {
	HeaderRectWidth  int
	HeaderRectHeight int
	RowRectHeight    int
	SvgWidth         int
	SvgHeight        int
	Group            string
	Years            []eventGanttHeader
	Months           []eventGanttHeader
	Days             []eventGanttHeader
	Rows             []eventGanttRow
}

type ganttRenderMetadata struct {
	isDayView            bool
	fullTaskDuration     int
	groupStartTime       time.Time
	groupEndTime         time.Time
	groupName            string
	rowTextInRectPadding int     // LR padding for text to be place in rect
	rectToTextMargin     int     // spacing between rect and text
	textYOffset          float32 // approximate text Y offset from Y (center)
	rowRectMargin        int     // spacing between two rect/row
	rowRectHeight        int
	headersOffset        int
	baseHeaderRectWidth  int // date or week rect width
	headerRectHeight     int
	headerRectMargin     int
	headerTextYOffset    int // somehow it needs +2 after dividing height by 2
}

func getGanttRenderMetadata(db *sql.DB, groupName string) (ganttRenderMetadata, error) {
	g := ganttRenderMetadata{
		isDayView:            true,
		groupName:            groupName,
		rowTextInRectPadding: 4,
		rectToTextMargin:     2,
		textYOffset:          6.3,
		rowRectMargin:        2,
		rowRectHeight:        10,
		headersOffset:        (24 + 2) * 3,
		baseHeaderRectWidth:  30,
		headerRectHeight:     24,
		headerRectMargin:     2,
		headerTextYOffset:    2,
	}
	query := `SELECT
		TO_CHAR(MIN(d), 'MM-DD-YYYY') AS startDate,
		TO_CHAR(MAX(d), 'MM-DD-YYYY') AS endDate,
		EXTRACT(DAY FROM MAX(d) - MIN(d)) + 1 AS taskFullDuration
	FROM (
		SELECT "start" d FROM events WHERE "group" = $1
		UNION ALL
		SELECT "end" d FROM events WHERE "group" = $1
	)`
	row := db.QueryRow(query, groupName)
	var groupStartTime, groupEndTime string
	if err := row.Scan(&groupStartTime, &groupEndTime, &g.fullTaskDuration); err != nil {
		return ganttRenderMetadata{}, fmt.Errorf("\nfailed to scan start/end time and full duration:\n%w", err)
	}
	var tGroupStartTime, tGroupEndTime time.Time
	tGroupStartTime, err := time.Parse(TimeLayout, groupStartTime)
	if err != nil {
		return g, fmt.Errorf("\nfailed to parse start time: %w", err)
	}
	tGroupEndTime, err = time.Parse(TimeLayout, groupEndTime)
	if err != nil {
		return g, fmt.Errorf("\nfailed to parse end time: %w", err)
	}
	g.groupStartTime = tGroupStartTime
	g.groupEndTime = tGroupEndTime

	query = `SELECT MIN(EXTRACT(DAY FROM "end" - start)) AS minTaskDuration FROM events WHERE "group" = $1`
	var minTaskDuration sql.NullInt64
	err = db.QueryRow(query, groupName).Scan(&minTaskDuration)
	if err != nil {
		if err == sql.ErrNoRows {
			return ganttRenderMetadata{}, fmt.Errorf("\nno events found for group %s:\n%w", groupName, err)
		}
		log.Printf("Failed to scan minimum task duration: %v; using default day view", err)
	} else if minTaskDuration.Valid && minTaskDuration.Int64 > 7 {
		g.isDayView = false
	}
	return g, nil
}

func getTextEstimateWidth(text string) int {
	return len([]rune(text)) * 7
}

func getRowRectWidth(startTime time.Time, endTime time.Time, headerWidthWithSpacing int) int {
	daysSpan := int(endTime.Sub(startTime).Hours()/24) + 1 // include starting day
	return daysSpan * headerWidthWithSpacing
}

func parseEventTimes(s, e, layout string) (time.Time, time.Time, error) {
	iStartTime, err := time.Parse(TimeLayout, s)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("not able to parse start time: %w", err)
	}
	iEndTime, err := time.Parse(TimeLayout, e)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("not able to parse end time: %w", err)
	}
	return iStartTime, iEndTime, nil
}

func processRectAndLine(spacing, rowIdx, rowRectWidth int, g ganttRenderMetadata) (int, int, int, int) {
	rectX := spacing * (g.baseHeaderRectWidth + g.headerRectMargin)
	rectY := (rowIdx * (g.rowRectHeight + g.rowRectMargin)) + g.headersOffset
	lineX1 := rectX + rowRectWidth
	lineX2 := rectX + rowRectWidth
	return rectX, rectY, lineX1, lineX2
}

func processOverFlowDays(maxWidth, rowEndWidth int, g ganttRenderMetadata) int {
	overFlowDays := float64(maxWidth-rowEndWidth) / float64(g.baseHeaderRectWidth+g.headerRectMargin)
	if overFlowDays < 0.2 {
		return 0
	}
	return int(math.Ceil(overFlowDays))
}

func getGanttEventRows(events []dbop.Event, g ganttRenderMetadata, debug bool) ([]eventGanttRow, int, error) {
	// TODO
	// support actual start, actual end rendering including start, stop button
	rSlice := []eventGanttRow{}

	var maxWidth, rowEndWidth int

	lineY1 := g.headersOffset - 2
	lineY2 := g.headersOffset + len(events)*(g.rowRectHeight+g.rowRectMargin)

	for idx, event := range events {
		iStartTime, iEndTime, err := parseEventTimes(event.Start, event.End, TimeLayout)
		if err != nil {
			return rSlice, -1, fmt.Errorf("processing event[%d]:\n%w", idx, err)
		}
		duration := iStartTime.Sub(g.groupStartTime)
		daysSpacing := int(duration.Hours() / 24)

		titleWidth := getTextEstimateWidth(event.Title)
		rowRectWidth := getRowRectWidth(iStartTime, iEndTime, g.baseHeaderRectWidth+g.headerRectMargin)

		var textX int
		var textY float32

		rectX, rectY, lineX1, lineX2 := processRectAndLine(daysSpacing, idx, rowRectWidth, g)

		// keeping text here since `maxWidth` is coupled with `textX`
		textY = float32(rectY) + g.textYOffset
		rowEndWidth = max(rowEndWidth, rectX+rowRectWidth)
		if rowRectWidth-g.rowRectMargin < titleWidth {
			// BUG titleWidth might not be divisible by 2 and will always floored, might get inconsistent spacing
			textX = rectX + rowRectWidth + g.rectToTextMargin + (titleWidth / 2)
			maxWidth = max(maxWidth, rectX+rowRectWidth+(g.rectToTextMargin*2)+titleWidth)
		} else {
			textX = (rowRectWidth / 2) + rectX
			maxWidth = max(maxWidth, rectX+rowRectWidth+g.rectToTextMargin)
		}

		var description string
		if debug {
			// BUG for now
			description = event.Description + "\n" + event.Start + "\n" + event.End
		} else {
			description = event.Description
		}
		e := eventGanttRow{
			RectX:       rectX,
			RectY:       rectY,
			RectWidth:   rowRectWidth,
			DataDetails: description,
			LineX1:      lineX1,
			LineY1:      lineY1,
			LineX2:      lineX2,
			LineY2:      lineY2,
			TextX:       textX,
			TextY:       textY,
			TextVal:     event.Title,
		}
		rSlice = append(rSlice, e)
	}

	overFlowDays := processOverFlowDays(maxWidth, rowEndWidth, g)
	return rSlice, overFlowDays, nil
}

func getGanttHeaders(g ganttRenderMetadata, overflowDays int) ([]eventGanttHeader, []eventGanttHeader, []eventGanttHeader) {
	startTime := g.groupStartTime
	endTime := g.groupEndTime.AddDate(0, 0, overflowDays)

	yearGanttHeader := []eventGanttHeader{}
	monthGanttHeader := []eventGanttHeader{}
	dayGanttHeader := []eventGanttHeader{}

	var iterYear, iterDay, dIdx, monthRectEndX, yearRectEndX int
	var cumDaysForMonth, cumDaysForYear int
	var iterMonth time.Month

	for !startTime.After(endTime) {
		iterYear, iterMonth, iterDay = startTime.Date()

		d := eventGanttHeader{
			RectX:     dIdx * (g.baseHeaderRectWidth + g.headerRectMargin),
			RectWidth: g.baseHeaderRectWidth,
			TextX:     dIdx*(g.baseHeaderRectWidth+g.headerRectMargin) + int(g.baseHeaderRectWidth/2),
			TextY:     2*(g.headerRectHeight+g.headerRectMargin) + int(g.headerRectHeight/2) + g.headerTextYOffset,
			TextVal:   strconv.Itoa(iterDay),
			DataDate:  startTime.Format("2006-01-02"), // YYYY-MM-DD
		}
		dayGanttHeader = append(dayGanttHeader, d)
		dIdx++
		cumDaysForMonth++
		cumDaysForYear++

		if startTime.AddDate(0, 0, 1).Month() != startTime.Month() {
			currentMonthWidth := (cumDaysForMonth * (g.baseHeaderRectWidth + g.headerRectMargin)) - g.headerRectMargin
			m := eventGanttHeader{
				RectX:     monthRectEndX,
				RectWidth: currentMonthWidth,
				TextX:     monthRectEndX + int(currentMonthWidth/2),
				TextY:     (g.headerRectHeight + g.headerRectMargin) + int(g.headerRectHeight/2) + g.headerTextYOffset,
				TextVal:   strconv.Itoa(int(iterMonth)),
			}
			monthGanttHeader = append(monthGanttHeader, m)
			monthRectEndX += currentMonthWidth + g.headerRectMargin
			cumDaysForMonth = 0
		}
		if startTime.AddDate(0, 0, 1).Year() != startTime.Year() {
			currentYearWidth := (cumDaysForYear * (g.baseHeaderRectWidth + g.headerRectMargin)) - g.headerRectMargin
			y := eventGanttHeader{
				RectX:     yearRectEndX,
				RectWidth: currentYearWidth,
				TextX:     yearRectEndX + int(currentYearWidth/2),
				TextY:     int(g.headerRectHeight/2) + g.headerTextYOffset,
				TextVal:   strconv.Itoa(iterYear)}
			yearGanttHeader = append(yearGanttHeader, y)
			yearRectEndX += currentYearWidth + g.headerRectMargin
			cumDaysForYear = 0
		}

		startTime = startTime.AddDate(0, 0, 1)
	}
	if cumDaysForMonth > 0 {
		currentMonthWidth := (cumDaysForMonth * (g.baseHeaderRectWidth + g.headerRectMargin)) - g.headerRectMargin
		m := eventGanttHeader{
			RectX:     monthRectEndX,
			RectWidth: currentMonthWidth,
			TextX:     monthRectEndX + int(currentMonthWidth/2),
			TextY:     (g.headerRectHeight + g.headerRectMargin) + int(g.headerRectHeight/2) + g.headerTextYOffset,
			TextVal:   strconv.Itoa(int(iterMonth)),
		}
		monthGanttHeader = append(monthGanttHeader, m)
	}
	if cumDaysForYear > 0 {
		currentYearWidth := (cumDaysForYear * (g.baseHeaderRectWidth + g.headerRectMargin)) - g.headerRectMargin
		y := eventGanttHeader{
			RectX:     yearRectEndX,
			RectWidth: currentYearWidth,
			TextX:     yearRectEndX + int(currentYearWidth/2),
			TextY:     int(g.headerRectHeight/2) + g.headerTextYOffset,
			TextVal:   strconv.Itoa(iterYear)}
		yearGanttHeader = append(yearGanttHeader, y)
	}
	return yearGanttHeader, monthGanttHeader, dayGanttHeader
}

func getGanttEventBase(events []dbop.Event, padding int, g ganttRenderMetadata, debug bool) (eventGanttBase, error) {
	// allow headerRectWidth, headerRectHeight, rowRectHeight to have default values
	e := eventGanttBase{
		HeaderRectWidth:  g.baseHeaderRectWidth,
		HeaderRectHeight: g.headerRectHeight,
		RowRectHeight:    g.rowRectHeight,
		Group:            g.groupName,
	}
	e.SvgHeight = g.headersOffset + len(events)*(e.RowRectHeight+padding)

	rSlice, overflowDays, err := getGanttEventRows(events, g, debug)
	if err != nil {
		return eventGanttBase{}, err
	}
	e.Rows = rSlice

	e.SvgWidth = (g.fullTaskDuration+overflowDays)*(e.HeaderRectWidth+g.headerRectMargin) + 1 // buffer

	y, m, d := getGanttHeaders(g, overflowDays)

	e.Years = y
	e.Months = m
	e.Days = d
	return e, nil
}

func renderGanttHTML(events []dbop.Event, file *os.File, g ganttRenderMetadata, debug bool) error {
	var t *template.Template
	var err error

	if g.isDayView {
		// t, err = template.ParseFiles("dayGantt.html")
		t, err = template.ParseGlob(filepath.Join(".", "day*.html"))
	} else {
		// t, err = template.ParseFiles("weekGantt.html")
		return fmt.Errorf("not supporting week view gantt")
	}

	if err != nil {
		return fmt.Errorf("failed to create template:\n\t%w", err)
	}
	data, err := getGanttEventBase(events, 2, g, debug)
	if err != nil {
		return err
	}
	// fmt.Printf("%+v\n", data)
	err = t.Execute(file, data)
	return err
}

func RenderGantt(renderTargetPath string, conn *sql.DB, debug bool) {
	groups, err := dbop.GetUniqueGroups(conn)
	if err != nil {
		log.Fatalf("Not able to query any group(s): %v", err)
	}
	for _, group := range groups {
		log.Printf("Processing group: %s", group)
		events, err := dbop.GetGanttGroupEvents(conn, group, true)
		if err != nil {
			log.Printf("Failed to get gantt group events: %v", err)
			continue
		}
		fileName := strings.Replace(group, " ", "-", -1) + "-gantt.html"
		file, err := os.Create(renderTargetPath + "/" + fileName) // truncates if exists
		if err != nil {
			log.Printf("Failed to create file %s: %v", fileName, err)
			continue
		}
		/*
			align to go idiom, only `Close` when file is created successfully;
			when err != nil, file is likely `nil`
			`nil.Close()` is potentially nil pointer dereference
		*/
		defer file.Close()

		g, err := getGanttRenderMetadata(conn, group)
		if err != nil {
			log.Fatalf("Failed to get gantt render metadata for group %s: %v", group, err)
		}
		err = renderGanttHTML(events, file, g, debug)
		if err != nil {
			log.Printf("Failed to process %s:\n\t%v", fileName, err)
		}
	}
}
