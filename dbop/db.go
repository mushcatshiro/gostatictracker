package dbop

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Status int

const (
	NOTSTARTED Status = iota
	INPROGRESS
	COMPLETED
	CANCELLED
)

var statusName = [...]string{
	"NOTSTARTED",
	"INPROGRESS",
	"COMPLETED",
	"CANCELLED",
}

func (s Status) String() string {
	if s < 0 || int(s) > len(statusName) {
		return "UNKNOWN"
	}
	return statusName[s]
}

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
