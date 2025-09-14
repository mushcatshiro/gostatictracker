package render

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mushcatshiro/gostatictracker/dbop"
)

func buildCalendarElmTree(fdom time.Time) []elm {
	var ce []elm
	dayOfFdom := fdom.Weekday()
	for _ = range dayOfFdom {
		ce = append(ce, elm{
			tag:   "div",
			attrs: attrsStruct{class: "day-cell padding"},
		})
	}
	ldom := fdom.AddDate(0, 1, -1)
	for s := fdom; !s.After(ldom); s.AddDate(0, 0, 1) {
		ce = append(ce, elm{
			tag:       "div",
			attrs:     attrsStruct{class: "day-cell"},
			innerText: strconv.Itoa(s.Day()),
		})
	}
	return ce
}

func buildEventElmTree(mg dbop.MonthGroup) elm {
	/*
			handle
			- start before month (css broken start)
			- continue after month (css broken end)
			- cross week (break into two divs)
		calculate start "grid-column-start"/"-end"
		calculate start row "grid-row-start: 2" (2nd week)
	*/
	var e elm
	return e
}

func renderCalendarHTML(mg dbop.MonthGroup, file *os.File) error {
	calendarElmTree := buildCalendarElmTree(mg.FirstDayOfMonth)
	// eventElmTree := buildEventElmTree(mg)
	cc := calendarContainer
	cc.childs = append(cc.childs, calendarElmTree...)

	htmlString := h(cc)
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
