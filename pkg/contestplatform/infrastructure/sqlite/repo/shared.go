package repo

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	_ "modernc.org/sqlite"
)

var idCounter atomic.Uint64

func OpenDatabase(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("create sqlite dir: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err = migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS problems (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			time_limit_ns INTEGER NOT NULL,
			memory_limit_bytes INTEGER NOT NULL,
			test_cases_json TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS submissions (
			id TEXT PRIMARY KEY,
			problem_id TEXT NOT NULL,
			language TEXT NOT NULL,
			source_code TEXT NOT NULL,
			verdict TEXT NOT NULL,
			test_results_json TEXT NOT NULL,
			created_at_unix INTEGER NOT NULL
		)`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("migrate sqlite database: %w", err)
		}
	}

	return nil
}

func newID(prefix string) string {
	sequence := idCounter.Add(1)
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), sequence)
}
