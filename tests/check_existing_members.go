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

	fmt.Println("📋 Checking existing members...")

	// Get all members
	var members []struct {
		ID        uint   `json:"id"`
		MemberNo  string `json:"member_no"`
		FullName  string `json:"full_name"`
		KtpNo     string `json:"ktp_no"`
	}

	err = db.Model(&struct{}{}).
		Table("members").
		Select("id, member_no, full_name, ktp_no").
		Find(&members).Error

	if err != nil {
		fmt.Printf("❌ Error getting members: %v\n", err)
		return
	}

	fmt.Printf("Total members in database: %d\n", len(members))
	fmt.Println("\n📋 Existing Members:")
	for i, member := range members {
		fmt.Printf("%d. ID: %d, MemberNo: %s, FullName: %s, KTPNo: %s\n",
			i+1, member.ID, member.MemberNo, member.FullName, member.KtpNo)
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

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func closeDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}
