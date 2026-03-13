package repositories

import (
	"nusvakspps/config"
	"nusvakspps/models"
	"time"

	"gorm.io/gorm"
)

type SavingsRepository struct {
	db *gorm.DB
}

func NewSavingsRepository() *SavingsRepository {
	return &SavingsRepository{
		db: config.GetDB(),
	}
}

// GetDB returns the database connection
func (r *SavingsRepository) GetDB() *gorm.DB {
	return r.db
}

// CreateTransaction creates a new savings transaction
func (r *SavingsRepository) CreateTransaction(transaction *models.SavingTransaction) error {
	return r.db.Create(transaction).Error
}

// FindTransactionByID finds a transaction by ID with preloaded relations
func (r *SavingsRepository) FindTransactionByID(id uint) (*models.SavingTransaction, error) {
	var transaction models.SavingTransaction
	err := r.db.Preload("SavingAccount.AccountType").Preload("SavingAccount.Member").First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// ListTransactions retrieves transactions with pagination
func (r *SavingsRepository) ListTransactions(offset, limit int) ([]models.SavingTransaction, int64, error) {
	var transactions []models.SavingTransaction
	var total int64

	err := r.db.Model(&models.SavingTransaction{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Model(&models.SavingTransaction{}).
		Select("saving_transactions.id, saving_transactions.saving_account_id, saving_transactions.transaction_type, saving_transactions.amount, saving_transactions.description, saving_transactions.balance_before, saving_transactions.balance_after, saving_transactions.transaction_date, saving_transactions.created_by, saving_transactions.created_at").
		Preload("SavingAccount.AccountType").
		Preload("SavingAccount.Member").
		Order("saving_transactions.transaction_date DESC, saving_transactions.created_at DESC").
		Offset(offset).Limit(limit).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// ListTransactionsByMemberID retrieves transactions for a specific member
func (r *SavingsRepository) ListTransactionsByMemberID(memberID uint64) ([]models.SavingTransaction, error) {
	var transactions []models.SavingTransaction

	err := r.db.Model(&models.SavingTransaction{}).
		Select("saving_transactions.id, saving_transactions.saving_account_id, saving_transactions.transaction_type, saving_transactions.amount, saving_transactions.description, saving_transactions.balance_before, saving_transactions.balance_after, saving_transactions.transaction_date, saving_transactions.created_by, saving_transactions.created_at").
		Joins("JOIN saving_accounts ON saving_transactions.saving_account_id = saving_accounts.id").
		Where("saving_accounts.member_id = ?", memberID).
		Preload("SavingAccount.AccountType").
		Preload("SavingAccount.Member").
		Order("saving_transactions.transaction_date DESC, saving_transactions.created_at DESC").
		Find(&transactions).Error

	return transactions, err
}

// SearchTransactions searches transactions by keyword
func (r *SavingsRepository) SearchTransactions(keyword string, offset, limit int) ([]models.SavingTransaction, int64, error) {
	var transactions []models.SavingTransaction
	var total int64

	query := r.db.Model(&models.SavingTransaction{}).
		Select("saving_transactions.id, saving_transactions.saving_account_id, saving_transactions.transaction_type, saving_transactions.amount, saving_transactions.description, saving_transactions.balance_before, saving_transactions.balance_after, saving_transactions.transaction_date, saving_transactions.created_by, saving_transactions.created_at").
		Preload("SavingAccount.AccountType").
		Preload("SavingAccount.Member")

	if keyword != "" {
		query = query.Joins("LEFT JOIN saving_accounts ON saving_transactions.saving_account_id = saving_accounts.id").
			Joins("LEFT JOIN members ON saving_accounts.member_id = members.id").
			Joins("LEFT JOIN saving_types ON saving_accounts.account_type_id = saving_types.id").
			Where("members.full_name ILIKE ? OR saving_types.code ILIKE ? OR saving_types.name ILIKE ? OR saving_transactions.description ILIKE ?",
				"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("saving_transactions.transaction_date DESC, saving_transactions.created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// UpdateTransaction updates an existing transaction
func (r *SavingsRepository) UpdateTransaction(transaction *models.SavingTransaction) error {
	return r.db.Save(transaction).Error
}

// DeleteTransaction deletes a transaction
func (r *SavingsRepository) DeleteTransaction(id uint) error {
	return r.db.Delete(&models.SavingTransaction{}, id).Error
}

// GetMemberBalanceByType gets a member's balance for a specific saving type
func (r *SavingsRepository) GetMemberBalanceByType(memberID uint, accountTypeID uint) (float64, error) {
	var balance float64

	query := r.db.Model(&models.SavingTransaction{}).
		Joins("JOIN saving_accounts ON saving_transactions.saving_account_id = saving_accounts.id").
		Where("saving_accounts.member_id = ? AND saving_accounts.account_type_id = ?", memberID, accountTypeID)

	err := query.Select("COALESCE(SUM(CASE WHEN transaction_type = 'credit' THEN amount ELSE -amount END), 0)").Scan(&balance).Error
	if err != nil {
		return 0, err
	}

	return balance, nil
}

// GetMemberBalanceByCode gets a member's balance for a specific saving type code (backward compatibility)
func (r *SavingsRepository) GetMemberBalanceByCode(memberID uint, accountTypeCode string) (float64, error) {
	// First find the saving type ID
	savingTypeRepo := NewSavingTypesRepository()
	savingType, err := savingTypeRepo.FindByCode(accountTypeCode)
	if err != nil {
		return 0, err
	}

	return r.GetMemberBalanceByType(memberID, savingType.ID)
}

// FindSavingAccount finds a saving account by member ID and account type ID
func (r *SavingsRepository) FindSavingAccount(memberID uint, accountTypeID uint) (*models.SavingAccount, error) {
	var savingAccount models.SavingAccount
	err := r.db.Where("member_id = ? AND account_type_id = ?", memberID, accountTypeID).
		Preload("AccountType").
		First(&savingAccount).Error
	if err != nil {
		return nil, err
	}
	return &savingAccount, nil
}

// CreateSavingAccount creates a new saving account
func (r *SavingsRepository) CreateSavingAccount(account *models.SavingAccount) error {
	return r.db.Create(account).Error
}

// UpdateSavingAccount updates an existing saving account
func (r *SavingsRepository) UpdateSavingAccount(account *models.SavingAccount) error {
	return r.db.Save(account).Error
}

// GetMemberAllBalances gets all balances for a member
func (r *SavingsRepository) GetMemberAllBalances(memberID uint) (map[uint]float64, error) {
	balances := make(map[uint]float64)

	var results []struct {
		AccountTypeID uint
		Balance       float64
	}

	query := r.db.Model(&models.SavingTransaction{}).
		Select("saving_accounts.account_type_id, COALESCE(SUM(CASE WHEN transaction_type = 'credit' THEN amount ELSE -amount END), 0) as balance").
		Joins("JOIN saving_accounts ON saving_transactions.saving_account_id = saving_accounts.id").
		Where("saving_accounts.member_id = ?", memberID).
		Group("saving_accounts.account_type_id")

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		balances[result.AccountTypeID] = result.Balance
	}

	return balances, nil
}

// ListTransactionsWithFilters retrieves transactions with multiple filters
func (r *SavingsRepository) ListTransactionsWithFilters(offset, limit int, memberID uint, accountTypeID uint, fromDate, toDate *time.Time) ([]models.SavingTransaction, int64, error) {
	var transactions []models.SavingTransaction
	var total int64

	query := r.db.Model(&models.SavingTransaction{}).
		Select("saving_transactions.id, saving_transactions.saving_account_id, saving_transactions.transaction_type, saving_transactions.amount, saving_transactions.description, saving_transactions.balance_before, saving_transactions.balance_after, saving_transactions.transaction_date, saving_transactions.created_by, saving_transactions.created_at").
		Preload("SavingAccount.AccountType").
		Preload("SavingAccount.Member").
		Order("saving_transactions.transaction_date DESC, saving_transactions.created_at DESC")

	// Apply filters
	if memberID > 0 {
		query = query.Joins("JOIN saving_accounts ON saving_transactions.saving_account_id = saving_accounts.id").
			Where("saving_accounts.member_id = ?", memberID)
	}

	if accountTypeID > 0 {
		if memberID > 0 {
			// Both member and account type filters
			query = query.Where("saving_accounts.member_id = ? AND saving_accounts.account_type_id = ?", memberID, accountTypeID)
		} else {
			// Only account type filter
			query = query.Joins("JOIN saving_accounts ON saving_transactions.saving_account_id = saving_accounts.id").
				Where("saving_accounts.account_type_id = ?", accountTypeID)
		}
	}

	if fromDate != nil {
		if memberID > 0 || accountTypeID > 0 {
			query = query.Where("saving_transactions.transaction_date >= ?", fromDate)
		} else {
			query = query.Where("saving_transactions.transaction_date >= ?", fromDate)
		}
	}

	if toDate != nil {
		if memberID > 0 || accountTypeID > 0 {
			query = query.Where("saving_transactions.transaction_date <= ?", toDate)
		} else {
			query = query.Where("saving_transactions.transaction_date <= ?", toDate)
		}
	}

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset(offset).Limit(limit).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate running balance for each transaction
	for i := range transactions {
		if i == 0 {
			// First transaction shows initial balance
			transactions[i].BalanceBefore = 0
		}
	}

	return transactions, total, nil
}

// GetAllMemberBalances gets all member balances for all saving types
func (r *SavingsRepository) GetAllMemberBalances() ([]map[string]interface{}, int64, error) {
	var results []map[string]interface{}
	var total int64

	// Get all active saving types
	savingTypeRepo := NewSavingTypesRepository()
	savingTypes, err := savingTypeRepo.List()
	if err != nil {
		return nil, 0, err
	}

	// Get all members
	var members []models.Member
	err = r.db.Where("is_active = ?", true).Find(&members).Error
	if err != nil {
		return nil, 0, err
	}

	total = int64(len(members))

	// For each member, get balances for all saving types
	for _, member := range members {
		memberData := map[string]interface{}{
			"member_id":   member.ID,
			"member_name": member.FullName,
			"member_no":  member.MemberNo,
		}

		// Get balances for each saving type
		balances, err := r.GetMemberAllBalances(member.ID)
		if err != nil {
			return nil, 0, err
		}

		// Add balances to member data using saving type codes
		for _, savingType := range savingTypes {
			balance, exists := balances[savingType.ID]
			if !exists {
				balance = 0
			}
			memberData[savingType.Code] = int64(balance)
		}

		results = append(results, memberData)
	}

	return results, total, nil
}
