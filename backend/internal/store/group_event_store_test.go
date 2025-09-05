package store

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupGroupEventTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	createGroupEventsTableSQL := `CREATE TABLE Group_Events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		event_time DATETIME NOT NULL,
		created_by INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createGroupEventsTableSQL); err != nil {
		t.Fatal(err)
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO Group_Events (group_id, title, description, event_time, created_by) VALUES (1, 'Test Event', 'Test Description', '2024-12-31 18:00:00', 1)")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestGetGroupEvents(t *testing.T) {
	db := setupGroupEventTestDB(t)
	defer db.Close()

	store := NewGroupEventStore(db)

	events, err := store.GetGroupEvents(1)
	if err != nil {
		t.Fatalf("GetGroupEvents failed: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if events[0].Title != "Test Event" {
		t.Errorf("Expected title 'Test Event', got %s", events[0].Title)
	}
}