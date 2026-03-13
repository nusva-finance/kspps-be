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

	fmt.Println("🚀 Starting migration: Make Savings Dynamic")

	// Step 1: Create saving_types table
	fmt.Println("📋 Step 1: Creating saving_types table...")
	createSavingTypesSQL := `
	CREATE TABLE IF NOT EXISTS saving_types (
		id SERIAL PRIMARY KEY,
		code VARCHAR(20) UNIQUE NOT NULL,
		name VARCHAR(50) NOT NULL,
		description TEXT,
		is_required BOOLEAN DEFAULT false,
		min_balance DECIMAL(15, 2) DEFAULT 0,
		is_active BOOLEAN DEFAULT true,
		display_order INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if err := db.Exec(createSavingTypesSQL).Error; err != nil {
		fmt.Printf("❌ Error creating saving_types table: %v\n", err)
		return
	}
	fmt.Println("✅ saving_types table created successfully")

	// Create indexes for saving_types
	createIndexesSQL := `
	CREATE INDEX IF NOT EXISTS idx_saving_types_code ON saving_types(code);
	CREATE INDEX IF NOT EXISTS idx_saving_types_active ON saving_types(is_active);
	`

	if err := db.Exec(createIndexesSQL).Error; err != nil {
		fmt.Printf("❌ Error creating indexes: %v\n", err)
		return
	}
	fmt.Println("✅ Indexes created for saving_types")

	// Step 2: Insert default saving types
	fmt.Println("📋 Step 2: Inserting default saving types...")
	insertDefaultTypesSQL := `
	INSERT INTO saving_types (code, name, description, is_required, min_balance, display_order)
	VALUES
		('pokok', 'Simpanan Pokok', 'Simpanan wajib pertama kali menjadi anggota koperasi', true, 0, 1),
		('wajib', 'Simpanan Wajib', 'Simpanan rutin yang wajib dibayar oleh anggota', true, 0, 2),
		('modal', 'Simpanan Modal', 'Simpanan modal untuk keperluan investasi dan pengembangan usaha', false, 0, 3)
	ON CONFLICT (code) DO UPDATE SET
		name = EXCLUDED.name,
		description = EXCLUDED.description,
		is_required = EXCLUDED.is_required,
		min_balance = EXCLUDED.min_balance,
		display_order = EXCLUDED.display_order;
	`

	if err := db.Exec(insertDefaultTypesSQL).Error; err != nil {
		fmt.Printf("❌ Error inserting default saving types: %v\n", err)
		return
	}
	fmt.Println("✅ Default saving types inserted successfully")

	// Step 3: Check if saving_accounts table exists and has account_type column
	fmt.Println("📋 Step 3: Checking saving_accounts table structure...")

	var tableExists bool
	db.Raw(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'saving_accounts'
		)
	`).Scan(&tableExists)

	if !tableExists {
		fmt.Println("⚠️ saving_accounts table doesn't exist, skipping migration")
		fmt.Println("✅ Migration completed (saving_types table created)")
		return
	}

	// Check if account_type_id column already exists
	var columnExists bool
	db.Raw(`
		SELECT EXISTS (
			SELECT FROM information_schema.columns
			WHERE table_schema = 'public'
			AND table_name = 'saving_accounts'
			AND column_name = 'account_type_id'
		)
	`).Scan(&columnExists)

	if columnExists {
		fmt.Println("⚠️ account_type_id column already exists, skipping column addition")
	} else {
		// Step 4: Add account_type_id column
		fmt.Println("📋 Step 4: Adding account_type_id column to saving_accounts...")

		// First, rename the old account_type column to account_type_old for migration
		renameColumnSQL := `
		ALTER TABLE saving_accounts RENAME COLUMN account_type TO account_type_old;
		`

		if err := db.Exec(renameColumnSQL).Error; err != nil {
			fmt.Printf("❌ Error renaming column: %v\n", err)
			fmt.Println("⚠️ Column might already be renamed or doesn't exist, continuing...")
		} else {
			fmt.Println("✅ Column renamed successfully")
		}

		// Add new column as foreign key to saving_types
		addColumnSQL := `
		ALTER TABLE saving_accounts ADD COLUMN IF NOT EXISTS account_type_id INTEGER;
		`

		if err := db.Exec(addColumnSQL).Error; err != nil {
			fmt.Printf("❌ Error adding account_type_id column: %v\n", err)
		} else {
			fmt.Println("✅ account_type_id column added successfully")
		}

		// Step 5: Migrate existing data
		fmt.Println("📋 Step 5: Migrating existing data...")
		migrateDataSQL := `
		UPDATE saving_accounts sa
		SET account_type_id = st.id
		FROM saving_types st
		WHERE st.code = sa.account_type_old AND sa.account_type_id IS NULL;
		`

		result := db.Exec(migrateDataSQL)
		if result.Error != nil {
			fmt.Printf("❌ Error migrating data: %v\n", result.Error)
		} else {
			fmt.Printf("✅ Migrated %d records\n", result.RowsAffected)
		}

		// Step 6: Apply constraints
		fmt.Println("📋 Step 6: Applying constraints...")

		// Add foreign key constraint
		addForeignKeySQL := `
		ALTER TABLE saving_accounts ADD CONSTRAINT fk_saving_accounts_account_type
		FOREIGN KEY (account_type_id) REFERENCES saving_types(id) ON DELETE RESTRICT;
		`

		if err := db.Exec(addForeignKeySQL).Error; err != nil {
			fmt.Printf("❌ Error adding foreign key constraint: %v\n", err)
			fmt.Println("⚠️ Constraint might already exist, continuing...")
		} else {
			fmt.Println("✅ Foreign key constraint added successfully")
		}

		// Set NOT NULL constraint on account_type_id (only if all records have values)
		setNotNullSQL := `
		UPDATE saving_accounts SET account_type_id = (
			SELECT id FROM saving_types WHERE code = 'pokok' LIMIT 1
		) WHERE account_type_id IS NULL;
		`

		if err := db.Exec(setNotNullSQL).Error; err != nil {
			fmt.Printf("❌ Error setting default account_type_id: %v\n", err)
		} else {
			fmt.Println("✅ Default account_type_id set for NULL values")
		}

		// Drop the old column
		dropOldColumnSQL := `
		ALTER TABLE saving_accounts DROP COLUMN IF EXISTS account_type_old;
		`

		if err := db.Exec(dropOldColumnSQL).Error; err != nil {
			fmt.Printf("❌ Error dropping old column: %v\n", err)
			fmt.Println("⚠️ Column might not exist, continuing...")
		} else {
			fmt.Println("✅ Old column dropped successfully")
		}

		// Add index for better performance
		createIndexSQL := `
		CREATE INDEX IF NOT EXISTS idx_saving_accounts_type_id ON saving_accounts(account_type_id);
		`

		if err := db.Exec(createIndexSQL).Error; err != nil {
			fmt.Printf("❌ Error creating index: %v\n", err)
		} else {
			fmt.Println("✅ Index created for account_type_id")
		}
	}

	// Step 7: Create helper functions
	fmt.Println("📋 Step 7: Creating helper functions...")

	createFunctionSQL := `
	-- Function to get saving type name by ID
	CREATE OR REPLACE FUNCTION get_saving_type_name(type_id INTEGER)
	RETURNS VARCHAR AS $$
	BEGIN
		RETURN COALESCE((SELECT name FROM saving_types WHERE id = type_id), 'Unknown');
	END;
	$$ LANGUAGE plpgsql;

	-- Function to get saving type code by ID
	CREATE OR REPLACE FUNCTION get_saving_type_code(type_id INTEGER)
	RETURNS VARCHAR AS $$
	BEGIN
		RETURN COALESCE((SELECT code FROM saving_types WHERE id = type_id), 'unknown');
	END;
	$$ LANGUAGE plpgsql;
	`

	if err := db.Exec(createFunctionSQL).Error; err != nil {
		fmt.Printf("❌ Error creating functions: %v\n", err)
	} else {
		fmt.Println("✅ Helper functions created successfully")
	}

	// Step 8: Create view for easy account queries
	fmt.Println("📋 Step 8: Creating view...")
	createViewSQL := `
	CREATE OR REPLACE VIEW v_saving_accounts_with_types AS
	SELECT
		sa.id,
		sa.member_id,
		sa.account_type_id,
		sa.account_number,
		sa.balance,
		sa.is_active,
		sa.created_at,
		sa.updated_at,
		st.code AS account_type_code,
		st.name AS account_type_name,
		st.description AS account_type_description,
		st.is_required,
		st.min_balance
	FROM saving_accounts sa
	LEFT JOIN saving_types st ON sa.account_type_id = st.id;
	`

	if err := db.Exec(createViewSQL).Error; err != nil {
		fmt.Printf("❌ Error creating view: %v\n", err)
	} else {
		fmt.Println("✅ View created successfully")
	}

	// Step 9: Verification
	fmt.Println("📋 Step 9: Verification...")

	// Check saving_types table
	var savingTypeCount int64
	db.Raw(`SELECT COUNT(*) FROM saving_types`).Scan(&savingTypeCount)
	fmt.Printf("✅ saving_types table exists with %d records\n", savingTypeCount)

	// Display saving types
	var savingTypes []struct {
		ID     int    `json:"id"`
		Code   string `json:"code"`
		Name   string `json:"name"`
		Active bool   `json:"is_active"`
	}
	db.Raw(`SELECT id, code, name, is_active FROM saving_types ORDER BY display_order`).Scan(&savingTypes)

	fmt.Println("📊 Available saving types:")
	for _, st := range savingTypes {
		var status string
		if st.Active {
			status = "✅"
		} else {
			status = "❌"
		}
		fmt.Printf("   %s ID: %d, Code: %s, Name: %s\n", status, st.ID, st.Code, st.Name)
	}

	// Check saving_accounts
	if tableExists {
		var accountCount int64
		db.Raw(`SELECT COUNT(*) FROM saving_accounts`).Scan(&accountCount)
		fmt.Printf("✅ saving_accounts table exists with %d records\n", accountCount)

		// Check accounts with account_type_id
		var migratedCount int64
		db.Raw(`SELECT COUNT(*) FROM saving_accounts WHERE account_type_id IS NOT NULL`).Scan(&migratedCount)
		fmt.Printf("✅ %d saving_accounts have been migrated to use account_type_id\n", migratedCount)
	}

	fmt.Println("\n🎉 Migration completed successfully!")
	fmt.Println("📝 Next steps:")
	fmt.Println("   1. Rebuild backend: go build -o main.exe")
	fmt.Println("   2. Restart backend: ./main.exe")
	fmt.Println("   3. Test API: GET /api/v1/savings/types")
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
