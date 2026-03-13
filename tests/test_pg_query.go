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

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}
	fmt.Println("✅ Connected to PostgreSQL database successfully!")

	// Count total users
	var totalUsers int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		log.Fatal("Error counting users:", err)
	}
	fmt.Printf("📊 Total users in database: %d\n", totalUsers)

	// Get all user details
	fmt.Println("\n👥 User Details:")
	fmt.Println("==========================================")
	rows, err := db.Query("SELECT id, username, email, full_name, is_active FROM users ORDER BY id")
	if err != nil {
		log.Fatal("Error querying users:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username, email, fullName string
		var isActive bool
		err := rows.Scan(&id, &username, &email, &fullName, &isActive)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		}
		fmt.Printf("ID: %d, Username: %s, Email: %s, Name: %s, Active: %v\n", id, username, email, fullName, isActive)
	}
	
	// Check user_roles table
	fmt.Println("\n🔗 User Roles:")
	fmt.Println("==========================================")
	roleRows, err := db.Query("SELECT ur.user_id, u.username, r.name as role_name FROM user_roles ur JOIN users u ON ur.user_id = u.id JOIN roles r ON ur.role_id = r.id ORDER BY ur.user_id")
	if err != nil {
		log.Println("Warning: Could not query user_roles:", err)
	} else {
		defer roleRows.Close()
		for roleRows.Next() {
			var userID int
			var username, roleName string
			err := roleRows.Scan(&userID, &username, &roleName)
			if err != nil {
				log.Fatal("Error scanning role row:", err)
			}
			fmt.Printf("User ID: %d, Username: %s, Role: %s\n", userID, username, roleName)
		}
	}
}
