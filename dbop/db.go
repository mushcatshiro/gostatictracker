package dbop

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"github.com/mushcatshiro/gostatictracker/common"
)

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

func ConnectDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
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
