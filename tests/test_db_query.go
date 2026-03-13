package main

import (
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "userkoperasi:nusva12345@tcp(103.127.98.221:1482)/nusvakspps?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}
	fmt.Println("✅ Connected to database successfully!")

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
	rows, err := db.Query("SELECT id, username, email, full_name, is_active FROM users")
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
}
