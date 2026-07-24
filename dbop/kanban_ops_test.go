package dbop

import (
	"testing"

	"github.com/mushcatshiro/gostatictracker/common"
)

func TestGetKanbanGroups(t *testing.T) {
	db := SetupTestTx(t)

	result, err := db.GetKanbanGroups("day view example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	notStartedKey := common.NOTSTARTED.String()
	if len(result[notStartedKey]) != 2 {
		t.Errorf(
			"expected 2 event in %s, got %d",
			notStartedKey,
			len(result[notStartedKey]),
		)
	}

	completedKey := common.COMPLETED.String()
	if len(result[completedKey]) != 2 {
		t.Errorf(
			"expected 2 event in %s, got %d",
			completedKey,
			len(result[completedKey]),
		)
	}

	inProgressKey := common.INPROGRESS.String()
	if len(result[inProgressKey]) != 2 {
		t.Errorf(
			"expected 2 event in %s, got %d",
			inProgressKey,
			len(result[inProgressKey]),
		)
	}

	cancelledKey := common.CANCELLED.String()
	if len(result[cancelledKey]) != 1 {
		t.Errorf(
			"expected 1 event in %s, got %d",
			cancelledKey,
			len(result[cancelledKey]),
		)
	}
}
