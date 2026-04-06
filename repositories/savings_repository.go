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

// FindTransactionByIDWithRekening finds a transaction by ID with member and latest rekening info
func (r *SavingsRepository) FindTransactionByIDWithRekening(id uint) (*TransactionListRow, error) {
	var result TransactionListRow

	query := `
		SELECT
			st.id,
			st.saving_account_id,
			st.transaction_type,
			st.amount,
			st.description,
			st.balance_before,
			st.balance_after,
			st.transaction_date,
			st.created_by,
			st.created_at,
			sa.member_id,
			m.member_no,
			m.full_name as member_name,
			sa.account_type_id,
			sty.code as account_type,
			sty.name as account_type_name,
			COALESCE(rt.idnusvarekening, 0) as rekening_id,
			COALESCE(r.namarekening, '') as rekening_name,
			COALESCE(r.norekening, '') as rekening_no
		FROM saving_transactions st
		JOIN saving_accounts sa ON st.saving_account_id = sa.id
		JOIN members m ON sa.member_id = m.id
		JOIN saving_types sty ON sa.account_type_id = sty.id
		LEFT JOIN LATERAL (
			SELECT idnusvarekening
			FROM rekening_transaction rt2
			WHERE rt2.tabletransaction = 'saving_transactions'
			  AND rt2.idtabletransaction = st.id
			ORDER BY rt2.created_at DESC
			LIMIT 1
		) rt ON true
		LEFT JOIN nusva_rekening r ON rt.idnusvarekening = r.idnusvarekening
		WHERE st.id = ?
	`

	err := r.db.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
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
	// Use Updates with specific fields to avoid GORM trying to update associations
	return r.db.Model(&models.SavingTransaction{}).
		Where("id = ?", transaction.ID).
		Updates(map[string]interface{}{
			"saving_account_id": transaction.SavingAccountID,
			"transaction_type":  transaction.TransactionType,
			"amount":            transaction.Amount,
			"description":       transaction.Description,
			"transaction_date":  transaction.TransactionDate,
			"updated_by":        transaction.UpdatedBy,
		}).Error
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

// TransactionListRow represents a row in the transaction list response
type TransactionListRow struct {
	ID               uint      `json:"id"`
	SavingAccountID  uint      `json:"saving_account_id"`
	TransactionType  string    `json:"transaction_type"`
	Amount           float64   `json:"amount"`
	Description      string    `json:"description"`
	BalanceBefore    float64   `json:"balance_before"`
	BalanceAfter     float64   `json:"balance_after"`
	TransactionDate  time.Time `json:"transaction_date"`
	CreatedBy        string    `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	MemberID         uint      `json:"member_id"`
	MemberNo         string    `json:"member_no"`
	MemberName       string    `json:"member_name"`
	AccountTypeID    uint      `json:"account_type_id"`
	AccountType      string    `json:"account_type"`
	AccountTypeName  string    `json:"account_type_name"`
	RekeningID       uint      `json:"rekening_id"`
	RekeningName     string    `json:"rekening_name"`
	RekeningNo       string    `json:"rekening_no"`
}

// ListAllTransactions retrieves all transactions for the Simpanan page (with member info and rekening info)
func (r *SavingsRepository) ListAllTransactions(limit int) ([]TransactionListRow, int64, error) {
	var transactions []TransactionListRow
	var total int64

	err := r.db.Model(&models.SavingTransaction{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			st.id,
			st.saving_account_id,
			st.transaction_type,
			st.amount,
			st.description,
			st.balance_before,
			st.balance_after,
			st.transaction_date,
			st.created_by,
			st.created_at,
			sa.member_id,
			m.member_no,
			m.full_name as member_name,
			sa.account_type_id,
			sty.code as account_type,
			sty.name as account_type_name,
			COALESCE(rt.idnusvarekening, 0) as rekening_id,
			COALESCE(r.namarekening, '') as rekening_name,
			COALESCE(r.norekening, '') as rekening_no
		FROM saving_transactions st
		JOIN saving_accounts sa ON st.saving_account_id = sa.id
		JOIN members m ON sa.member_id = m.id
		JOIN saving_types sty ON sa.account_type_id = sty.id
		LEFT JOIN rekening_transaction rt ON rt.tabletransaction = 'saving_transactions' AND rt.idtabletransaction = st.id AND rt.transactiontype = 'Insert'
		LEFT JOIN nusva_rekening r ON rt.idnusvarekening = r.idnusvarekening
		ORDER BY st.transaction_date DESC, st.created_at DESC
	`

	if limit > 0 {
		query += " LIMIT ?"
		err = r.db.Raw(query, limit).Scan(&transactions).Error
	} else {
		err = r.db.Raw(query).Scan(&transactions).Error
	}

	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}
