# Migration Guide: Dynamic Savings Types System

## Overview
This migration transforms the savings system from hardcoded account types (`'pokok'`, `'wajib'`, `'manasuka'`) to a dynamic system using a new `saving_types` table.

## 📋 Migration Steps

### 1. **Run Database Migration**

Run the SQL migration script:

```bash
psql -U your_username -d your_database -f migrations/005_make_savings_dynamic.sql
```

This will:
- Create `saving_types` table
- Insert 3 default saving types (pokok, wajib, manasuka)
- Migrate existing `saving_accounts` to use foreign keys
- Update all references

### 2. **Rebuild Backend**

```bash
cd backend
go build -o main.exe
```

### 3. **Restart Backend Server**

```bash
./main.exe
```

### 4. **Rebuild Frontend (if needed)**

```bash
cd frontend
npm run build
```

## 🎯 What Changes

### Database Changes

**Before:**
```sql
CREATE TABLE saving_accounts (
    ...
    account_type VARCHAR(20) CHECK (account_type IN ('pokok', 'wajib', 'manasuka')),
    ...
);
```

**After:**
```sql
CREATE TABLE saving_types (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE,
    name VARCHAR(50),
    description TEXT,
    is_required BOOLEAN,
    min_balance DECIMAL,
    is_active BOOLEAN,
    display_order INTEGER
);

CREATE TABLE saving_accounts (
    ...
    account_type_id INTEGER REFERENCES saving_types(id),
    ...
);
```

### Backend Changes

**New Models:**
- `SavingType` - Model for dynamic saving types
- Updated `SavingAccount` - Uses foreign key to `SavingType`
- Updated `SavingTransaction` - Helper methods for backward compatibility

**New Repositories:**
- `SavingTypesRepository` - Manage saving types
- Updated `SavingsRepository` - Use dynamic types

**New Handlers:**
- `GetSavingTypes`, `GetSavingTypeByID`
- `CreateSavingType`, `UpdateSavingType`, `DeleteSavingType`
- `InitializeSavingTypes` - Initialize default types

**New Routes:**
- `GET /api/v1/savings/types` - Get all saving types
- `POST /api/v1/savings/types` - Create new saving type
- `PUT /api/v1/savings/types/:id` - Update saving type
- `DELETE /api/v1/savings/types/:id` - Delete saving type
- `POST /api/v1/savings/types/initialize` - Initialize defaults

### Frontend Changes

**New Service Method:**
- `savingsService.getSavingTypes()` - Fetch available types

**Updated Components:**
- `Savings.tsx` - Dynamic columns based on saving types
- Form components - Dynamic dropdown options

## 🚀 Usage Examples

### Get All Saving Types (Frontend)
```typescript
const savingTypes = await savingsService.getSavingTypes()
// Returns: [{ id: 1, code: 'pokok', name: 'Simpanan Pokok', ... }, ...]
```

### Create New Saving Type (API)
```bash
POST /api/v1/savings/types
{
  "code": "pendidikan",
  "name": "Simpanan Pendidikan",
  "description": "Simpanan untuk keperluan pendidikan",
  "is_required": false,
  "min_balance": 0,
  "is_active": true,
  "display_order": 4
}
```

### Initialize Default Saving Types (API)
```bash
POST /api/v1/savings/types/initialize
```

## 📊 Default Saving Types

After migration, these saving types will be automatically created:

| ID | Code | Name | Description | Required |
|----|------|------|-------------|-----------|
| 1  | pokok | Simpanan Pokok | Simpanan wajib pertama kali menjadi anggota koperasi | Yes |
| 2  | wajib | Simpanan Wajib | Simpanan rutin yang wajib dibayar oleh anggota | Yes |
| 3  | manasuka | Simpanan Manasuka | Simpanan sukarela yang dapat ditarik kapan saja | No |

## 🔧 Manual Data Migration (if needed)

If you have existing data and the automatic migration fails:

```sql
-- Check existing saving accounts
SELECT * FROM saving_accounts;

-- Manually update account types
UPDATE saving_accounts
SET account_type_id = (
  SELECT id FROM saving_types WHERE code = account_type_old
)
WHERE account_type_id IS NULL;
```

## 🐛 Troubleshooting

### Issue: "Invalid account type" error
**Solution:** Initialize default saving types:
```bash
curl -X POST http://localhost:8080/api/v1/savings/types/initialize
```

### Issue: Columns still hardcoded
**Solution:** Ensure frontend is rebuilt and backend is restarted after migration

### Issue: Foreign key constraint errors
**Solution:** Check that all saving accounts have valid `account_type_id`:
```sql
SELECT * FROM saving_accounts WHERE account_type_id IS NULL;
```

## 📝 Adding New Saving Types

1. **Via API:**
```bash
curl -X POST http://localhost:8080/api/v1/savings/types \
  -H "Content-Type: application/json" \
  -d '{
    "code": "investasi",
    "name": "Simpanan Investasi",
    "description": "Simpanan dengan return investasi",
    "is_required": false,
    "min_balance": 1000000,
    "is_active": true,
    "display_order": 5
  }'
```

2. **Via Database (Direct):**
```sql
INSERT INTO saving_types (code, name, description, is_required, min_balance, is_active, display_order)
VALUES ('investasi', 'Simpanan Investasi', 'Simpanan dengan return investasi', false, 1000000, true, 5);
```

## ✅ Verification Steps

After migration, verify:

1. **Saving types table exists:**
```sql
SELECT * FROM saving_types;
```

2. **Saving accounts have foreign keys:**
```sql
SELECT * FROM saving_accounts WHERE account_type_id IS NULL;
```

3. **API returns saving types:**
```bash
curl http://localhost:8080/api/v1/savings/types
```

4. **Frontend displays dynamic columns:**
- Open Savings page
- Check if columns match saving types
- Verify data is displayed correctly

## 🔄 Rollback Plan

If you need to rollback:

```sql
-- Drop new table
DROP TABLE IF EXISTS saving_types CASCADE;

-- Revert saving_accounts
ALTER TABLE saving_accounts ADD COLUMN account_type VARCHAR(20);
UPDATE saving_accounts SET account_type = 'pokok' WHERE account_type IS NULL;
ALTER TABLE saving_accounts ALTER COLUMN account_type SET NOT NULL;
ALTER TABLE saving_accounts ADD CONSTRAINT check_account_type
  CHECK (account_type IN ('pokok', 'wajib', 'manasuka'));

-- Drop account_type_id
ALTER TABLE saving_accounts DROP COLUMN account_type_id;
```

## 📞 Support

If you encounter issues:
1. Check backend logs: `./main.exe`
2. Check browser console for frontend errors
3. Verify database schema: `\d saving_accounts`
4. Check API endpoints: `/api/v1/savings/types`
