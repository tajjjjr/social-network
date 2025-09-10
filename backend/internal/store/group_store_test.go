package store

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE Groups (
		id TEXT PRIMARY KEY,
		public_id TEXT UNIQUE,
		type TEXT NOT NULL,
		group_id TEXT,
		user_id INTEGER,
		title TEXT,
		content TEXT,
		role TEXT DEFAULT 'member',
		status TEXT DEFAULT 'active',
		privacy TEXT DEFAULT 'public',
		image TEXT,
		data TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestCreateGroup(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewGroupStore(db)

	group := &models.Group{
		CreatorID:   1,
		Title:       "Test Group",
		Description: "This is a test group.",
		Privacy:     "public",
	}

	createdGroup, err := store.CreateGroup(group)
	if err != nil {
		t.Fatalf("CreateGroup failed: %v", err)
	}

	if createdGroup.ID == "" {
		t.Error("Expected created group to have an ID")
	}

	var title string
	err = db.QueryRow("SELECT title FROM Groups WHERE type = 'group' AND public_id = ?", createdGroup.PublicID).Scan(&title)
	if err != nil {
		t.Fatalf("Failed to query created group: %v", err)
	}

	if title != "Test Group" {
		t.Errorf("Expected group title to be 'Test Group', got '%s'", title)
	}
}
