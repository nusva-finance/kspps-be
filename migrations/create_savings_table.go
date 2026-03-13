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

	fmt.Println("📋 Creating saving_transactions table...")

	// Create saving_transactions table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS saving_transactions (
		id SERIAL PRIMARY KEY,
		saving_account_id INTEGER NOT NULL REFERENCES saving_accounts(id) ON DELETE CASCADE,
		transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('credit', 'debit')),
		amount DECIMAL(15, 2) NOT NULL,
		description TEXT,
		balance_before DECIMAL(15, 2),
		balance_after DECIMAL(15, 2),
		transaction_date DATE NOT NULL,
		created_by VARCHAR(100),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_saving_txns_account ON saving_transactions(saving_account_id);
	CREATE INDEX IF NOT EXISTS idx_saving_txns_date ON saving_transactions(transaction_date);
	`

	if err := db.Exec(createTableSQL).Error; err != nil {
		fmt.Printf("❌ Error creating saving_transactions table: %v\n", err)
		return
	}

	fmt.Println("✅ saving_transactions table created successfully")

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

	if tableName != "" {
		fmt.Println("✅ Confirmed: saving_transactions table exists")
	} else {
		fmt.Println("❌ Error: saving_transactions table was not created")
	}

	// Check indexes
	var indexes []struct {
		IndexName string `json:"index_name"`
	}
	err = db.Raw(`
		SELECT indexname as index_name
		FROM pg_indexes
		WHERE tablename = 'saving_transactions'
	`).Scan(&indexes).Error

	if err != nil {
		fmt.Printf("❌ Error checking indexes: %v\n", err)
		return
	}

	fmt.Printf("✅ Found %d indexes on saving_transactions table\n", len(indexes))
	for _, idx := range indexes {
		fmt.Printf("   - %s\n", idx.IndexName)
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
