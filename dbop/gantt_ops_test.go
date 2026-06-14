package dbop

import (
	"fmt"
	"testing"
	"time"

	"github.com/mushcatshiro/gostatictracker/mock"
	"github.com/mushcatshiro/gostatictracker/models"
)

func compareTimePtrs(t1, t2 *time.Time) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}
	return t1.Equal(*t2) // Safely compares the actual time values
}

func verifyResult(g, r models.GanttRenderMetadata) error {
	if !compareTimePtrs(g.GroupStartTime, r.GroupStartTime) {
		return fmt.Errorf(
			"mismatched group start time: %v, expected %v",
			g.GroupStartTime, r.GroupStartTime,
		)
	}
	if !compareTimePtrs(g.GroupEndTime, r.GroupEndTime) {
		return fmt.Errorf(
			"mismatched group end time: %v, expected %v",
			g.GroupEndTime, r.GroupEndTime,
		)
	}
	if g.IsDayView != r.IsDayView {
		return fmt.Errorf(
			"mismatched is day view: %t, expected %t",
			g.IsDayView, r.IsDayView,
		)
	}
	if g.Divisor != r.Divisor {
		return fmt.Errorf(
			"mismatched divisor: %d, expected %d",
			g.Divisor, r.Divisor,
		)
	}
	return nil
}

func TestGetGanttRenderMetadata(t *testing.T) {
	tx := SetupTestTx(t)

	dveG, err := GetGanttRenderMetadata(tx, "day view example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dveMock := mock.DayViewMockData
	dveR := models.GanttRenderMetadata{
		GroupStartTime: dveMock[0].Start,
		GroupEndTime:   dveMock[len(dveMock)-1].End,
		IsDayView:      true,
		Divisor:        24,
	}

	err = verifyResult(dveG, dveR)
	if err != nil {
		t.Errorf(
			"Incorrect gantt render metadata for day view example, expected\n"+
				"  GroupStartTime %s GroupEndTime %s IsDayView %t Divisor %d; got \n"+
				"  GroupStartTime %s GroupEndTime %s IsDayView %t Divisor %d;\n"+
				"  %v",
			dveR.GroupStartTime, dveR.GroupEndTime, dveR.IsDayView, dveR.Divisor,
			dveG.GroupStartTime, dveG.GroupEndTime, dveG.IsDayView, dveG.Divisor,
			err,
		)
	}
	// day view overflow example
	// week view example
	wveG, err := GetGanttRenderMetadata(tx, "week view example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wveMock := mock.WeekViewMockData
	wveR := models.GanttRenderMetadata{
		GroupStartTime: wveMock[0].Start,
		GroupEndTime:   wveMock[len(wveMock)-1].End,
		IsDayView:      false,
		Divisor:        24 * 7,
	}

	err = verifyResult(wveG, wveR)
	if err != nil {
		t.Errorf(
			"Incorrect gantt render metadata for week view example, expected\n"+
				"  GroupStartTime %s GroupEndTime %s IsDayView %t Divisor %d; got \n"+
				"  GroupStartTime %s GroupEndTime %s IsDayView %t Divisor %d;\n"+
				"  %v",
			wveR.GroupStartTime, wveR.GroupEndTime, wveR.IsDayView, wveR.Divisor,
			wveG.GroupStartTime, wveG.GroupEndTime, wveG.IsDayView, wveG.Divisor,
			err,
		)
	}
}
