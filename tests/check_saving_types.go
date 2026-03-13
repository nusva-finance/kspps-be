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

	fmt.Println("📋 Checking saving_types table structure...")

	// Check saving_types structure
	var columns []struct {
		ColumnName string `json:"column_name"`
		DataType   string `json:"data_type"`
	}
	db.Raw(`
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_schema = 'public'
		AND table_name = 'saving_types'
		ORDER BY ordinal_position
	`).Scan(&columns)

	fmt.Println("📊 Table structure:")
	for _, col := range columns {
		fmt.Printf("   - %s (%s)\n", col.ColumnName, col.DataType)
	}

	// Get all saving types
	fmt.Println("\n📊 Current saving types in database:")
	var savingTypes []struct {
		ID          int     `json:"id"`
		Code        string  `json:"code"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		IsRequired  bool    `json:"is_required"`
		MinBalance  float64 `json:"min_balance"`
		IsActive    bool    `json:"is_active"`
		DisplayOrder int     `json:"display_order"`
	}
	db.Raw(`SELECT id, code, name, description, is_required, min_balance, is_active, display_order FROM saving_types ORDER BY display_order`).Scan(&savingTypes)

	for _, st := range savingTypes {
		var status string
		if st.IsActive {
			status = "✅ Active"
		} else {
			status = "❌ Inactive"
		}
		var required string
		if st.IsRequired {
			required = "Required"
		} else {
			required = "Optional"
		}
		fmt.Printf("   ID: %d\n", st.ID)
		fmt.Printf("     Code: %s\n", st.Code)
		fmt.Printf("     Name: %s\n", st.Name)
		fmt.Printf("     Description: %s\n", st.Description)
		fmt.Printf("     Status: %s (%s)\n", status, required)
		fmt.Printf("     Min Balance: %.2f\n", st.MinBalance)
		fmt.Printf("     Display Order: %d\n", st.DisplayOrder)
		fmt.Println()
	}

	// Check if saving_accounts table exists and its structure
	fmt.Println("📋 Checking saving_accounts table...")
	var tableExists bool
	db.Raw(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'saving_accounts'
		)
	`).Scan(&tableExists)

	if tableExists {
		fmt.Println("✅ saving_accounts table exists")

		// Check columns
		var accountColumns []struct {
			ColumnName string `json:"column_name"`
			DataType   string `json:"data_type"`
		}
		db.Raw(`
			SELECT column_name, data_type
			FROM information_schema.columns
			WHERE table_schema = 'public'
			AND table_name = 'saving_accounts'
			ORDER BY ordinal_position
		`).Scan(&accountColumns)

		fmt.Println("📊 saving_accounts structure:")
		for _, col := range accountColumns {
			fmt.Printf("   - %s (%s)\n", col.ColumnName, col.DataType)
		}

		// Get accounts
		var accountCount int64
		db.Raw(`SELECT COUNT(*) FROM saving_accounts`).Scan(&accountCount)
		fmt.Printf("📊 Total saving_accounts: %d\n", accountCount)

		// Check accounts with account_type_id
		var migratedCount int64
		db.Raw(`SELECT COUNT(*) FROM saving_accounts WHERE account_type_id IS NOT NULL`).Scan(&migratedCount)
		fmt.Printf("📊 Accounts with account_type_id: %d\n", migratedCount)

		// Sample account
		if accountCount > 0 {
			var sampleAccounts []struct {
				ID             int    `json:"id"`
				MemberID       int    `json:"member_id"`
				AccountTypeID  *int   `json:"account_type_id"`
				AccountNumber  string `json:"account_number"`
				Balance        float64 `json:"balance"`
			}
			db.Raw(`SELECT id, member_id, account_type_id, account_number, balance FROM saving_accounts LIMIT 3`).Scan(&sampleAccounts)

			fmt.Println("📊 Sample accounts:")
			for _, acc := range sampleAccounts {
				fmt.Printf("   ID: %d, MemberID: %d, AccountTypeID: %v, Balance: %.2f\n",
					acc.ID, acc.MemberID, acc.AccountTypeID, acc.Balance)
			}
		}
	} else {
		fmt.Println("❌ saving_accounts table doesn't exist")
	}

	fmt.Println("\n✅ Database check completed!")
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
