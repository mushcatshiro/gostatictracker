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

func buildEventElmTree(e []models.CalendarEvent) []elm {
	/*
		handle
			- start before month (css broken start)
			- continue after month (css broken end)
			- cross week (break into two divs)
		calculate start "grid-column-start"/"-end"
		calculate start row "grid-row-start: 2" (2nd week)
	*/
	var ce []elm
	return ce
}

func renderCalendarHTML(mg models.MonthGroup, file *os.File) error {
	calendarElmTree := buildCalendarElmTree(mg.FirstDayOfMonth)
	eventElmTree := buildEventElmTree(mg.Events)

	cc := calendarContainer
	cw := calendarWeekdays
	cc.childs = append(cc.childs, cw...)
	cc.childs = append(cc.childs, calendarElmTree...)

	cb := calendarBody
	ct := calendarTitle
	ct.innerText = mg.FirstDayOfMonth.Format("January 2006")
	cb.childs = append(cb.childs, ct)
	cb.childs = append(cb.childs, cc)
	cb.childs = append(cb.childs, eventElmTree...)

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
	}

}
