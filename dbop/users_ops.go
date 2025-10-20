package dbop

import (
	"database/sql"
	"fmt"
)

func InsertUser(db *sql.DB, uid, email, ipaddress, role string) error {
	stmt := `INSERT INTO users (google_id, email, ipaddress, role)
	VALUES ($1, $2, $3, $4)`
	err := db.QueryRow(stmt, uid, email, ipaddress, role)
	if err != nil {
		return fmt.Errorf("failed to insert user %s", email)
	}
	return nil
}
