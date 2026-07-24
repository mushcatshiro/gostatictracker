package dbop

import (
	"fmt"
)

func (db *DB) InsertUser(uid, email, ipaddress, role string) error {
	qInsertUser, err := db.getSql("insert-user.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(qInsertUser, uid, email, ipaddress, role)
	if err != nil {
		return fmt.Errorf("%v: failed to insert user %s", err, email)
	}
	return nil
}
