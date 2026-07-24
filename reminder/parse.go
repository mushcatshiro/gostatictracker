package reminder

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mushcatshiro/gostatictracker/models"
)

func validate(seg, segName string, start, end int) (string, error) {
	var s string
	if strings.Contains(seg, "*/") {
		segs := strings.Split(seg, "*/")
		if len(segs) != 2 {
			return "", fmt.Errorf("unexpected value %s", seg)
		}
		s = segs[1]
	} else {
		s = seg
	}
	parsed, err := strconv.Atoi(s)
	if err != nil {
		return "", fmt.Errorf("fail to parse %s segment %s", segName, seg)
	}
	if parsed < start || parsed >= end {
		return "", fmt.Errorf("value %d not in [%d - %d)", parsed, start, end)
	}
	return seg, nil
}

// YYYY-MM-DDTHH:MM:SS+00:00
func ParseInput(inp string) (models.Schedule, error) {
	// support
	// every 28 day, month, year
	// on every 28th, month, year
	// specific date
	// requires start date, matters to every 28th day
	var s models.Schedule
	segments := strings.Split(inp, " ")
	if len(segments) != 5 {
		return s, errors.New("schedule must be space separated 5 parts ints * * * * *")
	}
	miSeg, err := validate(segments[0], "minute", 0, 60)
	if err != nil {
		return s, err
	}
	s.Minute = miSeg
	hSeg, err := validate(segments[1], "hour", 0, 60)
	if err != nil {
		return s, err
	}
	s.Hour = hSeg
	domSeg, err := validate(segments[2], "day of month", 1, 32)
	if err != nil {
		return s, err
	}
	s.DayOfMonth = domSeg
	moSeg, err := validate(segments[3], "month", 1, 12)
	if err != nil {
		return s, err
	}
	s.Month = moSeg
	dowSeg, err := validate(segments[4], "month", 0, 7)
	if err != nil {
		return s, err
	}
	s.Month = dowSeg
	return s, nil
}
