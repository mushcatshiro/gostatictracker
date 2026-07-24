package dbop

import (
	"database/sql"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/mushcatshiro/gostatictracker/mock"
)

var globalTestDB *DB

func TestMain(m *testing.M) {
	var err error
	globalTestDB, err = NewDB("sqlite", "", "", ".", "test", "")
	if err != nil {
		log.Fatalf("Connection to DB failed: %v", err)
	}
	err = globalTestDB.InitDB(false, true)
	if err != nil {
		log.Fatalf("Failed to initialize test database: %s", err)
	}

	mockData := slices.Concat(
		mock.DayViewMockData[:],
		mock.DayViewOverflowMockData[:],
		mock.WeekViewMockData[:],
	)
	for _, record := range mockData {
		_, err := globalTestDB.InsertRecord(record)
		if err != nil {
			log.Fatalf("Failed to insert mock event: %v", err)
		}
	}

	exitCode := m.Run()
	if globalTestDB != nil {
		globalTestDB.Close()
	}
	err = os.Remove("test.db")
	if err != nil {
		log.Printf("failed to delete test.db, got %v", err)
	}
	os.Exit(exitCode)
}

func SetupTestTx(t *testing.T) *DB {
	t.Helper()

	if globalTestDB == nil {
		t.Skip("Database not initialized, skipping test")
	}
	sqlDB, ok := globalTestDB.DBTX.(*sql.DB)
	if !ok {
		t.Fatalf("globalTestDB.DBTX is not a *sql.DB connection pool")
	}

	tx, err := sqlDB.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	t.Cleanup(func() {
		tx.Rollback()
	})
	return &DB{
		DBTX:       tx,
		queryAsset: globalTestDB.queryAsset,
		identifier: globalTestDB.identifier,
	}
}
