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

	createGroupsTableSQL := `CREATE TABLE Groups (
		id TEXT PRIMARY KEY,
		public_id TEXT UNIQUE,
		type TEXT NOT NULL CHECK (type IN ('group', 'member', 'request', 'event', 'post', 'comment')),
		group_id TEXT,
		user_id INTEGER,
		title TEXT,
		content TEXT,
		role TEXT DEFAULT 'member' CHECK (role IN ('admin', 'member')),
		status TEXT DEFAULT 'active' CHECK (status IN ('active', 'pending', 'rejected', 'going', 'not_going')),
		privacy TEXT DEFAULT 'public' CHECK (privacy IN ('public', 'private')),
		image TEXT,
		data TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
		FOREIGN KEY (group_id) REFERENCES Groups(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(createGroupsTableSQL); err != nil {
		t.Fatal(err)
	}

	// Insert a group first
	_, err = db.Exec("INSERT INTO Groups (id, public_id, type, user_id, title, content, privacy) VALUES ('1', 'group-1', 'group', 1, 'Test Group', 'Test Description', 'public')")
	if err != nil {
		t.Fatal(err)
	}

	// Insert an event
	_, err = db.Exec("INSERT INTO Groups (id, public_id, type, group_id, user_id, title, content, data) VALUES ('event-1', 'event-pub-1', 'event', '1', 1, 'Test Event', 'Test Description', '2024-12-31T18:00:00Z')")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestGetGroupEvents(t *testing.T) {
	db := setupGroupEventTestDB(t)
	defer db.Close()

	store := NewGroupEventStore(db)

	events, err := store.GetGroupEvents("group-1")
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