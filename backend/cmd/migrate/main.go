package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <database_path>")
	}

	dbPath := os.Args[1]
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	fmt.Println("Running migration to add public_id to Groups table...")

	// Check if public_id column exists
	var columnExists bool
	err = db.QueryRow(`
		SELECT COUNT(*) > 0 
		FROM pragma_table_info('Groups') 
		WHERE name = 'public_id'
	`).Scan(&columnExists)
	if err != nil {
		log.Fatal("Failed to check column existence:", err)
	}

	if !columnExists {
		// Add public_id column
		_, err = db.Exec("ALTER TABLE Groups ADD COLUMN public_id TEXT")
		if err != nil {
			log.Fatal("Failed to add public_id column:", err)
		}
		fmt.Println("Added public_id column to Groups table")
	}

	// Create unique index if it doesn't exist
	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_groups_public_id ON Groups(public_id)")
	if err != nil {
		log.Fatal("Failed to create unique index:", err)
	}
	fmt.Println("Created unique index on public_id")

	// Backfill public_id for existing records
	rows, err := db.Query("SELECT id FROM Groups WHERE public_id IS NULL OR public_id = ''")
	if err != nil {
		log.Fatal("Failed to query records without public_id:", err)
	}
	defer rows.Close()

	var recordsUpdated int
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			log.Printf("Failed to scan record: %v", err)
			continue
		}

		publicID := uuid.New().String()
		_, err = db.Exec("UPDATE Groups SET public_id = ? WHERE id = ?", publicID, id)
		if err != nil {
			log.Printf("Failed to update record %d: %v", id, err)
			continue
		}
		recordsUpdated++
	}

	fmt.Printf("Updated %d records with new public_id values\n", recordsUpdated)
	fmt.Println("Migration completed successfully!")
}
