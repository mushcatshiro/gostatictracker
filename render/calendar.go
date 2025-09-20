package render

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
)

func buildCalendarElmTree(fdom time.Time) []elm {
	var ce []elm
	dayOfFdom := fdom.Weekday()
	for range int(dayOfFdom) {
		ce = append(ce, elm{
			tag:   "div",
			attrs: attrsStruct{class: "day-cell padding"},
		})
	}
	ldom := fdom.AddDate(0, 1, -1)
	for s := fdom; !s.After(ldom); s = s.AddDate(0, 0, 1) {
		ce = append(ce, elm{
			tag:       "div",
			attrs:     attrsStruct{class: "day-cell"},
			innerText: strconv.Itoa(s.Day()),
		})
	}
	return ce
}

func calculateStartRow(fdom, es time.Time) int {
	gridStartDate := fdom.AddDate(0, 0, -1*int(fdom.Weekday()))
	eventDay := es.Truncate(24 * time.Hour) // ignore time of day and prevent daylight saving issues
	gridStartDay := gridStartDate.Truncate(24 * time.Hour)
	duration := eventDay.Sub(gridStartDay)
	daysSinceGridStart := int(duration.Hours() / 24)
	row := (daysSinceGridStart / 7) + 1
	return row
}

func generateRowDiv(fdom, fdor time.Time, duration int, title string) elm {
	var e elm
	startRow := calculateStartRow(fdom, fdor)
	styleString := fmt.Sprintf(
		"--start-col: %s; --start-row: %s; --span: %s;",
		strconv.Itoa(int(fdor.Weekday())+1),
		strconv.Itoa(startRow),
		strconv.Itoa(duration),
	)
	e = elm{
		tag: "div",
		attrs: attrsStruct{
			class: "event-div",
			style: styleString,
		},
		innerText: title,
	}
	return e
}

func buildEventElmTree(fdom time.Time, ce []models.CalendarEvent) []elm {
	/*
		handle
			- start before month (css broken start)
			- continue after month (css broken end)
			- cross week (break into two divs)
		calculate start "grid-column-start"/"-end"
		calculate start row "grid-row-start: 2" (2nd week)
	*/
	var el []elm
	for _, e := range ce {
		// check if e spans into next row
		sdow := e.Start.Weekday()
		// remaining days in week
		rdiw := 7 - int(sdow)
		eventDuration := int(e.End.Sub(e.Start).Hours()/24) + 1
		if eventDuration > rdiw {
			// start segment
			el = append(el, generateRowDiv(fdom, e.Start, rdiw, e.Title + " " + e.Group))
			// middle segment
			eventDuration -= rdiw
			fdor := e.Start.AddDate(0, 0, rdiw)
			for eventDuration >= 7 {
				el = append(el, generateRowDiv(fdom, fdor, 7, e.Title + " " + e.Group))
				eventDuration -= 7
				fdor = fdor.AddDate(0, 0, 7)
			}
			// end segment
			if eventDuration > 0 {
				el = append(el, generateRowDiv(fdom, fdor, eventDuration, e.Title + " " + e.Group))
			}
		} else {
			el = append(el, generateRowDiv(fdom, e.Start, eventDuration, e.Title + " " + e.Group))
		}
	}
	return el
}

func renderCalendarHTML(mg models.MonthGroup, file *os.File) error {
	calendarElmTree := buildCalendarElmTree(mg.FirstDayOfMonth)
	eventElmTree := buildEventElmTree(mg.FirstDayOfMonth, mg.Events)

	cc := calendarContainer
	cw := calendarWeekdays
	cc.childs = append(cc.childs, cw...)
	cc.childs = append(cc.childs, calendarElmTree...)
	cc.childs = append(cc.childs, eventElmTree...)

	cb := calendarBody
	ct := calendarTitle
	ct.innerText = mg.FirstDayOfMonth.Format("January 2006")
	cb.childs = append(cb.childs, ct)
	cb.childs = append(cb.childs, cc)

	ch := calendarHeader
	ch.childs = append(ch.childs, cb)

	htmlString := h(ch)
	_, err := file.Write([]byte(htmlString))
	if err != nil {
		return err
	}
	return nil
}

func RenderCalendar(month, year int, renderTargetPath string, conn *sql.DB) {
	/*
		query by month of interest - then generate - color code by group;
		paradigm shift, to render by month(s) instead of group centric
	*/

	monthGroups, err := dbop.GetCalendarMonthGroups(conn, month, year)
	if err != nil {
		log.Fatalf("Not able to query any group(s): %v", err)
	}
	for _, mg := range monthGroups {
		fmt.Printf("%+v\n", mg)
		fileName := strconv.Itoa(int(mg.FirstDayOfMonth.Month())) +
			"-" + strconv.Itoa(mg.FirstDayOfMonth.Year()) +
			"-calendar.html"
		file, err := os.Create(renderTargetPath + "/" + fileName) // truncates if exists
		if err != nil {
			log.Printf("Failed to create file %s: %v", fileName, err)
			continue
		}
		err = renderCalendarHTML(mg, file)
		if err != nil {
			log.Fatalf("Failed to render for %d-%d", mg.FirstDayOfMonth.Month(), mg.FirstDayOfMonth.Year())
		}
	}

}
