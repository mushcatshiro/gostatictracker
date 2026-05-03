package mock

import (
	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/models"
)

var DayViewMockData = [...]models.Event{
	{
		Start:       common.ParseStringDate("12-27-2024 00:00", false, false),
		End:         common.ParseStringDate("01-01-2025 00:00", false, false),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 1",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.COMPLETED,
		Priority:    common.DONOW,
	},
	{
		Start:       common.ParseStringDate("12-29-2024 00:00", false, false),
		End:         common.ParseStringDate("01-02-2025 00:00", false, false),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 2",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.COMPLETED,
		Priority:    common.DONOW,
	},
	{
		Start:       common.ParseStringDate("12-31-2024 00:00", false, false),
		End:         common.ParseStringDate("01-04-2025 00:00", false, false),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 3",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.INPROGRESS,
		Priority:    common.DOLATER,
	},
	{
		Start:       common.ParseStringDate("12-31-2024 00:00", false, false),
		End:         common.ParseStringDate("01-05-2025 00:00", false, false),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 4",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.INPROGRESS,
		Priority:    common.DOLATER,
	},
	{
		Start:       common.ParseStringDate("12-31-2024 00:00", false, false),
		End:         common.ParseStringDate("01-04-2025 00:00", false, false),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 5",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.CANCELLED,
		Priority:    common.ELIMINATE,
	},
	{
		Start:       common.ParseStringDate("12-30-2024 00:00", false, false),
		End:         common.ParseStringDate("01-05-2025 00:00", false, false),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 6",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.NOTSTARTED,
		Priority:    common.NOPRIORITY,
	},
	{
		Start:       common.ParseStringDate("12-28-2024 00:00", false, false),
		End:         common.ParseStringDate("01-05-2025 00:00", false, false),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 7",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.NOTSTARTED,
		Priority:    common.DELEGATE,
	},
}

var DayViewOverflowMockData = [...]models.Event{
	{
		Start:       common.ParseStringDate("01-04-2025 00:00", false, false),
		End:         common.ParseStringDate("01-05-2025 00:00", false, false),
		Group:       "day view overflow example",
		AllDay:      false,
		Title:       "Mock long long long long long long long long long task 1",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.CANCELLED,
		Priority:    common.DONOW,
	},
}

var WeekViewMockData = [...]models.Event{
	{
		Start:       common.ParseStringDate("12-22-2024 00:00", false, false),
		End:         common.ParseStringDate("01-03-2025 00:00", false, false),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 1",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.COMPLETED,
		Priority:    common.DOLATER,
	},
	{
		Start:       common.ParseStringDate("12-23-2024 00:00", false, false),
		End:         common.ParseStringDate("01-02-2025 00:00", false, false),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 2",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.COMPLETED,
		Priority:    common.DELEGATE,
	},
	{
		Start:       common.ParseStringDate("12-23-2024 00:00", false, false),
		End:         common.ParseStringDate("01-06-2025 00:00", false, false),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 3",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.INPROGRESS,
		Priority:    common.DONOW,
	},
	{
		Start:       common.ParseStringDate("12-23-2024 00:00", false, false),
		End:         common.ParseStringDate("12-30-2024 00:00", false, false),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 4",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.INPROGRESS,
		Priority:    common.DONOW,
	},
	{
		Start:       common.ParseStringDate("01-13-2025 00:00", false, false),
		End:         common.ParseStringDate("01-24-2025 00:00", false, false),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 5",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.NOTSTARTED,
		Priority:    common.ELIMINATE,
	},
	{
		Start:       common.ParseStringDate("01-30-2025 00:00", false, false),
		End:         common.ParseStringDate("02-05-2025 00:00", false, false),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 6",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.CANCELLED,
		Priority:    common.NOPRIORITY,
	},
}

func GetDayViewMockDataAs[T any](converter func(models.Event) T) []T {
	result := make([]T, len(DayViewMockData))
	for idx, item := range DayViewMockData {
		result[idx] = converter(item)
	}
	return result
}

func GetDayViewOverflowMockDataAs[T any](converter func(models.Event) T) []T {
	result := make([]T, len(DayViewOverflowMockData))
	for idx, item := range DayViewOverflowMockData {
		result[idx] = converter(item)
	}
	return result
}

func GetWeekViewMockDataAs[T any](converter func(models.Event) T) []T {
	result := make([]T, len(WeekViewMockData))
	for idx, item := range WeekViewMockData {
		result[idx] = converter(item)
	}
	return result
}
