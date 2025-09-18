package mock

import (
	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/models"
)

var DayViewMockData = [...]models.GanttEvent{
	{
		Start:       common.ParseStringDate("12-27-2024 00:00"),
		End:         common.ParseStringDate("01-01-2025 00:00"),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 1",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
	{
		Start:       common.ParseStringDate("12-29-2024 00:00"),
		End:         common.ParseStringDate("01-02-2025 00:00"),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 2",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
	{
		Start:       common.ParseStringDate("12-31-2024 00:00"),
		End:         common.ParseStringDate("01-04-2025 00:00"),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 3",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
	{
		Start:       common.ParseStringDate("12-31-2024 00:00"),
		End:         common.ParseStringDate("01-05-2025 00:00"),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 4",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
	{
		Start:       common.ParseStringDate("12-31-2024 00:00"),
		End:         common.ParseStringDate("01-04-2025 00:00"),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 5",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
	{
		Start:       common.ParseStringDate("12-30-2024 00:00"),
		End:         common.ParseStringDate("01-05-2025 00:00"),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 6",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
	{
		Start:       common.ParseStringDate("12-28-2024 00:00"),
		End:         common.ParseStringDate("01-05-2025 00:00"),
		Group:       "day view example",
		AllDay:      false,
		Title:       "Mock task 7",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
}

var DayViewOverflowMockData = [...]models.GanttEvent{
	{
		Start:       common.ParseStringDate("01-04-2025 00:00"),
		End:         common.ParseStringDate("01-05-2025 00:00"),
		Group:       "day view overflow example",
		AllDay:      false,
		Title:       "Mock task 1",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
}

var WeekViewMockData = [...]models.GanttEvent{
	{
		Start:       common.ParseStringDate("12-22-2024 00:00"),
		End:         common.ParseStringDate("01-03-2025 00:00"),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 1",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
	{
		Start:       common.ParseStringDate("12-23-2024 00:00"),
		End:         common.ParseStringDate("01-02-2025 00:00"),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 2",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
	{
		Start:       common.ParseStringDate("12-23-2024 00:00"),
		End:         common.ParseStringDate("01-06-2025 00:00"),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 3",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
	{
		Start:       common.ParseStringDate("12-23-2024 00:00"),
		End:         common.ParseStringDate("12-30-2024 00:00"),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 4",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
	{
		Start:       common.ParseStringDate("01-13-2025 00:00"),
		End:         common.ParseStringDate("01-24-2025 00:00"),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 5",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
	{
		Start:       common.ParseStringDate("01-30-2025 00:00"),
		End:         common.ParseStringDate("02-05-2025 00:00"),
		Group:       "week view example",
		AllDay:      false,
		Title:       "Mock task 6",
		URL:         "http://example.com/event1",
		Description: "This is a mock event for testing purposes.",
	},
}
