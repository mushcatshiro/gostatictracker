package render

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/snabb/isoweek"
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
	groupStartTime       *time.Time
	groupEndTime         *time.Time
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
	divisor              int
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
		divisor:              24,
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
	if err := row.Scan(&g.groupStartTime, &g.groupEndTime); err != nil {
		return ganttRenderMetadata{}, fmt.Errorf("\nfailed to scan start/end time and full duration:\n%w", err)
	}

	query = `SELECT MIN(("end"::date - start::date)+1) AS minTaskDuration FROM events WHERE "group" = $1`
	var minTaskDuration sql.NullInt64
	err := db.QueryRow(query, groupName).Scan(&minTaskDuration)
	if err != nil {
		if err == sql.ErrNoRows {
			return ganttRenderMetadata{}, fmt.Errorf("\nno events found for group %s:\n%w", groupName, err)
		}
		log.Printf("Failed to scan minimum task duration: %v; using default day view", err)
	} else if minTaskDuration.Valid && minTaskDuration.Int64 >= 7 {
		g.isDayView = false
		g.divisor = 7 * 24
	}
	return g, nil
}

func getTextEstimateWidth(text string) int {
	return len([]rune(text)) * 8
}

func getRowRectWidth(startTime time.Time, endTime time.Time, headerWidthWithSpacing int, divisor float64) int {
	daysSpan := int(endTime.Sub(startTime).Hours()/divisor) + 1 // include starting day
	return daysSpan * headerWidthWithSpacing
}

func formatEventTimes(s, e *time.Time, isDayView bool) (time.Time, time.Time) {
	var iStartTime, iEndTime time.Time
	if !isDayView {
		sy, sw := s.ISOWeek()
		iStartTime = isoweek.StartTime(sy, sw, time.UTC)
		ey, ew := e.ISOWeek()
		iEndTime = isoweek.StartTime(ey, ew, time.UTC) // is `time.UTC` the correct time location?
	} else {
		iStartTime = *s
		iEndTime = *e
	}
	return iStartTime, iEndTime
}

func processRectAndLine(spacing, rowIdx, rowRectWidth int, g ganttRenderMetadata) (int, int, int, int) {
	rectX := spacing * (g.baseHeaderRectWidth + g.headerRectMargin)
	rectY := (rowIdx * (g.rowRectHeight + g.rowRectMargin)) + g.headersOffset
	lineX1 := rectX + rowRectWidth
	lineX2 := rectX + rowRectWidth
	return rectX, rectY, lineX1, lineX2
}

func processOverFlowUnits(maxWidth, rowEndWidth int, g ganttRenderMetadata) int {
	// visually still blocking last few chars due to padding 0.5rem
	diff := maxWidth - rowEndWidth
	var overflowUnits int
	if diff <= 0 {
		return overflowUnits
	}

	for diff > 0 {
		diff = diff - (g.baseHeaderRectWidth + g.headerRectMargin)
		overflowUnits++
	}

	return overflowUnits
}

