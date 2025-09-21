package dbop

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	NOTSTARTED = iota
	INPROGRESS
	COMPLETED
	CANCELLED
)

const (
	NOPRIORITY = iota
	ELIMINATE
	DELEGATE
	DOLATER
	DONOW
)

func GenerateConnStr(username, password , dbhost, dbname string) (string, error) {
	if username == "" || dbname == "" || password == "" {
		return "", fmt.Errorf("username, password, dbname must be provided")
	}
	if dbhost == "" {
		fmt.Printf("`dbhost` is not provided using localhost...")
		dbhost = "localhost" // Default to localhost if no host is provided
	}
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", username, password, dbhost, dbname)
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
