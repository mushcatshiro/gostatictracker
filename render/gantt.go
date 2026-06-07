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

func processRectAndLine(spacing, rowIdx, rowRectWidth int, g models.GanttRenderMetadata) (int, int, int, int) {
	rectX := spacing * (g.BaseHeaderRectWidth + g.HeaderRectMargin)
	rectY := (rowIdx * (g.RowRectHeight + g.RowRectMargin)) + g.HeadersOffset
	lineX1 := rectX + rowRectWidth
	lineX2 := rectX + rowRectWidth
	return rectX, rectY, lineX1, lineX2
}

func processOverFlowUnits(maxWidth, rowEndWidth int, g models.GanttRenderMetadata) int {
	// visually still blocking last few chars due to padding 0.5rem
	diff := maxWidth - rowEndWidth
	var overflowUnits int
	if diff <= 0 {
		return overflowUnits
	}

	for diff > 0 {
		diff = diff - (g.BaseHeaderRectWidth + g.HeaderRectMargin)
		overflowUnits++
	}

	return overflowUnits
}

func getGanttRows(events []models.GanttEvent, g models.GanttRenderMetadata, debug bool) ([]eventGanttRow, int) {
	// TODO
	// support actual start, actual end rendering including start, stop button
	rSlice := []eventGanttRow{}

	var maxWidth, rowEndWidth int

	lineY1 := g.HeadersOffset - 2
	lineY2 := g.HeadersOffset + len(events)*(g.RowRectHeight+g.RowRectMargin)

	for idx, event := range events {
		iStartTime, iEndTime := formatEventTimes(event.Start, event.End, g.IsDayView)
		if !g.IsDayView && idx == 0 {
			g.GroupStartTime = &iStartTime
		}
		duration := iStartTime.Sub(*g.GroupStartTime)
		daysSpacing := int(duration.Hours() / float64(g.Divisor))

		titleWidth := getTextEstimateWidth(event.Title)
		rowRectWidth := getRowRectWidth(iStartTime, iEndTime, g.BaseHeaderRectWidth+g.HeaderRectMargin, float64(g.Divisor))

		var textX int
		var textY float32

		rectX, rectY, lineX1, lineX2 := processRectAndLine(daysSpacing, idx, rowRectWidth, g)

		// keeping text here since `maxWidth` is coupled with `textX`
		textY = float32(rectY) + g.TextYOffset
		rowEndWidth = max(rowEndWidth, rectX+rowRectWidth)
		if rowRectWidth-g.RowRectMargin-(2*g.RectToTextMargin) > titleWidth {
			textX = rectX + g.RectToTextMargin
			maxWidth = max(maxWidth, rectX+rowRectWidth)
		} else {
			textX = rectX + rowRectWidth + g.RectToTextMargin
			maxWidth = max(maxWidth, rectX+rowRectWidth+(g.RectToTextMargin*2)+titleWidth)
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

func getGanttHeaders(g models.GanttRenderMetadata, overflowUnits int) ([]eventGanttHeader, []eventGanttHeader, []eventGanttHeader) {
	startTime := *g.GroupStartTime
	endTime := *g.GroupEndTime

	var overflowDays int
	if g.IsDayView {
		overflowDays = overflowUnits
	} else {
		for overflowDays < overflowUnits*7 {
			overflowDays = overflowDays + 7
		}

	}
	endTime = endTime.AddDate(0, 0, overflowDays)
	if !g.IsDayView {
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
		if g.IsDayView {
			d = eventGanttHeader{
				RectX:     dIdx * (g.BaseHeaderRectWidth + g.HeaderRectMargin),
				RectWidth: g.BaseHeaderRectWidth,
				TextX:     dIdx*(g.BaseHeaderRectWidth+g.HeaderRectMargin) + int(g.BaseHeaderRectWidth/2),
				TextY:     2*(g.HeaderRectHeight+g.HeaderRectMargin) + int(g.HeaderRectHeight/2) + g.HeaderTextYOffset,
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
					RectX:     dIdx * (g.BaseHeaderRectWidth + g.HeaderRectMargin),
					RectWidth: g.BaseHeaderRectWidth,
					TextX:     dIdx*(g.BaseHeaderRectWidth+g.HeaderRectMargin) + int(g.BaseHeaderRectWidth/2),
					TextY:     2*(g.HeaderRectHeight+g.HeaderRectMargin) + int(g.HeaderRectHeight/2) + g.HeaderTextYOffset,
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
			currentMonthWidth := (cumForMonth * (g.BaseHeaderRectWidth + g.HeaderRectMargin)) - g.HeaderRectMargin
			m := eventGanttHeader{
				RectX:     monthRectEndX,
				RectWidth: currentMonthWidth,
				TextX:     monthRectEndX + int(currentMonthWidth/2),
				TextY:     (g.HeaderRectHeight + g.HeaderRectMargin) + int(g.HeaderRectHeight/2) + g.HeaderTextYOffset,
				TextVal:   strconv.Itoa(int(iterMonth)),
			}
			monthGanttHeader = append(monthGanttHeader, m)
			monthRectEndX += currentMonthWidth + g.HeaderRectMargin
			cumForMonth = 0
		}
		if startTime.AddDate(0, 0, 1).Year() != startTime.Year() {
			currentYearWidth := (cumForYear * (g.BaseHeaderRectWidth + g.HeaderRectMargin)) - g.HeaderRectMargin
			y := eventGanttHeader{
				RectX:     yearRectEndX,
				RectWidth: currentYearWidth,
				TextX:     yearRectEndX + int(currentYearWidth/2),
				TextY:     int(g.HeaderRectHeight/2) + g.HeaderTextYOffset,
				TextVal:   strconv.Itoa(iterYear)}
			yearGanttHeader = append(yearGanttHeader, y)
			yearRectEndX += currentYearWidth + g.HeaderRectMargin
			cumForYear = 0
		}

		startTime = startTime.AddDate(0, 0, 1)
	}
	if cumForMonth > 0 {
		currentMonthWidth := (cumForMonth * (g.BaseHeaderRectWidth + g.HeaderRectMargin)) - g.HeaderRectMargin
		m := eventGanttHeader{
			RectX:     monthRectEndX,
			RectWidth: currentMonthWidth,
			TextX:     monthRectEndX + int(currentMonthWidth/2),
			TextY:     (g.HeaderRectHeight + g.HeaderRectMargin) + int(g.HeaderRectHeight/2) + g.HeaderTextYOffset,
			TextVal:   strconv.Itoa(int(iterMonth)),
		}
		monthGanttHeader = append(monthGanttHeader, m)
	}
	if cumForYear > 0 {
		currentYearWidth := (cumForYear * (g.BaseHeaderRectWidth + g.HeaderRectMargin)) - g.HeaderRectMargin
		y := eventGanttHeader{
			RectX:     yearRectEndX,
			RectWidth: currentYearWidth,
			TextX:     yearRectEndX + int(currentYearWidth/2),
			TextY:     int(g.HeaderRectHeight/2) + g.HeaderTextYOffset,
			TextVal:   strconv.Itoa(iterYear)}
		yearGanttHeader = append(yearGanttHeader, y)
	}
	return yearGanttHeader, monthGanttHeader, dayGanttHeader
}

func getGanttEventBase(events []models.GanttEvent, g models.GanttRenderMetadata, debug bool) eventGanttBase {
	// allow headerRectWidth, headerRectHeight, rowRectHeight to have default values
	e := eventGanttBase{
		HeaderRectWidth:  g.BaseHeaderRectWidth,
		HeaderRectHeight: g.HeaderRectHeight,
		RowRectHeight:    g.RowRectHeight,
		Group:            g.GroupName,
	}
	e.SvgHeight = g.HeadersOffset + len(events)*(e.RowRectHeight+g.RowRectMargin) + 2*8

	rSlice, overflowUnits := getGanttRows(events, g, debug)
	e.Rows = rSlice

	y, m, d := getGanttHeaders(g, overflowUnits)
	e.SvgWidth = len(d)*(e.HeaderRectWidth+g.HeaderRectMargin) + 17 // buffer

	e.Years = y
	e.Months = m
	e.Days = d
	return e
}

func RenderGanttHTML(events []models.GanttEvent, file *os.File, g models.GanttRenderMetadata, debug bool) error {
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

func buildGanttElmTree(events []models.GanttEvent, g models.GanttRenderMetadata) elm {
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
	headerElm := elm{tag: "h2", innerText: g.GroupName}
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

		g, err := dbop.GetGanttRenderMetadata(conn, group)
		if err != nil {
			log.Fatalf("Failed to get gantt render metadata for group %s: %v", group, err)
		}
		err = RenderGanttHTML(events, file, g, debug)
		if err != nil {
			log.Printf("Failed to process %s:\n\t%v", fileName, err)
		}
	}
}

func RenderGanttV2(events []models.GanttEvent, g models.GanttRenderMetadata, groupName string) (string, error) {
	if len(events) == 0 {
		return "", fmt.Errorf("no rows found for group %s", groupName)
	}
	htmlString := h(buildGanttElmTree(events, g))
	return htmlString, nil
}
