package common

import (
	"errors"
	"fmt"
	"log"
	"time"
)

const TimeLayout = "01-02-2006 15:04" // MM-DD-YYYY hh:mm
const Html5TimeLayout = "2006-01-02T15:04"
const BlogTimeLayout = "2006-01-02T15:04:05-07:00"
const BlogTimeLayoutNoTtz = "2006-01-02 15:04"

func ParseStringDate(value string, html5, blog bool) *time.Time {
	if value == "" {
		return nil
	}
	var parsed time.Time
	var err error
	if html5 {
		parsed, err = time.Parse(Html5TimeLayout, value)
	} else if blog {
		parsed, err = time.Parse(BlogTimeLayout, value)
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
