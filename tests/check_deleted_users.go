package main

import (
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=103.127.98.221 port=1482 user=userkoperasi password=nusva12345 dbname=nusvakspps sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Check if deleted_at column exists
	var hasColumn bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'deleted_at')").Scan(&hasColumn)
	if err != nil {
		log.Fatal("Error checking deleted_at column:", err)
	}

	if hasColumn {
		fmt.Println("✅ deleted_at column exists in users table")
		
		// Count non-deleted users (deleted_at is NULL)
		var activeUsers int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&activeUsers)
		if err != nil {
			log.Fatal("Error counting active users:", err)
		}
		fmt.Printf("📊 Active users (deleted_at IS NULL): %d\n", activeUsers)

		// Count deleted users (deleted_at is NOT NULL)
		var deletedUsers int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NOT NULL").Scan(&deletedUsers)
		if err != nil {
			log.Fatal("Error counting deleted users:", err)
		}
		fmt.Printf("🗑️  Deleted users (deleted_at IS NOT NULL): %d\n", deletedUsers)

		// Get all users with deleted_at status
		fmt.Println("\n👥 All Users with Status:")
		fmt.Println("==========================================")
		rows, err := db.Query("SELECT id, username, email, full_name, is_active, deleted_at FROM users ORDER BY id")
		if err != nil {
			log.Fatal("Error querying users:", err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var username, email, fullName string
			var isActive bool
			var deletedAt sql.NullTime
			err := rows.Scan(&id, &username, &email, &fullName, &isActive, &deletedAt)
			if err != nil {
				log.Fatal("Error scanning row:", err)
			}
			status := "Active"
			if deletedAt.Valid {
				status = fmt.Sprintf("DELETED at %v", deletedAt.Time)
			}
			fmt.Printf("ID: %d, Username: %s, Email: %s, Name: %s, Active: %v, Status: %s\n", 
				id, username, email, fullName, isActive, status)
		}
	} else {
		fmt.Println("❌ deleted_at column does NOT exist in users table")
		fmt.Println("📊 Total users:", getRowCount(db, "users"))
	}
}

func getRowCount(db *sql.DB, tableName string) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
	if err != nil {
		log.Fatal("Error counting rows in", tableName, ":", err)
	}
	return count
}
