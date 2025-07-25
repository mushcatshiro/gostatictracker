package render

import (
	"testing"
	"time"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/mock"

	"github.com/stretchr/testify/assert"
)

func mockGetGanttRenderMetadata(mockData []dbop.Event, isDayView bool) ganttRenderMetadata {
	groupStartTime, _ := time.Parse(TimeLayout, mockData[0].Start)
	groupEndTime, _ := time.Parse(TimeLayout, mockData[len(mockData)-1].End)
	var divisor int
	if isDayView {
		divisor = 24
	} else {
		divisor = 7 * 24
	}
	return ganttRenderMetadata{
		isDayView:            isDayView,
		groupStartTime:       groupStartTime,
		groupEndTime:         groupEndTime,
		groupName:            "test",
		rowTextInRectPadding: 4,
		rectToTextMargin:     2,
		textYOffset:          6.3,
		rowRectMargin:        2,
		rowRectHeight:        10,
		headersOffset:        (24 + 2) * 3,
		baseHeaderRectWidth:  30,
		headerRectHeight:     24,
		headerRectMargin:     2,
		headerTextYOffset:    2,
		divisor:              divisor,
	}
}

func TestGetTextEstimateWidth(t *testing.T) {
	assert.Equal(t, 4*7, getTextEstimateWidth("Test"))
	assert.Equal(t, 5*7, getTextEstimateWidth("Test "))
}

func TestGetRowRectWidth(t *testing.T) {
	t1, _ := time.Parse(TimeLayout, "12-27-2024")
	t2, _ := time.Parse(TimeLayout, "12-29-2024")
	t.Run("", func(t *testing.T) {
		result := getRowRectWidth(t1, t2, 26, float64(24))
		assert.Equal(t, 3*26, result)
	})

	t3, _ := time.Parse(TimeLayout, "12-27-2024")
	t4, _ := time.Parse(TimeLayout, "12-27-2024")
	t.Run("", func(t *testing.T) {
		result := getRowRectWidth(t3, t4, 26, float64(24))
		assert.Equal(t, 1*26, result)
	})
}

func TestGetGanttRows(t *testing.T) {
	// hardcoded `ganttRenderMetadata`
	g0 := mockGetGanttRenderMetadata(mock.DayViewMockData[:], true)

	// trick for converting a defined size array into a slice
	r0 := eventGanttRow{
		RectX:       0,
		RectY:       78,
		RectWidth:   192,
		DataDetails: "This is a mock event for testing purposes.\n12-27-2024\n01-01-2025",
		LineX1:      192,
		LineY1:      76,
		LineX2:      192,
		LineY2:      162,
		TextX:       96,
		TextY:       84.3,
		TextVal:     "Mock task 1",
	}
	r1 := eventGanttRow{
		RectX:       64,
		RectY:       90,
		RectWidth:   160,
		DataDetails: "This is a mock event for testing purposes.\n12-29-2024\n01-02-2025",
		LineX1:      224,
		LineY1:      76,
		LineX2:      224,
		LineY2:      162,
		TextX:       144,
		TextY:       96.3,
		TextVal:     "Mock task 2",
	}
	r2 := eventGanttRow{
		RectX:       128,
		RectY:       102,
		RectWidth:   160,
		DataDetails: "This is a mock event for testing purposes.\n12-31-2024\n01-04-2025",
		LineX1:      288,
		LineY1:      76,
		LineX2:      288,
		LineY2:      162,
		TextX:       208,
		TextY:       108.3,
		TextVal:     "Mock task 3",
	}
	r3 := eventGanttRow{
		RectX:       128,
		RectY:       114,
		RectWidth:   192,
		DataDetails: "This is a mock event for testing purposes.\n12-31-2024\n01-05-2025",
		LineX1:      320,
		LineY1:      76,
		LineX2:      320,
		LineY2:      162,
		TextX:       224,
		TextY:       120.3,
		TextVal:     "Mock task 4",
	}
	r4 := eventGanttRow{
		RectX:       128,
		RectY:       126,
		RectWidth:   160,
		DataDetails: "This is a mock event for testing purposes.\n12-31-2024\n01-04-2025",
		LineX1:      288,
		LineY1:      76,
		LineX2:      288,
		LineY2:      162,
		TextX:       208,
		TextY:       132.3,
		TextVal:     "Mock task 5",
	}
	r5 := eventGanttRow{
		RectX:       96,
		RectY:       138,
		RectWidth:   224,
		DataDetails: "This is a mock event for testing purposes.\n12-30-2024\n01-05-2025",
		LineX1:      320,
		LineY1:      76,
		LineX2:      320,
		LineY2:      162,
		TextX:       208,
		TextY:       144.3,
		TextVal:     "Mock task 6",
	}
	r6 := eventGanttRow{
		RectX:       32,
		RectY:       150,
		RectWidth:   288,
		DataDetails: "This is a mock event for testing purposes.\n12-28-2024\n01-05-2025",
		LineX1:      320,
		LineY1:      76,
		LineX2:      320,
		LineY2:      162,
		TextX:       176,
		TextY:       156.3,
		TextVal:     "Mock task 7",
	}

	t.Run("base case without overflow `getGanttEventRows`", func(t *testing.T) {
		resultGanttEventRow, resultOverflowDay, err := getGanttRows(mock.DayViewMockData[:], g0, true)
		assert.NoError(t, err, "Not expecting parsing related error: %v", err)
		assert.Equal(t, 0, resultOverflowDay)
		assert.Equal(t, 7, len(resultGanttEventRow))
		assert.Equal(t, r0, resultGanttEventRow[0])
		assert.Equal(t, r1, resultGanttEventRow[1])
		assert.Equal(t, r2, resultGanttEventRow[2])
		assert.Equal(t, r3, resultGanttEventRow[3])
		assert.Equal(t, r4, resultGanttEventRow[4])
		assert.Equal(t, r5, resultGanttEventRow[5])
		assert.Equal(t, r6, resultGanttEventRow[6])
	})

	g1 := mockGetGanttRenderMetadata(mock.DayViewOverflowMockData[:], true)

	r7 := eventGanttRow{
		RectX:       0,
		RectY:       78,
		RectWidth:   64,
		DataDetails: "This is a mock event for testing purposes.\n01-04-2025\n01-05-2025",
		LineX1:      64,
		LineY1:      76,
		LineX2:      64,
		LineY2:      90,
		TextX:       104,
		TextY:       84.3,
		TextVal:     "Mock task 1",
	}
	t.Run("base case without overflow `getGanttEventRows`", func(t *testing.T) {
		resultGanttEventRow, resultOverflowDay, err := getGanttRows(mock.DayViewOverflowMockData[:], g1, true)
		assert.NoError(t, err, "Not expecting parsing related error: %v", err)
		assert.Equal(t, 3, resultOverflowDay)
		assert.Equal(t, 1, len(resultGanttEventRow))
		assert.Equal(t, r7, resultGanttEventRow[0])
	})

}

