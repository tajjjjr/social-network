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

	createGroupsTableSQL := `CREATE TABLE Groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		creator_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		avatar TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createGroupMembersTableSQL := `CREATE TABLE Group_Members (
		group_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		is_accepted BOOLEAN DEFAULT 0,
		PRIMARY KEY (group_id, user_id)
	);`

	_, err = db.Exec(createGroupsTableSQL)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(createGroupMembersTableSQL)
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
	}

	createdGroup, err := store.CreateGroup(group)
	if err != nil {
		t.Fatalf("CreateGroup failed: %v", err)
	}

	if createdGroup.ID == 0 {
		t.Error("Expected created group to have an ID")
	}

	var title string
	err = db.QueryRow("SELECT title FROM Groups WHERE id = ?", createdGroup.ID).Scan(&title)
	if err != nil {
		t.Fatalf("Failed to query created group: %v", err)
	}

	if title != "Test Group" {
		t.Errorf("Expected group title to be 'Test Group', got '%s'", title)
	}
}
