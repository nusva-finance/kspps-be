package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	// Debug: Print environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	log.Printf("Database Config - Host: %s, Port: %s, User: %s, DB: %s", host, port, user, dbName)

	// Check if critical config is missing
	if host == "" || port == "" || user == "" || dbName == "" {
		return fmt.Errorf("missing required database configuration. Please check .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbName,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// Run migrations - check what has been run and run new ones
	log.Println("Running database migrations...")
	err = runMigrations(DB)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database connected and migrations completed successfully")
	return nil
}

// addDeletedAtColumn adds deleted_at column to users table
func addDeletedAtColumn(db *gorm.DB) error {
	var hasColumn bool
	err := db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'deleted_at')").Scan(&hasColumn).Error
	if err != nil {
		return err
	}

	if hasColumn {
		return nil // Column already exists
	}

	// Add deleted_at column
	err = db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP").Error
	if err != nil {
		return err
	}

	// Create index on deleted_at for better query performance
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)").Error
	if err != nil {
		log.Printf("Warning: Failed to create index on deleted_at: %v", err)
	}

	return nil
}

// createKategoriBarangsTable creates kategori_barang table
func createKategoriBarangsTable(db *gorm.DB) error {
	// Check if table exists
	var tableExists bool
	err := db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'kategori_barang')").Scan(&tableExists).Error
	if err != nil {
		return err
	}

	if tableExists {
		return nil // Table already exists
	}

	// Create table (tanpa 's' di nama tabel)
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS kategori_barang (
			id SERIAL PRIMARY KEY,
			kategori VARCHAR(100) UNIQUE NOT NULL,
			aktif BOOLEAN DEFAULT true,
			created_by VARCHAR(100),
			updated_by VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP
		)
	`).Error

	if err != nil {
		return err
	}

	// Create indexes
	err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_kategori_barang_kategori ON kategori_barang(kategori)`).Error
	if err != nil {
		log.Printf("Warning: Failed to create index: %v", err)
	}

	return nil
}

// createMarginSetupsTable creates margin_setups table
func createMarginSetupsTable(db *gorm.DB) error {
	// Check if table exists
	var tableExists bool
	err := db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'margin_setups')").Scan(&tableExists).Error
	if err != nil {
		return err
	}

	if tableExists {
		return nil // Table already exists
	}

	// Create table
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS margin_setups (
			category VARCHAR(50) NOT NULL,
			tenor INTEGER NOT NULL,
			margin DECIMAL(5,4) NOT NULL,
			deleted_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error

	if err != nil {
		return err
	}

	// Create indexes
	err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_margin_setups_category ON margin_setups(category)`).Error
	if err != nil {
		log.Printf("Warning: Failed to create index: %v", err)
	}

	err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_margin_setups_tenor ON margin_setups(tenor)`).Error
	if err != nil {
		log.Printf("Warning: Failed to create index: %v", err)
	}

	err = db.Exec(`
		CREATE UNIQUE CONSTRAINT unique_category_tenor UNIQUE (category, tenor)
	`).Error

	if err != nil {
		log.Printf("Warning: Failed to create unique constraint: %v", err)
	}

	// Insert sample data
	err = db.Exec(`
		INSERT INTO margin_setups (category, tenor, margin) VALUES
		('Elektronik', 3, 0.05),
		('Elektronik', 6, 0.07),
		('Elektronik', 12, 0.10),
		('Pakaian', 3, 0.04),
		('Pakaian', 6, 0.06),
		('Kendaraan', 12, 0.12),
		('Kendaraan', 24, 0.15),
		('Kendaraan', 36, 0.18)
		ON CONFLICT (category, tenor) DO NOTHING
	`).Error

	if err != nil {
		log.Printf("Warning: Failed to insert sample data: %v", err)
	}

	return nil
}

