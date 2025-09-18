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
		log.Fatal("Environment varibale `TEST_DATABASE_URL is ''\n")
	}
	conn, err := ConnectDB(testDbUrl)
	if err != nil {
		log.Fatalf("Connection to DB failed: %v", err)
	}
	err = InitDB(conn)
	if err != nil {
		log.Fatalf("Failed to initialize test database: %s", err)
	}
	// insert Mock
	mockData := slices.Concat(
		mock.DayViewMockData[:],
		mock.DayViewOverflowMockData[:],
		mock.WeekViewMockData[:],
	)
	for _, event := range mockData {
		_, err := InsertEvent(conn, event.ToEvent())
		if err != nil {
			log.Fatalf("Failed to insert mock event: %v", err)
		}
	}
	exitCode := m.Run()
	conn.Close()
	os.Exit(exitCode)
}
