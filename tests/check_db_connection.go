package main

import (
	"database/sql"
	"fmt"
	
	_ "github.com/lib/pq"
)

func main() {
	// Test multiple possible database connections
	databases := []struct {
		name string
		connStr string
	}{
		{
			name: "Remote Database (from .env)",
			connStr: "host=103.127.98.221 port=1482 user=userkoperasi password=nusva12345 dbname=nusvakspps sslmode=disable",
		},
		{
			name: "Localhost Database",
			connStr: "host=localhost port=5432 user=postgres password= dbname=nusvakspps sslmode=disable",
		},
		{
			name: "Localhost with default user",
			connStr: "host=localhost port=5432 user=postgres password=postgres dbname=nusvakspps sslmode=disable",
		},
	}

	for _, dbInfo := range databases {
		fmt.Printf("\n🔍 Testing: %s\n", dbInfo.name)
		fmt.Printf("   Connection: %s\n", dbInfo.connStr)
		
		db, err := sql.Open("postgres", dbInfo.connStr)
		if err != nil {
			fmt.Printf("   ❌ Error opening connection: %v\n\n", err)
			continue
		}
		defer db.Close()

		// Test connection
		if err := db.Ping(); err != nil {
			fmt.Printf("   ❌ Failed to connect: %v\n\n", err)
			continue
		}

		fmt.Printf("   ✅ Connected successfully!\n")
		
		// Check users table
		var userCount int
		err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
		if err != nil {
			fmt.Printf("   ❌ Error counting users: %v\n\n", err)
			continue
		}
		
		fmt.Printf("   📊 Total users in database: %d\n", userCount)
		
		if userCount > 0 {
			fmt.Printf("   👥 Sample users:\n")
			rows, err := db.Query("SELECT id, username, email FROM users LIMIT 3")
			if err != nil {
				fmt.Printf("   ❌ Error querying users: %v\n\n", err)
				continue
			}
			defer rows.Close()
			
			for rows.Next() {
				var id int
				var username, email string
				err := rows.Scan(&id, &username, &email)
				if err != nil {
					fmt.Printf("   ❌ Error scanning row: %v\n", err)
					continue
				}
				fmt.Printf("      - ID: %d, Username: %s, Email: %s\n", id, username, email)
			}
		}
		
		fmt.Println()
		break
	}
}
