package common

import (
	"log"
	"time"
)

const TimeLayout = "01-02-2006 15:04" // MM-DD-YYYY hh:mm

func ParseStringDate(value string, html5 bool) *time.Time {
	if value == "" {
		return nil
	}
	var parsed time.Time
	var err error
	if html5 {
		parsed, err = time.Parse("2006-01-02T15:04", value)
	} else {
		parsed, err = time.Parse(TimeLayout, value)
	}
	if err != nil {
		log.Printf("unexpected format %s results in %v", value, err)
		return nil
	}
	return &parsed
}

func ParseTimeToString(value *time.Time, html5 bool) string {
	if value != nil {
		if html5 {
			return value.Format("2006-01-02T15:04")
		}
		return value.Format(TimeLayout)
	}
	return ""
}
