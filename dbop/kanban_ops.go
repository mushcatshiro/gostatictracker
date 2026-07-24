package dbop

import (
	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/models"
)

func (db *DB) GetKanbanGroups(groupName string) (map[string][]models.Record, error) {
	out := make(map[string][]models.Record)

	for idx := common.NOTSTARTED; idx <= common.CANCELLED; idx++ {
		out[idx.String()] = []models.Record{}
	}

	allEvents, err := db.ReadRecords(models.Record{Group: groupName})
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
