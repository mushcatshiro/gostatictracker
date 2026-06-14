package dbop

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"github.com/mushcatshiro/gostatictracker/common"
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

func ConnectDB(connStr, dbType string) (*sql.DB, error) {
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
