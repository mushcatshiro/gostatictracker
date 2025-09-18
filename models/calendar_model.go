package models

import "time"

type CalendarEvent struct {
	Title string
	Group string
	Start time.Time
	End   time.Time
}

type MonthGroup struct {
	FirstDayOfMonth time.Time
	Events          []CalendarEvent
}
