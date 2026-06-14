package dbop

import (
	"database/sql"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/mushcatshiro/gostatictracker/mock"
)

var conn *sql.DB

func TestMain(m *testing.M) {
	testDbUrl := os.Getenv("TEST_DATABASE_URL")
	if testDbUrl == "" {
		log.Fatal("Environment varibale `TEST_DATABASE_URL cannot be empty\n")
	}
	testDbType := os.Getenv("TEST_DATABASE_TYPE")
	if testDbType == "" {
		log.Fatal("Environment varibale `TEST_DATABASE_TYPE cannot be empty\n")
	}

	var err error
	conn, err = ConnectDB(testDbUrl, testDbType)
	if err != nil {
		log.Fatalf("Connection to DB failed: %v", err)
	}
	err = InitDB(conn, false, true)
	if err != nil {
		log.Fatalf("Failed to initialize test database: %s", err)
	}

	mockData := slices.Concat(
		mock.DayViewMockData[:],
		mock.DayViewOverflowMockData[:],
		mock.WeekViewMockData[:],
	)
	for _, event := range mockData {
		_, err := InsertEvent(conn, event)
		if err != nil {
			log.Fatalf("Failed to insert mock event: %v", err)
		}
	}

	exitCode := m.Run()
	conn.Close()
	os.Exit(exitCode)
}

func SetupTestTx(t *testing.T) *sql.Tx {
	t.Helper()

	if conn == nil {
		t.Skip("Database not initialized, skipping test")
	}
	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	t.Cleanup(func() {
		tx.Rollback()
	})
	return tx
}
