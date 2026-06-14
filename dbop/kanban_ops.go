package dbop

import (
	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/models"
)

func GetKanbanGroups(db DBTX, groupName string) (map[string][]models.Event, error) {
	out := make(map[string][]models.Event)

	for idx := common.NOTSTARTED; idx <= common.CANCELLED; idx++ {
		out[idx.String()] = []models.Event{}
	}

	fc := models.FilterCols{
		Group: groupName,
	}

	allEvents, err := ReadFilteredEvents(db, fc)
	if err != nil {
		return nil, err
	}

	for _, e := range allEvents {
		statusKey := common.Status(e.Status).String()

		if _, exists := out[statusKey]; exists {
			out[statusKey] = append(out[statusKey], e)
		}
	}

	return out, nil
}
