package main

import (
	"fmt"

	"github.com/joho/godotenv"

	"nusvakspps/config"
	"nusvakspps/models"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Initialize database connection first
	if err := config.InitDB(); err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		return
	}

	// Get database connection
	db := config.GetDB()

	// Check if members table exists and has data
	var count int64
	db.Model(&models.Member{}).Count(&count)
	fmt.Printf("Total members in database: %d\n", count)

	if count > 0 {
		// Get first few members
		var members []models.Member
		db.Limit(5).Find(&members)

		fmt.Println("\nFirst 5 members:")
		for i, member := range members {
			fmt.Printf("%d. ID: %d, MemberNo: %s, FullName: %s, NIK: %s, IsActive: %v\n",
				i+1, member.ID, member.MemberNo, member.FullName, member.NIK, member.IsActive)
		}
	} else {
		fmt.Println("\n⚠️  Members table is empty!")
		fmt.Println("You need to create members through the UI or insert data directly.")
	}

	// Check table structure
	if !db.Migrator().HasTable(&models.Member{}) {
		fmt.Println("\n❌ Members table doesn't exist!")
	} else {
		fmt.Println("\n✅ Members table exists")
	}
}
