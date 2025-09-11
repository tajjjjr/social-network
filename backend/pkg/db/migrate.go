package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// Migrate applies SQL migrations and backfills data.
func Migrate(db *sql.DB) {
	// Apply SQL migrations
	migrations := []string{
		"0001_add_public_id_to_groups.up.sql",
		"0002_add_public_id_to_group_events.up.sql",
	}

	for _, migration := range migrations {
		p := filepath.Join("pkg/db/migrations", migration)
		fmt.Printf("Applying migration: %s\n", p)
		stmt, err := os.ReadFile(p)
		if err != nil {
			log.Fatalf("failed to read migration file %s: %v", p, err)
		}

		if _, err := db.Exec(string(stmt)); err != nil {
			log.Fatalf("failed to apply migration %s: %v", migration, err)
		}
	}

	// Backfill public_id for groups

	rows, err := db.Query("SELECT id FROM groups WHERE public_id IS NULL")
	if err != nil {
		log.Fatalf("failed to query groups: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("failed to scan group id: %v", err)
			continue
		}
		publicID := uuid.New().String()
		_, err := db.Exec("UPDATE groups SET public_id = ? WHERE id = ?", publicID, id)
		if err != nil {
			log.Printf("failed to update group %d: %v", id, err)
		}
	}

	// Backfill public_id for group_events

	rows, err = db.Query("SELECT id FROM group_events WHERE public_id IS NULL")
	if err != nil {
		log.Fatalf("failed to query group_events: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("failed to scan group_event id: %v", err)
			continue
		}
		publicID := uuid.New().String()
		_, err := db.Exec("UPDATE group_events SET public_id = ? WHERE id = ?", publicID, id)
		if err != nil {
			log.Printf("failed to update group_event %d: %v", id, err)
		}
	}

	fmt.Println("Migrations and backfill complete.")
}
