package common

import (
	"log"
	"time"
)

const TimeLayout = "01-02-2006 15:04" // MM-DD-YYYY hh:mm

func ParseStringDate(value string) *time.Time {
	if value == "" {
		return nil
	}
	parsed, err := time.Parse(TimeLayout, value)
	if err != nil {
		log.Fatalf("unexpected format %s results in %v", value, err)
		return nil
	}
	return &parsed
}
