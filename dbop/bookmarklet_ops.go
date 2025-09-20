package dbop

import (
	"database/sql"
	"fmt"

	"github.com/mushcatshiro/gostatictracker/models"
)

func GetSpecificGroupEvents(db *sql.DB, groupName string) ([]models.Bookmarklet, error) {
	var bs []models.Bookmarklet
	readQuery := `SELECT title, description, url, inserttime FROM events WHERE "group" = $1`
	rows, err := db.Query(readQuery, groupName)
	if err != nil {
		return bs, fmt.Errorf("failed to get bookmarklet(s): %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var b models.Bookmarklet
		if err := rows.Scan(&b.Title, &b.Description, &b.URL, &b.InsertTime); err != nil {
			return bs, fmt.Errorf("failed to scan bookmarklet: %v", err)
		}
		bs = append(bs, b)
	}
	if err := rows.Err(); err != nil {
		return bs, fmt.Errorf("error reading events: %v", err)
	}
	return bs, nil
}
