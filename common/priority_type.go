package common

import (
	"fmt"
	"strconv"
	"strings"
)

type Priority int

const (
	NOPRIORITY = iota
	ELIMINATE
	DELEGATE
	DOLATER
	DONOW
)

var priorityName = [...]string{
	"No Priority",
	"Eliminate",
	"Delegate",
	"Do Later",
	"Do Now",
}

func (p Priority) String() string {
	if p < 0 || int(p) > len(priorityName) {
		return "Unknown Priority"
	}
	return priorityName[p]
}

var priorityMap = make(map[string]Priority)

func init() {
	for i, name := range priorityName {
		priorityMap[name] = Priority(i)
	}
}

func ParsePriority(input string) (Priority, error) {
	for k, v := range priorityMap {
		if strings.EqualFold(input, k) {
			return v, nil
		}
	}
	if i, err := strconv.Atoi(input); err == nil {
		if i >= int(NOPRIORITY) && i <= int(DONOW) {
			return Priority(i), nil
		}
	}
	return -1, fmt.Errorf("invalid status: %q", input)
}
