package dbop

import (
	"fmt"
)

func InsertUser(db DBTX, uid, email, ipaddress, role string) error {
	stmt := `INSERT INTO users (google_id, email, ipaddress, role)
	VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(stmt, uid, email, ipaddress, role)
	if err != nil {
		return fmt.Errorf("%v: failed to insert user %s", err, email)
	}
	return nil
}
