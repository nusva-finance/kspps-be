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

	fmt.Println("📋 Checking saving_transactions table...")

	// Check if table exists
	var tableName string
	err = db.Raw(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name = 'saving_transactions'
	`).Scan(&tableName).Error

	if err != nil {
		fmt.Printf("❌ Error checking table existence: %v\n", err)
		return
	}

	if tableName == "" {
		fmt.Println("❌ Table 'saving_transactions' does not exist in database!")
		return
	}

	fmt.Println("✅ Table 'saving_transactions' exists in database")

	// Get count of records
	var count int64
	err = db.Table("saving_transactions").Count(&count).Error
	if err != nil {
		fmt.Printf("❌ Error counting records: %v\n", err)
		return
	}

	fmt.Printf("📊 Total records in saving_transactions: %d\n", count)

	// Get all records
	var transactions []struct {
		ID              uint   `json:"id"`
		SavingAccountID uint   `json:"saving_account_id"`
		MemberID        uint   `json:"member_id"`
		MemberName      string `json:"member_name"`
		TransactionType string `json:"transaction_type"`
		Amount          int64  `json:"amount"`
		Description     string `json:"description"`
		TransactionDate string `json:"transaction_date"`
		CreatedAt       string `json:"created_at"`
		BalanceBefore   string `json:"balance_before"`
		BalanceAfter    string `json:"balance_after"`
	}

	err = db.Table("saving_transactions").
		Select("st.id, st.saving_account_id, st.transaction_type, st.amount, st.description, st.transaction_date, st.created_at, sa.member_id, m.full_name as member_name, st.balance_before, st.balance_after").
		Joins("JOIN saving_accounts sa ON st.saving_account_id = sa.id").
		Joins("JOIN members m ON sa.member_id = m.id").
		Find(&transactions).Error

	if err != nil {
		fmt.Printf("❌ Error fetching transactions: %v\n", err)
		return
	}

	if len(transactions) == 0 {
		fmt.Println("⚠️  No records found in saving_transactions table")
	} else {
		fmt.Println("\n📋 Saving Transactions in Database:")
		fmt.Println("=====================================")
		for i, t := range transactions {
			fmt.Printf("%d. ID: %d\n", i+1, t.ID)
			fmt.Printf("   Saving Account ID: %d\n", t.SavingAccountID)
			fmt.Printf("   Member ID: %d\n", t.MemberID)
			fmt.Printf("   Member Name: %s\n", t.MemberName)
			fmt.Printf("   Transaction Type: %s\n", t.TransactionType)
			fmt.Printf("   Amount: Rp %d\n", t.Amount)
			fmt.Printf("   Description: %s\n", t.Description)
			fmt.Printf("   Balance Before: %s\n", t.BalanceBefore)
			fmt.Printf("   Balance After: %s\n", t.BalanceAfter)
			fmt.Printf("   Transaction Date: %s\n", t.TransactionDate)
			fmt.Printf("   Created At: %s\n", t.CreatedAt)
			fmt.Println("-------------------------------------")
		}
	}

	// Also check if member exists
	fmt.Println("\n👤 Checking member data for Member ID 5...")
	var member struct {
		ID        uint   `json:"id"`
		MemberNo  string `json:"member_no"`
		FullName  string `json:"full_name"`
		KtpNo     string `json:"ktp_no"`
	}

	err = db.Table("members").
		Select("id, member_no, full_name, ktp_no").
		Where("id = ?", 5).
		First(&member).Error

	if err != nil {
		fmt.Printf("❌ Error finding member: %v\n", err)
		return
	}

	fmt.Println("✅ Member found:")
	fmt.Printf("   ID: %d\n", member.ID)
	fmt.Printf("   Member No: %s\n", member.MemberNo)
	fmt.Printf("   Full Name: %s\n", member.FullName)
	fmt.Printf("   KTP No: %s\n", member.KtpNo)
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
