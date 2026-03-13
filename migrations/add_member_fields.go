package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Migration to update members table structure
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

	fmt.Println("🔄 Starting members table migration...")

	// Create schema_migrations table if not exists
	createSchemaMigrationsTable(db)

	// Check if this migration has already run
	version := int64(2024030901) // Migration version for member fields update
	if hasMigrationRun(db, version) {
		fmt.Println("✅ Migration already run, skipping.")
		return
	}

	// Start transaction for atomic migration
	tx := db.Begin()
	if tx.Error != nil {
		fmt.Printf("❌ Error starting transaction: %v\n", tx.Error)
		return
	}

	// Step 1: Add new columns if they don't exist
	migrations := []struct {
		columnName string
		dataType   string
		comment    string
	}{
		// Split join_date into join_year and join_month
		{"join_year", "CHAR(2)", "YY - tahun masuk"},
		{"join_month", "CHAR(2)", "MM - bulan masuk"},

		// Add missing emergency contact fields
		{"emergency_name", "VARCHAR(100)", "Nama kontak darurat"},
		{"emergency_relation", "VARCHAR(50)", "Hubungan dengan anggota"},
		{"emergency_phone", "VARCHAR(20)", "No. telp kontak darurat"},
		{"emergency_address", "TEXT", "Alamat kontak darurat"},

		// Add work information
		{"company_name", "VARCHAR(100)", "Nama perusahaan tempat bekerja"},
		{"job_title", "VARCHAR(100)", "Nama jabatan"},

		// Add bank information
		{"bank_account_no", "VARCHAR(30)", "Nomor rekening bank"},
		{"bank_name", "VARCHAR(50)", "Nama bank"},

		// Rename existing columns
		{"address_ktp", "TEXT", "Alamat sesuai KTP (will be renamed from address)"},
	}

	for _, migration := range migrations {
		addColumnIfNotExists(tx, migration.columnName, migration.dataType)
	}

	// Step 2: Rename existing columns (check if old column exists first)
	// Rename 'address' to 'address_ktp'
	var addressExists bool
	tx.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'members' AND column_name = 'address')").Scan(&addressExists)
	if addressExists {
		err = tx.Exec("ALTER TABLE members RENAME COLUMN address TO address_ktp").Error
		if err != nil {
			log.Printf("⚠️  Warning renaming address to address_ktp: %v", err)
		} else {
			fmt.Println("✓ Renamed address to address_ktp")
		}
	}

	// Rename 'nik' to 'ktp_no'
	var nikExists bool
	tx.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'members' AND column_name = 'nik')").Scan(&nikExists)
	if nikExists {
		err = tx.Exec("ALTER TABLE members RENAME COLUMN nik TO ktp_no").Error
		if err != nil {
			log.Printf("⚠️  Warning renaming nik to ktp_no: %v", err)
		} else {
			fmt.Println("✓ Renamed nik to ktp_no")
		}
	}

	// Rename 'npwp' to 'npwp_no'
	var npwpExists bool
	tx.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'members' AND column_name = 'npwp')").Scan(&npwpExists)
	if npwpExists {
		err = tx.Exec("ALTER TABLE members RENAME COLUMN npwp TO npwp_no").Error
		if err != nil {
			log.Printf("⚠️  Warning renaming npwp to npwp_no: %v", err)
		} else {
			fmt.Println("✓ Renamed npwp to npwp_no")
		}
	}

	// Step 3: Add UNIQUE constraint on ktp_no (check if exists first)
	var constraintExists bool
	tx.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE table_name = 'members' AND constraint_name = 'members_ktp_no_key')").Scan(&constraintExists)
	if !constraintExists {
		err = tx.Exec("ALTER TABLE members ADD CONSTRAINT members_ktp_no_key UNIQUE (ktp_no)").Error
		if err != nil {
			log.Printf("⚠️  Warning adding UNIQUE constraint on ktp_no: %v", err)
		} else {
			fmt.Println("✓ Added UNIQUE constraint on ktp_no")
		}
	}

	// Step 4: Update existing data to populate join_year and join_month from join_date
	fmt.Println("🔄 Updating existing data to populate join_year and join_month...")
	err = tx.Exec(`
		UPDATE members
		SET join_year = EXTRACT(YEAR FROM join_date)::CHAR(2),
		    join_month = EXTRACT(MONTH FROM join_date)::CHAR(2)
		WHERE join_year IS NULL OR join_month IS NULL
	`).Error
	if err != nil {
		log.Printf("⚠️  Warning updating join_year/join_month: %v", err)
	}

	// Step 4: Update existing data to populate join_year and join_month from join_date
	fmt.Println("🔄 Updating existing data to populate join_year and join_month...")
	err = tx.Exec(`
		UPDATE members
		SET join_year = EXTRACT(YEAR FROM join_date)::CHAR(2),
		    join_month = EXTRACT(MONTH FROM join_date)::CHAR(2)
		WHERE join_year IS NULL OR join_month IS NULL
	`).Error
	if err != nil {
		log.Printf("⚠️  Warning updating join_year/join_month: %v", err)
	}

	// Record migration
	err = recordMigration(tx, version)
	if err != nil {
		log.Printf("⚠️  Warning recording migration: %v", err)
		tx.Rollback()
		fmt.Printf("❌ Error rolling back migration: %v\n", err)
		return
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		fmt.Printf("❌ Error committing migration: %v\n", err)
		tx.Rollback()
		return
	}

	fmt.Println("✅ Migration completed successfully!")
	fmt.Println("\n📋 Updated members table structure:")
	printTableStructure(db)
}

func addColumnIfNotExists(db *gorm.DB, columnName, dataType string) {
	// Check if column exists
	var exists bool
	err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns
			WHERE table_name = 'members' AND column_name = ?
		)
	`, columnName).Scan(&exists).Error

	if err != nil {
		log.Printf("❌ Error checking column %s: %v", columnName, err)
		return
	}

	if !exists {
		fmt.Printf("➕ Adding column: %s (%s)\n", columnName, dataType)
		err := db.Exec(fmt.Sprintf("ALTER TABLE members ADD COLUMN %s %s", columnName, dataType)).Error
		if err != nil {
			log.Printf("❌ Error adding column %s: %v", columnName, err)
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
		log.Printf("Warning: Failed to create schema_migrations table: %v", err)
	}
}

func hasMigrationRun(db *gorm.DB, version int64) bool {
	var hasRun bool
	err := db.Table("schema_migrations").Where("version = ?", version).Select("version").First(&hasRun).Error
	if err != nil && err.Error() != "record not found" {
		log.Printf("Warning: Error checking migration status: %v", err)
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

	for rows.Next() {
		var columnName, dataType string
		var maxLength *int

		err := rows.Scan(&columnName, &dataType, &maxLength)
		if err != nil {
			log.Printf("❌ Error scanning row: %v\n", err)
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
