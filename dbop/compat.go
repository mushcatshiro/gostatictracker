package dbop

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// AgnosticTime wraps time.Time to seamlessly handle both Postgres and SQLite
type AgnosticTime struct {
	time.Time
}

// Scan implements the sql.Scanner interface
func (a *AgnosticTime) Scan(value interface{}) error {
	if value == nil {
		return nil // Handle NULLs if they ever occur, though you guaranteed they won't
	}

	switch v := value.(type) {
	case time.Time:
		// POSTGRES: The driver successfully maps to native time.Time
		a.Time = v
		return nil

	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			a.Time = t
			return nil
		}

		t, err = time.Parse("2006-01-02 15:04:05 -0700 MST", v)
		if err == nil {
			a.Time = t
			return nil
		}

		t, err = time.Parse("2006-01-02 15:04:05", v)
		if err == nil {
			a.Time = t
			return nil
		}

		return fmt.Errorf("failed to parse string time %q", v)

	case []byte:
		// SQLITE/SQL: Sometimes standard library sql passes raw bytes
		t, err := time.Parse(time.RFC3339, string(v))
		if err != nil {
			return fmt.Errorf("failed to parse byte time: %w", err)
		}
		a.Time = t
		return nil

	default:
		return fmt.Errorf("cannot scan type %T into AgnosticTime", value)
	}
}

// Optional: Implement driver.Valuer if you want this struct to also handle INSERTs
func (a AgnosticTime) Value() (driver.Value, error) {
	return a.Time, nil
}