func TestGetGanttHeaders(t *testing.T) {
	g0 := mockGetGanttRenderMetadata(mock.DayViewMockData[:], true)
	yh0 := eventGanttHeader{
		RectX:     0,
		RectWidth: 158,
		TextX:     79,
		TextY:     14,
		TextVal:   "2024",
	}
	mh0 := eventGanttHeader{
		RectX:     0,
		RectWidth: 158,
		TextX:     79,
		TextY:     40,
		TextVal:   "12",
	}
	dh0 := eventGanttHeader{
		RectX:     0,
		RectWidth: 30,
		TextX:     15,
		TextY:     66,
		TextVal:   "27",
		DataDate:  "2024-12-27",
	}

	t.Run("day gantt view", func(t *testing.T) {
		yh, mh, dh := getGanttHeaders(g0, 0)
		assert.Equal(t, 2, len(yh))
		assert.Equal(t, 2, len(mh))
		assert.Equal(t, 10, len(dh))
		assert.Equal(t, yh0, yh[0])
		assert.Equal(t, mh0, mh[0])
		assert.Equal(t, dh0, dh[0])
	})

	g1 := mockGetGanttRenderMetadata(mock.DayViewOverflowMockData[:], true)
	yh0 = eventGanttHeader{
		RectX:     0,
		RectWidth: 158,
		TextX:     79,
		TextY:     14,
		TextVal:   "2025",
	}
	mh0 = eventGanttHeader{
		RectX:     0,
		RectWidth: 158,
		TextX:     79,
		TextY:     40,
		TextVal:   "1",
	}
	dh0 = eventGanttHeader{
		RectX:     0,
		RectWidth: 30,
		TextX:     15,
		TextY:     66,
		TextVal:   "4",
		DataDate:  "2025-01-04",
	}

	t.Run("day gantt overflow view", func(t *testing.T) {
		yh, mh, dh := getGanttHeaders(g1, 3)
		assert.Equal(t, 1, len(yh))
		assert.Equal(t, 1, len(mh))
		assert.Equal(t, 5, len(dh))
		assert.Equal(t, yh0, yh[0])
		assert.Equal(t, mh0, mh[0])
		assert.Equal(t, dh0, dh[0])
	})

	g2 := mockGetGanttRenderMetadata(mock.WeekViewMockData[:], false)
	yh0 = eventGanttHeader{
		RectX:     0,
		RectWidth: 94,
		TextX:     47,
		TextY:     14,
		TextVal:   "2024",
	}
	mh0 = eventGanttHeader{
		RectX:     0,
		RectWidth: 94,
		TextX:     47,
		TextY:     40,
		TextVal:   "12",
	}
	dh0 = eventGanttHeader{
		RectX:     0,
		RectWidth: 30,
		TextX:     15,
		TextY:     66,
		TextVal:   "16",
		DataDate:  "2024-12-22",
	}
	dArr := [11]string{"16", "23", "1", "6", "13", "20", "27", "3", "10", "17", "24"}

	t.Run("week gantt view", func(t *testing.T) {
		yh, mh, dh := getGanttHeaders(g2, 3)
		assert.Equal(t, 2, len(yh))
		assert.Equal(t, 3, len(mh))
		assert.Equal(t, 11, len(dh))
		assert.Equal(t, yh0, yh[0])
		assert.Equal(t, mh0, mh[0])
		assert.Equal(t, dh0, dh[0])
		assert.Equal(t, dArr[1], dh[1].TextVal)
		assert.Equal(t, dArr[2], dh[2].TextVal)
		assert.Equal(t, dArr[3], dh[3].TextVal)
		assert.Equal(t, dArr[4], dh[4].TextVal)
		assert.Equal(t, dArr[5], dh[5].TextVal)
		assert.Equal(t, dArr[6], dh[6].TextVal)
		assert.Equal(t, dArr[7], dh[7].TextVal)
		assert.Equal(t, dArr[8], dh[8].TextVal)
		assert.Equal(t, dArr[9], dh[9].TextVal)
		assert.Equal(t, dArr[10], dh[10].TextVal)
	})

}
