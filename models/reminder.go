package models

import "time"

type Schedule struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
}

func (s *Schedule) Next() *time.Time {
	// need to handle if month dont of day of month
	return nil
}
