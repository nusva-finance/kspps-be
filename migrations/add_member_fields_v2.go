package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Migration to update members table structure (Simple & Safe version)
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

	fmt.Println("🔄 Starting members table migration (v2)...")

	// Create schema_migrations table if not exists
	createSchemaMigrationsTable(db)

	// Check if this migration has already run
	version := int64(2024030902) // New version
	if hasMigrationRun(db, version) {
		fmt.Println("✅ Migration already run, skipping.")
		return
	}

	// Step 1: Add new columns one by one (no transaction for safety)
	newColumns := []struct {
		columnName string
		dataType   string
	}{
		{"join_year", "CHAR(2)"},
		{"join_month", "CHAR(2)"},
		{"emergency_name", "VARCHAR(100)"},
		{"emergency_relation", "VARCHAR(50)"},
		{"emergency_phone", "VARCHAR(20)"},
		{"emergency_address", "TEXT"},
		{"company_name", "VARCHAR(100)"},
		{"job_title", "VARCHAR(100)"},
		{"bank_account_no", "VARCHAR(30)"},
		{"bank_name", "VARCHAR(50)"},
	}

	for _, col := range newColumns {
		addColumnIfNotExists(db, col.columnName, col.dataType)
	}

	// Step 2: Check and handle column renaming
	// Check if 'address' exists and 'address_ktp' doesn't exist
	var addressExists, addressKtpExists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'members' AND column_name = 'address')").Scan(&addressExists)
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'members' AND column_name = 'address_ktp')").Scan(&addressKtpExists)

	if addressExists && !addressKtpExists {
		fmt.Println("🔄 Renaming address to address_ktp...")
		err := db.Exec("ALTER TABLE members RENAME COLUMN address TO address_ktp").Error
		if err != nil {
			fmt.Printf("❌ Error renaming address to address_ktp: %v\n", err)
		} else {
			fmt.Println("✅ Renamed address to address_ktp")
		}
	}

	// Check if 'nik' exists and 'ktp_no' doesn't exist
	var nikExists, ktpNoExists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'members' AND column_name = 'nik')").Scan(&nikExists)
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'members' AND column_name = 'ktp_no')").Scan(&ktpNoExists)

	if nikExists && !ktpNoExists {
		fmt.Println("🔄 Renaming nik to ktp_no...")
		err := db.Exec("ALTER TABLE members RENAME COLUMN nik TO ktp_no").Error
		if err != nil {
			fmt.Printf("❌ Error renaming nik to ktp_no: %v\n", err)
		} else {
			fmt.Println("✅ Renamed nik to ktp_no")

			// Add UNIQUE constraint on ktp_no
			fmt.Println("🔄 Adding UNIQUE constraint on ktp_no...")
			err = db.Exec("ALTER TABLE members ADD CONSTRAINT members_ktp_no_key UNIQUE (ktp_no)").Error
			if err != nil {
				fmt.Printf("❌ Error adding UNIQUE constraint: %v\n", err)
			} else {
				fmt.Println("✅ Added UNIQUE constraint on ktp_no")
			}
		}
	}

	// Check if 'npwp' exists and 'npwp_no' doesn't exist
	var npwpExists, npwpNoExists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'members' AND column_name = 'npwp')").Scan(&npwpExists)
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'members' AND column_name = 'npwp_no')").Scan(&npwpNoExists)

	if npwpExists && !npwpNoExists {
		fmt.Println("🔄 Renaming npwp to npwp_no...")
		err := db.Exec("ALTER TABLE members RENAME COLUMN npwp TO npwp_no").Error
		if err != nil {
			fmt.Printf("❌ Error renaming npwp to npwp_no: %v\n", err)
		} else {
			fmt.Println("✅ Renamed npwp to npwp_no")
		}
	}

	// Step 3: Update existing data to populate join_year and join_month
	fmt.Println("🔄 Updating existing data to populate join_year and join_month...")
	result := db.Exec(`
		UPDATE members
		SET join_year = EXTRACT(YEAR FROM join_date)::CHAR(2),
		    join_month = EXTRACT(MONTH FROM join_date)::CHAR(2)
		WHERE (join_year IS NULL OR join_month IS NULL) AND join_date IS NOT NULL
	`)
	if result.Error != nil {
		fmt.Printf("⚠️  Warning updating join_year/join_month: %v\n", result.Error)
	} else {
		fmt.Println("✅ Updated join_year and join_month for existing records")
	}

	// Record migration
	err = recordMigration(db, version)
	if err != nil {
		fmt.Printf("❌ Error recording migration: %v\n", err)
		return
	}

	fmt.Println("✅ Migration completed successfully!")
	printTableStructure(db)
}

func addColumnIfNotExists(db *gorm.DB, columnName, dataType string) {
	var exists bool
	err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns
			WHERE table_name = 'members' AND column_name = ?
		)
	`, columnName).Scan(&exists).Error

	if err != nil {
		fmt.Printf("❌ Error checking column %s: %v\n", columnName, err)
		return
	}

	if !exists {
		fmt.Printf("➕ Adding column: %s (%s)\n", columnName, dataType)
		err := db.Exec(fmt.Sprintf("ALTER TABLE members ADD COLUMN %s %s", columnName, dataType)).Error
		if err != nil {
			fmt.Printf("❌ Error adding column %s: %v\n", columnName, err)
		}
	} else {
		fmt.Printf("✓ Column %s already exists\n", columnName)
	}
}

func createSchemaMigrationsTable(db *gorm.DB) {
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	if err != nil {
		fmt.Printf("Warning: Failed to create schema_migrations table: %v\n", err)
	}
}

func hasMigrationRun(db *gorm.DB, version int64) bool {
	var hasRun bool
	err := db.Table("schema_migrations").Where("version = ?", version).Select("version").First(&hasRun).Error
	if err != nil && err.Error() != "record not found" {
		fmt.Printf("Warning: Error checking migration status: %v\n", err)
	}
	return hasRun
}

func recordMigration(db *gorm.DB, version int64) error {
	return db.Exec(`
		INSERT INTO schema_migrations (version, applied_at)
		VALUES (?, CURRENT_TIMESTAMP)
		ON CONFLICT (version) DO NOTHING
	`, version).Error
}

func printTableStructure(db *gorm.DB) {
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

	fmt.Println("\n📋 Updated members table structure:")
	for rows.Next() {
		var columnName, dataType string
		var maxLength *int

		err := rows.Scan(&columnName, &dataType, &maxLength)
		if err != nil {
			continue
		}

		lengthInfo := ""
		if maxLength != nil {
			lengthInfo = fmt.Sprintf("(%d)", *maxLength)
		}

		fmt.Printf("  %-25s %-25s %s\n", columnName, dataType, lengthInfo)
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
