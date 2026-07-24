package common

import (
	"errors"
	"fmt"
	"time"
)

const DefaultTimeZone = "Asia/Singapore"
const TimeLayout = "01-02-2006 15:04" // MM-DD-YYYY hh:mm
const Html5TimeLayout = "2006-01-02T15:04"
const BlogTimeLayout = "2006-01-02T15:04:05-07:00"
const BlogTimeLayoutNoTtz = "2006-01-02 15:04"

func ParseStringDate(value, timezone string) (time.Time, error) {
	var parsed time.Time
	if value == "" {
		return parsed, nil
	}
	var tz string
	if timezone == "" {
		tz = DefaultTimeZone
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return parsed, nil
	}
	parsed, err = time.ParseInLocation(Html5TimeLayout, value, loc)
	if err != nil {
		return parsed, err
	}
	return parsed.UTC(), nil
}

func MustParseStringDate(value, timezone string) *time.Time {
	parsed, err := ParseStringDate(value, timezone)
	if err != nil {
		panic(err)
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

func ParseStringTimeToDifferentFormat(value string, iFormat, eFormat string) (string, error) {
	if value == "" {
		return "", errors.New("empty string given")
	}
	parsed, err := time.Parse(iFormat, value)
	if err != nil {
		return "", fmt.Errorf("fail to parse input format %s to time.Time, %v", iFormat, err)
	}
	return parsed.Format(eFormat), nil
}