// createRekeningTransactionsTable creates rekening_transactions table
func createRekeningTransactionsTable(db *gorm.DB) error {
	// Check if table exists
	var tableExists bool
	err := db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'rekening_transactions')").Scan(&tableExists).Error
	if err != nil {
		return err
	}

	if tableExists {
		return nil // Table already exists
	}

	// Create table
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS rekening_transactions (
			id SERIAL PRIMARY KEY,
			rekening_id INTEGER NOT NULL,
			transaction_type VARCHAR(50) NOT NULL,
			amount DECIMAL(18,2) NOT NULL,
			description TEXT,
			reference_id INTEGER,
			reference_type VARCHAR(50),
			created_by VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error

	if err != nil {
		return err
	}

	// Create indexes
	err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_rekening_transactions_rekening_id ON rekening_transactions(rekening_id)`).Error
	if err != nil {
		log.Printf("Warning: Failed to create index: %v", err)
	}

	err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_rekening_transactions_type ON rekening_transactions(transaction_type)`).Error
	if err != nil {
		log.Printf("Warning: Failed to create index: %v", err)
	}

	return nil
}

// runMigrations handles all database migrations in a controlled way
func runMigrations(db *gorm.DB) error {
	// Migration 1: Add deleted_at column to users table
	if err := runMigrationOnce(db, 2024030301, addDeletedAtColumn); err != nil {
		return fmt.Errorf("migration %d failed: %w", 2024030301, err)
	}

	// Migration 2: Create kategori_barangs table
	if err := runMigrationOnce(db, 2024031501, createKategoriBarangsTable); err != nil {
		return fmt.Errorf("migration %d failed: %w", 2024031501, err)
	}

	// Migration 3: Create margin_setups table
	if err := runMigrationOnce(db, 2024031502, createMarginSetupsTable); err != nil {
		return fmt.Errorf("migration %d failed: %w", 2024031502, err)
	}

	// Migration 4: Create rekening_transaction table
	if err := runMigrationOnce(db, 2025031801, createRekeningTransactionsTable); err != nil {
		return fmt.Errorf("migration %d failed: %w", 2025031801, err)
	}

	return nil
}

// runMigrationOnce runs a migration function only if version hasn't been run yet
func runMigrationOnce(db *gorm.DB, version int64, migrationFunc func(*gorm.DB) error) error {
	// Check if migration has already been run
	var hasRun bool
	err := db.Table("schema_migrations").Where("version = ?", version).Select("version").First(&hasRun).Error
	if err != nil {
		// Create schema_migrations table if it doesn't exist
		if !isTableNotFoundError(err) {
			log.Printf("Creating schema_migrations table: %v", err)
			createErr := db.Exec(`
				CREATE TABLE IF NOT EXISTS schema_migrations (
					version BIGINT PRIMARY KEY,
					applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				)
			`).Error
			if createErr != nil {
				return fmt.Errorf("failed to create schema_migrations table: %w", createErr)
			}
		} else {
			return fmt.Errorf("failed to check migration status: %w", err)
		}
	}

	// If migration already run, skip it
	if hasRun {
		log.Printf("Migration %d already run, skipping", version)
		return nil
	}

	// Run migration
	log.Printf("Running migration %d...", version)
	err = migrationFunc(db)
	if err != nil {
		return fmt.Errorf("migration %d failed: %w", version, err)
	}

	// Record that migration was run
	err = db.Exec(`
		INSERT INTO schema_migrations (version, applied_at)
		VALUES (?, CURRENT_TIMESTAMP)
		ON CONFLICT (version) DO NOTHING
	`, version).Error
	if err != nil {
		log.Printf("Warning: Failed to record migration %d: %v", version)
	}

	return nil
}

// isTableNotFoundError checks if error is because a table doesn't exist
func isTableNotFoundError(err error) bool {
	return err != nil && (err.Error() == "relation \"schema_migrations\" does not exist" ||
		err.Error() == `pq: relation "schema_migrations" does not exist`)
}

func GetDB() *gorm.DB {
	return DB
}
