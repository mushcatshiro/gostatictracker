package dbop

/*
func updateStatus(db DBTX, id int64, status common.Status) error {
	e, err := ReadEventById(db, id)
	if err != nil {
		return err
	}
	e.Status = status
	if status == common.COMPLETED {
		ae := time.Now()
		e.ActualEnd = &ae
	}
	err = UpdateEvent(db, e)
	if err != nil {
		return err
	}
	return nil
}

func UpdateStatusSafe(db *sql.DB, id int64, status common.Status) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// If tx.Commit() is successful later, this Rollback does nothing.
	// If the function returns early due to an error, this ensures the database
	// safely reverts the transaction and releases any locks.
	defer tx.Rollback()

	err = updateStatus(tx, id, status)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
*/
