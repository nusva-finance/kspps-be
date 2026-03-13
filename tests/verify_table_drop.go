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

	fmt.Println("🔍 Verifying saving_transactions table...")
	fmt.Println("=====================================")

	// Check all tables containing 'saving'
	var tables []string
	err = db.Raw(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND (table_name LIKE '%saving%' OR table_name LIKE '%savings%')
		ORDER BY table_name
	`).Scan(&tables).Error

	if err != nil {
		fmt.Printf("❌ Error checking tables: %v\n", err)
		return
	}

	if len(tables) == 0 {
		fmt.Println("❌ No tables found containing 'saving' or 'savings'")
		return
	}

	fmt.Println("📋 Tables found:")
	for _, table := range tables {
		fmt.Printf("   - %s\n", table)
	}

	// Check specific tables
	var savingExists bool

	for _, table := range tables {
		if table == "saving_transactions" {
			savingExists = true
		}
	}

	fmt.Println("\n📊 Table Status:")
	fmt.Printf("   'saving_transactions' (should exist): ")
	if savingExists {
		fmt.Println("✅ EXISTS")
	} else {
		fmt.Println("❌ DOESN'T EXIST!")
	}

	// If saving_transactions exists, show its data
	if savingExists {
		var count int64
		err = db.Table("saving_transactions").Count(&count).Error
		if err != nil {
			fmt.Printf("❌ Error counting records: %v\n", err)
			return
		}

		fmt.Printf("\n📊 Records in 'saving_transactions': %d\n", count)

		if count > 0 {
			var transactions []struct {
				ID              uint   `json:"id"`
				SavingAccountID uint   `json:"saving_account_id"`
				MemberID        uint   `json:"member_id"`
				MemberName      string `json:"member_name"`
				TransactionType string `json:"transaction_type"`
				Amount          int64  `json:"amount"`
				Description     string `json:"description"`
				BalanceBefore   string `json:"balance_before"`
				BalanceAfter    string `json:"balance_after"`
				TransactionDate string `json:"transaction_date"`
				CreatedAt       string `json:"created_at"`
			}

			err = db.Table("saving_transactions").
				Select("st.id, st.saving_account_id, st.transaction_type, st.amount, st.description, st.balance_before, st.balance_after, st.transaction_date, st.created_at, sa.member_id, m.full_name as member_name").
				Joins("JOIN saving_accounts sa ON st.saving_account_id = sa.id").
				Joins("JOIN members m ON sa.member_id = m.id").
				Order("st.created_at DESC").
				Find(&transactions).Error

			if err != nil {
				fmt.Printf("❌ Error fetching transactions: %v\n", err)
				return
			}

			fmt.Println("\n📋 Latest Transactions:")
			fmt.Println("=====================================")
			for i, t := range transactions {
				fmt.Printf("%d. ID: %d\n", i+1, t.ID)
				fmt.Printf("   Saving Account ID: %d\n", t.SavingAccountID)
				fmt.Printf("   Member: %s (ID: %d)\n", t.MemberName, t.MemberID)
				fmt.Printf("   Type: %s\n", t.TransactionType)
				fmt.Printf("   Amount: Rp %d\n", t.Amount)
				fmt.Printf("   Balance Before: %s\n", t.BalanceBefore)
				fmt.Printf("   Balance After: %s\n", t.BalanceAfter)
				fmt.Printf("   Description: %s\n", t.Description)
				fmt.Printf("   Date: %s\n", t.TransactionDate)
				fmt.Printf("   Created: %s\n", t.CreatedAt)
				fmt.Println("-------------------------------------")
			}
		}
	}

	fmt.Println("\n✅ Verification Complete!")
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
