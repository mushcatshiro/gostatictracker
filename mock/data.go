package mock

import (
	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/models"
)

var DayViewMockData = [...]models.Record{
	{
		Start:       common.MustParseStringDate("2024-12-27T00:00", ""),
		End:         common.MustParseStringDate("2025-01-01T00:00", ""),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 1",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.COMPLETED,
		Priority:    common.DONOW,
	},
	{
		Start:       common.MustParseStringDate("2024-12-29T00:00", ""),
		End:         common.MustParseStringDate("2025-01-02T00:00", ""),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 2",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.COMPLETED,
		Priority:    common.DONOW,
	},
	{
		Start:       common.MustParseStringDate("2024-12-31T00:00", ""),
		End:         common.MustParseStringDate("2025-01-04T00:00", ""),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 3",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.INPROGRESS,
		Priority:    common.DOLATER,
	},
	{
		Start:       common.MustParseStringDate("2024-12-31T00:00", ""),
		End:         common.MustParseStringDate("2025-01-05T00:00", ""),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 4",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.INPROGRESS,
		Priority:    common.DOLATER,
	},
	{
		Start:       common.MustParseStringDate("2024-12-31T00:00", ""),
		End:         common.MustParseStringDate("2025-01-04T00:00", ""),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 5",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.CANCELLED,
		Priority:    common.ELIMINATE,
	},
	{
		Start:       common.MustParseStringDate("2024-12-30T00:00", ""),
		End:         common.MustParseStringDate("2025-01-05T00:00", ""),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 6",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.NOTSTARTED,
		Priority:    common.NOPRIORITY,
	},
	{
		Start:       common.MustParseStringDate("2024-12-28T00:00", ""),
		End:         common.MustParseStringDate("2025-01-05T00:00", ""),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 7",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.NOTSTARTED,
		Priority:    common.DELEGATE,
	},
}

var DayViewOverflowMockData = [...]models.Record{
	{
		Start:       common.MustParseStringDate("2025-01-04T00:00", ""),
		End:         common.MustParseStringDate("2025-01-05T00:00", ""),
		Group:       "day view overflow example",
		AllDay:      false,
		Title:       "Mock long long long long long long long long long task 1",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.CANCELLED,
		Priority:    common.DONOW,
	},
}

var WeekViewMockData = [...]models.Record{
	{
		Start:       common.MustParseStringDate("2024-12-22T00:00", ""),
		End:         common.MustParseStringDate("2025-01-03T00:00", ""),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 1",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.COMPLETED,
		Priority:    common.DOLATER,
	},
	{
		Start:       common.MustParseStringDate("2024-12-23T00:00", ""),
		End:         common.MustParseStringDate("2025-01-02T00:00", ""),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 2",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.COMPLETED,
		Priority:    common.DELEGATE,
	},
	{
		Start:       common.MustParseStringDate("2024-12-23T00:00", ""),
		End:         common.MustParseStringDate("2025-01-06T00:00", ""),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 3",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.INPROGRESS,
		Priority:    common.DONOW,
	},
	{
		Start:       common.MustParseStringDate("2024-12-23T00:00", ""),
		End:         common.MustParseStringDate("2024-12-30T00:00", ""),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 4",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.INPROGRESS,
		Priority:    common.DONOW,
	},
	{
		Start:       common.MustParseStringDate("2025-01-13T00:00", ""),
		End:         common.MustParseStringDate("2025-01-24T00:00", ""),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 5",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.NOTSTARTED,
		Priority:    common.ELIMINATE,
	},
	{
		Start:       common.MustParseStringDate("2025-01-30T00:00", ""),
		End:         common.MustParseStringDate("2025-02-05T00:00", ""),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 6",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
		Status:      common.CANCELLED,
		Priority:    common.NOPRIORITY,
	},
}

func GetDayViewMockDataAs[T any](converter func(models.Record) T) []T {
	result := make([]T, len(DayViewMockData))
	for idx, item := range DayViewMockData {
		result[idx] = converter(item)
	}
	return result
}

func GetDayViewOverflowMockDataAs[T any](converter func(models.Record) T) []T {
	result := make([]T, len(DayViewOverflowMockData))
	for idx, item := range DayViewOverflowMockData {
		result[idx] = converter(item)
	}
	return result
}

func GetWeekViewMockDataAs[T any](converter func(models.Record) T) []T {
	result := make([]T, len(WeekViewMockData))
	for idx, item := range WeekViewMockData {
		result[idx] = converter(item)
	}
	return result
}