func getGanttRows(events []models.GanttEvent, g ganttRenderMetadata, debug bool) ([]eventGanttRow, int) {
	// TODO
	// support actual start, actual end rendering including start, stop button
	rSlice := []eventGanttRow{}

	var maxWidth, rowEndWidth int

	lineY1 := g.headersOffset - 2
	lineY2 := g.headersOffset + len(events)*(g.rowRectHeight+g.rowRectMargin)

	for idx, event := range events {
		iStartTime, iEndTime := formatEventTimes(event.Start, event.End, g.isDayView)
		if !g.isDayView && idx == 0 {
			g.groupStartTime = &iStartTime
		}
		duration := iStartTime.Sub(*g.groupStartTime)
		daysSpacing := int(duration.Hours() / float64(g.divisor))

		titleWidth := getTextEstimateWidth(event.Title)
		rowRectWidth := getRowRectWidth(iStartTime, iEndTime, g.baseHeaderRectWidth+g.headerRectMargin, float64(g.divisor))

		var textX int
		var textY float32

		rectX, rectY, lineX1, lineX2 := processRectAndLine(daysSpacing, idx, rowRectWidth, g)

		// keeping text here since `maxWidth` is coupled with `textX`
		textY = float32(rectY) + g.textYOffset
		rowEndWidth = max(rowEndWidth, rectX+rowRectWidth)
		if rowRectWidth-g.rowRectMargin-(2*g.rectToTextMargin) > titleWidth {
			textX = rectX + g.rectToTextMargin
			maxWidth = max(maxWidth, rectX+rowRectWidth)
		} else {
			textX = rectX + rowRectWidth + g.rectToTextMargin
			maxWidth = max(maxWidth, rectX+rowRectWidth+(g.rectToTextMargin*2)+titleWidth)
		}

		var description string
		if debug {
			description = event.Description + "\n" +
				event.Start.Format(common.TimeLayout) + "\n" +
				event.End.Format(common.TimeLayout)
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

	overFlowUnits := processOverFlowUnits(maxWidth, rowEndWidth, g)
	return rSlice, overFlowUnits
}

func getGanttHeaders(g ganttRenderMetadata, overflowUnits int) ([]eventGanttHeader, []eventGanttHeader, []eventGanttHeader) {
	startTime := *g.groupStartTime
	endTime := *g.groupEndTime

	var overflowDays int
	if g.isDayView {
		overflowDays = overflowUnits
	} else {
		for overflowDays < overflowUnits*7 {
			overflowDays = overflowDays + 7
		}

	}
	endTime = endTime.AddDate(0, 0, overflowDays)
	if !g.isDayView {
		ey, ew := endTime.ISOWeek()
		endTime = isoweek.StartTime(ey, ew, time.UTC)
	}

	yearGanttHeader := []eventGanttHeader{}
	monthGanttHeader := []eventGanttHeader{}
	dayGanttHeader := []eventGanttHeader{}

	var iterYear, iterDay, iterWeek, trackWeek, dIdx, monthRectEndX, yearRectEndX int
	var cumForMonth, cumForYear int
	var iterMonth time.Month

	for !startTime.After(endTime) {
		iterYear, iterMonth, iterDay = startTime.Date()
		_, iterWeek = startTime.ISOWeek()

		var d eventGanttHeader
		if g.isDayView {
			d = eventGanttHeader{
				RectX:     dIdx * (g.baseHeaderRectWidth + g.headerRectMargin),
				RectWidth: g.baseHeaderRectWidth,
				TextX:     dIdx*(g.baseHeaderRectWidth+g.headerRectMargin) + int(g.baseHeaderRectWidth/2),
				TextY:     2*(g.headerRectHeight+g.headerRectMargin) + int(g.headerRectHeight/2) + g.headerTextYOffset,
				TextVal:   strconv.Itoa(iterDay),
				DataDate:  startTime.Format("2006-01-02"), // YYYY-MM-DD, BUG?
			}
			dayGanttHeader = append(dayGanttHeader, d)
			dIdx++
			cumForMonth++
			cumForYear++
		} else {
			if trackWeek != iterWeek {
				// limitation: last monday of the year + remaining days in jan,
				// startTime will result in the next year's first day
				t := isoweek.StartTime(iterYear, iterWeek, time.UTC)
				_, _, iterDay = t.Date()
				d = eventGanttHeader{
					RectX:     dIdx * (g.baseHeaderRectWidth + g.headerRectMargin),
					RectWidth: g.baseHeaderRectWidth,
					TextX:     dIdx*(g.baseHeaderRectWidth+g.headerRectMargin) + int(g.baseHeaderRectWidth/2),
					TextY:     2*(g.headerRectHeight+g.headerRectMargin) + int(g.headerRectHeight/2) + g.headerTextYOffset,
					TextVal:   strconv.Itoa(iterDay),
					DataDate:  startTime.Format("2006-01-02"), // YYYY-MM-DD, BUG?
				}
				trackWeek = iterWeek
				dayGanttHeader = append(dayGanttHeader, d)
				dIdx++
				cumForMonth++
				cumForYear++
			}
		}

		if startTime.AddDate(0, 0, 1).Month() != startTime.Month() {
			currentMonthWidth := (cumForMonth * (g.baseHeaderRectWidth + g.headerRectMargin)) - g.headerRectMargin
			m := eventGanttHeader{
				RectX:     monthRectEndX,
				RectWidth: currentMonthWidth,
				TextX:     monthRectEndX + int(currentMonthWidth/2),
				TextY:     (g.headerRectHeight + g.headerRectMargin) + int(g.headerRectHeight/2) + g.headerTextYOffset,
				TextVal:   strconv.Itoa(int(iterMonth)),
			}
			monthGanttHeader = append(monthGanttHeader, m)
			monthRectEndX += currentMonthWidth + g.headerRectMargin
			cumForMonth = 0
		}
		if startTime.AddDate(0, 0, 1).Year() != startTime.Year() {
			currentYearWidth := (cumForYear * (g.baseHeaderRectWidth + g.headerRectMargin)) - g.headerRectMargin
			y := eventGanttHeader{
				RectX:     yearRectEndX,
				RectWidth: currentYearWidth,
				TextX:     yearRectEndX + int(currentYearWidth/2),
				TextY:     int(g.headerRectHeight/2) + g.headerTextYOffset,
				TextVal:   strconv.Itoa(iterYear)}
			yearGanttHeader = append(yearGanttHeader, y)
			yearRectEndX += currentYearWidth + g.headerRectMargin
			cumForYear = 0
		}

		startTime = startTime.AddDate(0, 0, 1)
	}
	if cumForMonth > 0 {
		currentMonthWidth := (cumForMonth * (g.baseHeaderRectWidth + g.headerRectMargin)) - g.headerRectMargin
		m := eventGanttHeader{
			RectX:     monthRectEndX,
			RectWidth: currentMonthWidth,
			TextX:     monthRectEndX + int(currentMonthWidth/2),
			TextY:     (g.headerRectHeight + g.headerRectMargin) + int(g.headerRectHeight/2) + g.headerTextYOffset,
			TextVal:   strconv.Itoa(int(iterMonth)),
		}
		monthGanttHeader = append(monthGanttHeader, m)
	}
	if cumForYear > 0 {
		currentYearWidth := (cumForYear * (g.baseHeaderRectWidth + g.headerRectMargin)) - g.headerRectMargin
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

func getGanttEventBase(events []models.GanttEvent, g ganttRenderMetadata, debug bool) eventGanttBase {
	// allow headerRectWidth, headerRectHeight, rowRectHeight to have default values
	e := eventGanttBase{
		HeaderRectWidth:  g.baseHeaderRectWidth,
		HeaderRectHeight: g.headerRectHeight,
		RowRectHeight:    g.rowRectHeight,
		Group:            g.groupName,
	}
	e.SvgHeight = g.headersOffset + len(events)*(e.RowRectHeight+g.rowRectMargin) + 2*8

	rSlice, overflowUnits := getGanttRows(events, g, debug)
	e.Rows = rSlice

	y, m, d := getGanttHeaders(g, overflowUnits)
	e.SvgWidth = len(d)*(e.HeaderRectWidth+g.headerRectMargin) + 17 // buffer

	e.Years = y
	e.Months = m
	e.Days = d
	return e
}

func RenderGanttHTML(events []models.GanttEvent, file *os.File, g ganttRenderMetadata, debug bool) error {
	var t *template.Template
	var err error

	t, err = template.ParseGlob(filepath.Join(".", "day*.html"))

	if err != nil {
		return fmt.Errorf("failed to create template:\n\t%w", err)
	}
	data := getGanttEventBase(events, g, debug)
	err = t.Execute(file, data)
	return err
}

func buildGanttElmTree(events []models.GanttEvent, g ganttRenderMetadata) elm {
	// should only take in GanttEventBase
	geb := getGanttEventBase(events, g, false)
	gss := buildGanttStyleString(geb.HeaderRectWidth, geb.HeaderRectHeight, geb.RowRectHeight)
	svgE := svgElm
	svgE.attrs.height = strconv.Itoa(geb.SvgHeight)
	svgE.attrs.width = strconv.Itoa(geb.SvgWidth)
	for _, egh := range geb.Years {
		e := buildHeaderGroup(egh.RectX, egh.RectWidth, egh.TextX, egh.TextY, egh.TextVal, "header year")
		svgE.childs = append(svgE.childs, e)
	}
	for _, egh := range geb.Months {
		e := buildHeaderGroup(egh.RectX, egh.RectWidth, egh.TextX, egh.TextY, egh.TextVal, "header month")
		svgE.childs = append(svgE.childs, e)
	}
	for _, egh := range geb.Days {
		e := buildHeaderGroup(egh.RectX, egh.RectWidth, egh.TextX, egh.TextY, egh.TextVal, "header date")
		svgE.childs = append(svgE.childs, e)
	}
	for _, egr := range geb.Rows {
		e := buildRowGroup(egr.RectX, egr.RectY, egr.RectWidth, egr.LineX1, egr.LineX2, egr.LineY1, egr.LineY2, egr.TextX, egr.TextY, egr.TextVal, "rows")
		svgE.childs = append(svgE.childs, e)
	}
	htmlBody := bodyElm
	headerElm := elm{tag: "h2", innerText: g.groupName}
	htmlBody.childs = append(htmlBody.childs, headerElm)
	htmlBody.childs = append(htmlBody.childs, svgE)
	return buildBaseHtml(gss, htmlBody, todayIndicatorScript)
}

func RenderGantt(conn *sql.DB, renderTargetPath string, debug bool) {
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
		fileName := strings.ReplaceAll(group, " ", "-") + "-gantt.html"
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
		err = RenderGanttHTML(events, file, g, debug)
		if err != nil {
			log.Printf("Failed to process %s:\n\t%v", fileName, err)
		}
	}
}

func RenderGanttHtml(events []models.GanttEvent, g ganttRenderMetadata) string {
	return h(buildGanttElmTree(events, g))
}

func RenderGanttV2(conn *sql.DB, groupName string) (string, error) {
	events, err := dbop.GetGanttGroupEvents(conn, groupName, true)
	if err != nil {
		return "", err
	}
	if len(events) == 0 {
		return "", fmt.Errorf("no rows found for group %s", groupName)
	}
	g, err := getGanttRenderMetadata(conn, groupName)
	if err != nil {
		return "", fmt.Errorf("Failed to get gantt render metadata for group %s: %v", groupName, err)
	}
	htmlString := RenderGanttHtml(events, g)
	return htmlString, nil
}
