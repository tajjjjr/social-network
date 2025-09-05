package store

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

func setupGroupPostTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Create tables
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

	createGroupPostsTableSQL := `CREATE TABLE Group_Posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT,
		image TEXT,
		like_count INTEGER DEFAULT 0,
		dislike_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createGroupCommentsTableSQL := `CREATE TABLE Group_Post_Comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		parent_comment_id INTEGER,
		content TEXT NOT NULL,
		image TEXT,
		like_count INTEGER DEFAULT 0,
		dislike_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createGroupPermissionsTableSQL := `CREATE TABLE Group_Permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		permission_type TEXT NOT NULL,
		granted_by INTEGER NOT NULL,
		granted_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	tables := []string{
		createUsersTableSQL,
		createGroupsTableSQL,
		createGroupPostsTableSQL,
		createGroupCommentsTableSQL,
		createGroupPermissionsTableSQL,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			t.Fatal(err)
		}
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO Users (id, email, password, firstname, lastname, nickname, avatar) VALUES (1, 'test@example.com', 'password', 'Test', 'User', 'testuser', 'avatar.jpg')")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO Groups (id, title, description, creator_id) VALUES (1, 'Test Group', 'Test Description', 1)")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestCreateGroupPost(t *testing.T) {
	db := setupGroupPostTestDB(t)
	defer db.Close()

	store := NewGroupPostStore(db)

	post := &models.GroupPost{
		GroupID: 1,
		UserID:  1,
		Content: "Test group post content",
		Image:   "test.jpg",
	}

	createdPost, err := store.CreateGroupPost(post)
	if err != nil {
		t.Fatalf("CreateGroupPost failed: %v", err)
	}

	if createdPost.ID == 0 {
		t.Error("Expected post ID to be set")
	}

	if createdPost.Content != "Test group post content" {
		t.Errorf("Expected content 'Test group post content', got %s", createdPost.Content)
	}
}

func TestGetGroupPosts(t *testing.T) {
	db := setupGroupPostTestDB(t)
	defer db.Close()

	store := NewGroupPostStore(db)

	// Create test posts
	post1 := &models.GroupPost{GroupID: 1, UserID: 1, Content: "Post 1"}
	post2 := &models.GroupPost{GroupID: 1, UserID: 1, Content: "Post 2"}

	_, err := store.CreateGroupPost(post1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.CreateGroupPost(post2)
	if err != nil {
		t.Fatal(err)
	}

	posts, err := store.GetGroupPosts(1, 1, 10, 0)
	if err != nil {
		t.Fatalf("GetGroupPosts failed: %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}
}

func TestCreateGroupPostComment(t *testing.T) {
	db := setupGroupPostTestDB(t)
	defer db.Close()

	store := NewGroupPostStore(db)

	// Create a post first
	post := &models.GroupPost{GroupID: 1, UserID: 1, Content: "Test post"}
	createdPost, err := store.CreateGroupPost(post)
	if err != nil {
		t.Fatal(err)
	}

	comment := &models.GroupPostComment{
		GroupPostID: createdPost.ID,
		UserID:      1,
		Content:     "Test comment",
	}

	createdComment, err := store.CreateGroupPostComment(comment)
	if err != nil {
		t.Fatalf("CreateGroupPostComment failed: %v", err)
	}

	if createdComment.ID == 0 {
		t.Error("Expected comment ID to be set")
	}

	if createdComment.Content != "Test comment" {
		t.Errorf("Expected content 'Test comment', got %s", createdComment.Content)
	}
}

func TestCanUserDeleteGroupContent(t *testing.T) {
	db := setupGroupPostTestDB(t)
	defer db.Close()

	store := NewGroupPostStore(db)

	// Test group creator can delete
	canDelete, err := store.CanUserDeleteGroupContent(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !canDelete {
		t.Error("Group creator should be able to delete content")
	}

	// Test non-creator cannot delete
	_, err = db.Exec("INSERT INTO Users (id, email, password) VALUES (2, 'user2@example.com', 'password')")
	if err != nil {
		t.Fatal(err)
	}

	canDelete, err = store.CanUserDeleteGroupContent(1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if canDelete {
		t.Error("Non-creator should not be able to delete content")
	}

	// Test admin can delete
	_, err = db.Exec("INSERT INTO Group_Permissions (group_id, user_id, permission_type, granted_by) VALUES (1, 2, 'admin', 1)")
	if err != nil {
		t.Fatal(err)
	}

	canDelete, err = store.CanUserDeleteGroupContent(1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !canDelete {
		t.Error("Group admin should be able to delete content")
	}
}
