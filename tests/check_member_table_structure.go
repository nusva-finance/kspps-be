package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Connect to database
	db, err := connectDB()
	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		return
	}
	defer closeDB(db)

	// Check members table structure
	fmt.Println("📋 Checking members table structure:")

	rows, err := db.Raw(`
		SELECT column_name, data_type, character_maximum_length
		FROM information_schema.columns
		WHERE table_name = 'members'
		ORDER BY ordinal_position
	`).Rows()

	if err != nil {
		fmt.Printf("❌ Error getting table structure: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var columnName, dataType string
		var maxLength *int

		err := rows.Scan(&columnName, &dataType, &maxLength)
		if err != nil {
			fmt.Printf("❌ Error scanning row: %v\n", err)
			continue
		}

		lengthInfo := ""
		if maxLength != nil {
			lengthInfo = fmt.Sprintf("(%d)", *maxLength)
		}

		fmt.Printf("  %-20s %-20s %s\n", columnName, dataType, lengthInfo)
	}
}

func connectDB() (*gorm.DB, error) {
	host := "103.127.98.221"
	port := "1482"
	user := "userkoperasi"
	password := "nusva12345"
	dbName := "nusvakspps"

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func closeDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}
