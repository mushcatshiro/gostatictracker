package dbop

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"github.com/mushcatshiro/gostatictracker/assets"
	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/models"
	_ "modernc.org/sqlite"
)

type DBTX interface {
	// Exec implicitly uses `context.Background()` and calls `ExecContext`
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row

	// Context variants (Best practice to include these for future-proofing)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type DB struct {
	DBTX
	queryAsset embed.FS
	identifier string
}

func (db *DB) Begin() (*sql.Tx, error) {
	sqlDB, ok := db.DBTX.(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("cannot begin transaction: underlying connection is not *sql.DB")
	}
	return sqlDB.Begin()
}

// Close safely attempts to close the database pool.
func (db *DB) Close() error {
	sqlDB, ok := db.DBTX.(*sql.DB)
	if !ok {
		return fmt.Errorf("cannot close connection: underlying connection is not *sql.DB")
	}
	return sqlDB.Close()
}

var ErrSqlNotFound = errors.New("file not found")

func GenerateConnStr(dbType, username, password, dbhost, dbname, sslmode string) (string, error) {
	var connStr string
	switch dbType {
	case "sqlite":
		if err := common.ValidateInputDirectory(dbhost); err != nil {
			return connStr, err
		}
		if !strings.HasSuffix(dbname, ".db") {
			dbname += ".db"
		}
		connStr = filepath.Join(dbhost, dbname)
	case "postgres":
		if username == "" || dbname == "" || password == "" {
			return "", fmt.Errorf("username, password, dbname must be provided for postgres database")
		}
		if dbhost == "" {
			fmt.Println("`dbhost` is not provided using localhost...")
			dbhost = "localhost" // Default to localhost if no host is provided
		}
		connStr = fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s", username, password, dbhost, dbname, sslmode)
	default:
		return connStr, fmt.Errorf("unexpected database type %s", dbType)
	}
	return connStr, nil
}

func CreateDBTX(connStr, dbType string) (DBTX, error) {
	db, err := sql.Open(dbType, connStr)
	if err != nil {
		return nil, err
	}
	// Check if the connection is established
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewDB(dbType, username, password, dbhost, dbname, sslmode string) (*DB, error) {
	var db DB
	connStr, err := GenerateConnStr(dbType, username, password, dbhost, dbname, sslmode)
	if err != nil {
		return &db, err
	}
	dbtx, err := CreateDBTX(connStr, dbType)
	if err != nil {
		return &db, err
	}
	var identifier string
	if dbType == "sqlite" {
		identifier = "?"
	} else {
		identifier = "$"
	}
	return &DB{
		DBTX:       dbtx,
		queryAsset: assets.SqlFS,
		identifier: identifier,
	}, nil
}

// consider this to ensure all queries are available
func (db *DB) verifyQueries() {}

func (db *DB) getSql(name string) (string, error) {
	file, err := db.queryAsset.Open(path.Join("sql", name))
	if err != nil {
		return "", fmt.Errorf("%s %w: %v", name, ErrSqlNotFound, err)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("not able to read %s: %v", name, err)
	}
	return string(content), nil
}

func AppendWhereClause(q, identifier string, f models.Record) (string, []any) {
	var conditions []string
	var args []any
	if f.ID > 0 {
		conditions = append(conditions, fmt.Sprintf(
			`"id" = %s%d`, identifier, len(conditions)+1))
		args = append(args, f.ID)
	}
	if f.Group != "" {
		conditions = append(conditions, fmt.Sprintf(
			`"group" = %s%d`, identifier, len(conditions)+1))
		args = append(args, f.Group)
	}
	if f.DefaultMode != "" {
		conditions = append(conditions, fmt.Sprintf(
			`"mode" = %s%d`, identifier, len(conditions)+1))
		args = append(args, f.DefaultMode)
	}
	if f.URL != "" {
		conditions = append(conditions, fmt.Sprintf(
			`"url" = %s%d`, identifier, len(conditions)+1))
		args = append(args, f.DefaultMode)
	}
	if f.Status > 0 && f.Status < 4 {
		conditions = append(conditions, fmt.Sprintf(
			`"status" = %s%d`, identifier, len(conditions)+1))
		args = append(args, f.Status)
	}

	var query string
	if len(conditions) > 0 {
		query = q + " WHERE " + strings.Join(conditions, " AND ")
	}

	return query, args
}
