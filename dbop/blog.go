package dbop

import "github.com/mushcatshiro/gostatictracker/models"

func (db *DB) ReadBlogRecords(f models.Record) ([]models.BlogRecord, error) {
	var bs []models.BlogRecord
	qSelectRecord, err := db.getSql("select-record.sql")
	if err != nil {
		return bs, err
	}
	q, args := AppendWhereClause(
		qSelectRecord, db.identifier, f,
	)
	rows, err := db.Query(q, args...)
	if err != nil {
		return bs, err
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
			return bs, err
		}
		bs = append(bs, r.ToBlogRecord())
	}
	if err := rows.Err(); err != nil {
		return bs, err
	}
	return bs, nil
}

func (db *DB) ReadBlogRecordByUrl(f models.Record) (models.BlogRecord, error) {
	var br models.BlogRecord
	qSelectRecord, err := db.getSql("select-blog-record.sql")
	if err != nil {
		return br, err
	}
	q, args := AppendWhereClause(
		qSelectRecord, db.identifier, f,
	)
	err = db.QueryRow(q, args...).Scan(
		&br.ID,
		&br.Title,
		&br.URL,
		&br.InsertTime,
		&br.Metadata,
	)
	return br, nil
}
