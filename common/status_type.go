package common

import (
	"fmt"
	"strconv"
	"strings"
)

type Status int

const (
	NOTSTARTED Status = iota
	INPROGRESS
	COMPLETED
	CANCELLED
)

var statusName = [...]string{
	"Not Started",
	"In Progress",
	"Completed",
	"Cancelled",
}

func (s Status) String() string {
	if s < 0 || int(s) > len(statusName) {
		return "Unknown Status"
	}
	return statusName[s]
}

var statusMap = make(map[string]Status)

func init() {
	for i, name := range statusName {
		statusMap[name] = Status(i)
	}
}

func ParseStatus(input string) (Status, error) {
	for k, v := range statusMap {
		if strings.EqualFold(input, k) {
			return v, nil
		}
	}
	if i, err := strconv.Atoi(input); err == nil {
		if i >= int(NOPRIORITY) && i <= int(DONOW) {
			return Status(i), nil
		}
	}
	return -1, fmt.Errorf("invalid status: %q", input)
}
