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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		creator_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createGroupMembersTableSQL := `CREATE TABLE Group_Members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		role TEXT DEFAULT 'member',
		is_accepted INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	tables := []string{createUsersTableSQL, createGroupsTableSQL, createGroupMembersTableSQL}

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

	_, err = db.Exec("INSERT INTO Groups (id, title, description, creator_id) VALUES (1, 'Test Group', 'Test Description', 1)")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO Group_Members (group_id, user_id, role, is_accepted) VALUES (1, 1, 'admin', 1)")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestGetGroupMembers(t *testing.T) {
	db := setupGroupMemberTestDB(t)
	defer db.Close()

	store := NewGroupMemberStore(db)

	members, err := store.GetGroupMembers(1)
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

	isMember, err := store.IsGroupMember(1, 1)
	if err != nil {
		t.Fatalf("IsGroupMember failed: %v", err)
	}

	if !isMember {
		t.Error("Expected user to be a group member")
	}

	isMember, err = store.IsGroupMember(1, 999)
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

	member, err := store.AddGroupMember(1, 2, "member")
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