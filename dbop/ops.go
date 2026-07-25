package dbop

import (
	"fmt"
	"time"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/models"
)

// TODO: create reserved group names in a separate table + their records esp
// urls
func (db *DB) InitDB(truncate, drop bool) error {
	// deletes all rows while keeping table structure intact
	qTruncTables, err := db.getSql("truncate-tables.sql")
	if err != nil {
		return err
	}
	qDropTables, err := db.getSql("drop-tables.sql")
	if err != nil {
		return err
	}
	qCreateRecordTable, err := db.getSql("create-record-table.sql")
	if err != nil {
		return err
	}
	qCreateUserTable, err := db.getSql("create-user-table.sql")
	if err != nil {
		return err
	}

	if truncate {
		_, err = db.Exec(qTruncTables)
		if err != nil {
			return fmt.Errorf("failed to truncate tables: %v", err)
		}
	}
	// removes the entire table definition with data, indexes and constraints
	if drop {
		_, err = db.Exec(qDropTables)
		if err != nil {
			return fmt.Errorf("failed to drop tables: %v", err)
		}
	}
	_, err = db.Exec(qCreateRecordTable)
	if err != nil {
		return fmt.Errorf("failed to create record table: %v", err)
	}
	_, err = db.Exec(qCreateUserTable)
	if err != nil {
		return fmt.Errorf("failed to create user table: %v", err)
	}
	return nil
}

func (db *DB) ReadRecord(f models.Record) (models.Record, error) {
	var r models.Record
	qSelectRecord, err := db.getSql("select-record.sql")
	if err != nil {
		return r, err
	}
	q, args := AppendWhereClause(
		qSelectRecord, db.identifier, models.Record{ID: f.ID},
	)

	err = db.QueryRow(q, args...).Scan(
		&r.ID,
		&r.Start,
		&r.End,
		&r.ActualStart,
		&r.ActualEnd,
		&r.InsertTime,
		&r.Group,
		&r.DefaultMode,
		&r.Repeat,
		&r.AllDay,
		&r.Title,
		&r.URL,
		&r.Description,
		&r.PID,
		&r.Priority,
		&r.Metadata,
		&r.Status,
	)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (db *DB) ReadRecords(f models.Record) ([]models.Record, error) {
	var rs []models.Record
	qSelectRecord, err := db.getSql("select-record.sql")
	if err != nil {
		return rs, err
	}
	q, args := AppendWhereClause(qSelectRecord, db.identifier, f)

	rows, err := db.Query(q, args...)
	if err != nil {
		return rs, err
	}
	defer rows.Close()
	for rows.Next() {
		var r models.Record
		if err := rows.Scan(
			&r.ID,
			&r.Start,
			&r.End,
			&r.ActualStart,
			&r.ActualEnd,
			&r.InsertTime,
			&r.Group,
			&r.DefaultMode,
			&r.Repeat,
			&r.AllDay,
			&r.Title,
			&r.URL,
			&r.Description,
			&r.PID,
			&r.Priority,
			&r.Metadata,
			&r.Status,
		); err != nil {
			return rs, err
		}
		rs = append(rs, r)
	}
	if err := rows.Err(); err != nil {
		return rs, err
	}
	return rs, nil
}

func (db *DB) ReadUniqueGroups() ([]string, error) {
	var gs []string
	qSelectUniqGroups, err := db.getSql("select-unique-groups.sql")
	if err != nil {
		return gs, err
	}
	rows, err := db.Query(qSelectUniqGroups)
	if err != nil {
		return gs, err
	}
	defer rows.Close()

	for rows.Next() {
		var group string
		if err := rows.Scan(&group); err != nil {
			return gs, err
		}
		gs = append(gs, group)
	}
	if err := rows.Err(); err != nil {
		return gs, err
	}
	return gs, nil
}

func (db *DB) InsertRecord(r models.Record) (int64, error) {
	var id int64
	qInsertRecord, err := db.getSql("insert-record.sql")
	if err != nil {
		return id, err
	}

	err = db.QueryRow(qInsertRecord,
		r.Start,
		r.End,
		r.ActualStart,
		r.ActualEnd,
		r.Group,
		r.DefaultMode,
		r.Repeat,
		r.AllDay,
		r.Title,
		r.URL,
		r.Description,
		r.PID,
		r.Priority,
		r.Metadata,
		r.Status,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db *DB) UpdateRecord(r models.Record) error {
	qUpdateRecord, err := db.getSql("update-record.sql")
	if err != nil {
		return err
	}

	result, err := db.Exec(qUpdateRecord,
		r.Start,
		r.End,
		r.ActualStart,
		r.ActualEnd,
		r.Group,
		r.AllDay,
		r.Title,
		r.URL,
		r.Description,
		r.PID,
		r.Priority,
		r.Metadata,
		r.Status,
		r.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no record found with ID %d", r.ID)
	}
	return nil
}

func (db *DB) UpdateRecordStatus(id int64, status common.Status) error {
	r, err := db.ReadRecord(models.Record{ID: id})
	if err != nil {
		return err
	}
	r.Status = status
	if status == common.COMPLETED {
		actlEnd := time.Now()
		r.ActualEnd = &actlEnd
	}
	return db.UpdateRecord(r)
}

func (db *DB) UpdateRecordPriority(id int64, priority common.Priority) error {
	r, err := db.ReadRecord(models.Record{ID: id})
	if err != nil {
		return err
	}
	r.Priority = priority
	return db.UpdateRecord(r)
}

func (db *DB) GetUngrouppedRecords() ([]models.Record, error) {
	return db.ReadRecords(models.Record{Group: "default"})
}
