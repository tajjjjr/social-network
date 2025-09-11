package store

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupGroupMemberTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	createUsersTableSQL := `CREATE TABLE Users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		firstname TEXT,
		lastname TEXT,
		nickname TEXT,
		avatar TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

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

	tables := []string{createUsersTableSQL, createGroupsTableSQL}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			t.Fatal(err)
		}
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO Users (id, email, password, firstname, lastname, nickname) VALUES (1, 'test@example.com', 'password', 'Test', 'User', 'testuser')")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO Groups (id, public_id, type, user_id, title, content, role, privacy) VALUES ('1', 'group-1', 'group', 1, 'Test Group', 'Test Description', 'admin', 'public')")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO Groups (id, type, group_id, user_id, role, status) VALUES ('member-1', 'member', '1', 1, 'admin', 'active')")
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestGetGroupMembers(t *testing.T) {
	db := setupGroupMemberTestDB(t)
	defer db.Close()

	store := NewGroupMemberStore(db)

	members, err := store.GetGroupMembers("group-1")
	if err != nil {
		t.Fatalf("GetGroupMembers failed: %v", err)
	}

	if len(members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(members))
	}

	if *members[0].FirstName != "Test" {
		t.Errorf("Expected firstname 'Test', got %s", *members[0].FirstName)
	}
}

func TestIsGroupMember(t *testing.T) {
	db := setupGroupMemberTestDB(t)
	defer db.Close()

	store := NewGroupMemberStore(db)

	isMember, err := store.IsGroupMember("group-1", 1)
	if err != nil {
		t.Fatalf("IsGroupMember failed: %v", err)
	}

	if !isMember {
		t.Error("Expected user to be a group member")
	}

	isMember, err = store.IsGroupMember("group-1", 999)
	if err != nil {
		t.Fatalf("IsGroupMember failed: %v", err)
	}

	if isMember {
		t.Error("Expected user to not be a group member")
	}
}

func TestAddGroupMember(t *testing.T) {
	db := setupGroupMemberTestDB(t)
	defer db.Close()

	store := NewGroupMemberStore(db)

	// Add another user
	_, err := db.Exec("INSERT INTO Users (id, email, password, firstname, lastname) VALUES (2, 'test2@example.com', 'password', 'Test2', 'User2')")
	if err != nil {
		t.Fatal(err)
	}

	member, err := store.AddGroupMember("group-1", 2, "member")
	if err != nil {
		t.Fatalf("AddGroupMember failed: %v", err)
	}

	if member.UserID != 2 {
		t.Errorf("Expected UserID 2, got %d", member.UserID)
	}

	if member.Role != "member" {
		t.Errorf("Expected role 'member', got %s", member.Role)
	}
}
