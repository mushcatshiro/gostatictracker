package dbop

import (
	"database/sql"
	"time"
)

func UpdateInsertTime(db *sql.DB, title string, timestamp *time.Time) (int64, error) {
	stmt := `UPDATE events
	SET insertTime = $1
	WHERE id = (SELECT id FROM events WHERE title = $2)
	RETURNING id;`

	var id int64
	err := db.QueryRow(stmt,
		timestamp,
		title,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
