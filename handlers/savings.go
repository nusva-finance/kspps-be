package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"nusvakspps/models"
	"nusvakspps/repositories"
)

type CreateSavingsTransactionRequest struct {
	AccountType     string `json:"account_type" binding:"required"`     // pokok, wajib, manasuka (code)
	MemberID        string `json:"member_id" binding:"required"`        // string from frontend
	MemberName      string `json:"member_name" binding:"required"`      // redundant but keeps transaction self-contained
	TransactionType string `json:"transaction_type" binding:"required"` // credit, debit
	Amount          int    `json:"amount" binding:"required"`           // integer amount (in rupiah)
	Description     string `json:"description" binding:"required"`
	TransactionDate string `json:"transaction_date" binding:"required"` // YYYY-MM-DD format
}

// GetSavingsTransactions retrieves all savings transactions with pagination, search, and filters
// GetSavingsTransactions retrieves all savings transactions with opening balance calculation
func GetSavingsTransactions(c *gin.Context) {
	page := 1
	limit := 10
	memberIDParam := c.Query("member_id")
	accountTypeCode := c.Query("account_type") 
	dateFrom := c.Query("date_from")      // Digunakan untuk filter mutasi & batas saldo awal
	dateTo := c.Query("date_to")          

	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	savingsRepo := repositories.NewSavingsRepository()

	// 1. Parse dates untuk filter list
	var fromDate, toDate *time.Time
	if dateFrom != "" {
		parsedDate, _ := time.Parse("2006-01-02", dateFrom)
		fromDate = &parsedDate
	}
	if dateTo != "" {
		parsedDate, _ := time.Parse("2006-01-02", dateTo)
		toDate = &parsedDate
	}

	// 2. Parse IDs
	memberIDParsed, _ := strconv.ParseUint(memberIDParam, 10, 32)
	memberID := uint(memberIDParsed)

	var accountTypeID uint
	if accountTypeCode != "" {
		savingTypeRepo := repositories.NewSavingTypesRepository()
		savingType, _ := savingTypeRepo.FindByCode(accountTypeCode)
		accountTypeID = savingType.ID
	}

	// --- LOGIKA HITUNG SALDO AWAL (OPENING BALANCE) ---
	var openingBalance int64 = 0
	if memberID > 0 && accountTypeID > 0 && dateFrom != "" {
		db := savingsRepo.GetDB()
		var result struct {
			Total float64
		}
		
		// Menghitung semua transaksi SEBELUM tanggal filter (dateFrom)
		err := db.Table("saving_transactions").
			Joins("JOIN saving_accounts ON saving_transactions.saving_account_id = saving_accounts.id").
			Where("saving_accounts.member_id = ?", memberID).
			Where("saving_accounts.account_type_id = ?", accountTypeID).
			Where("saving_transactions.transaction_date < ?", dateFrom). // Kunci: < dateFrom
			Select("COALESCE(SUM(CASE WHEN transaction_type = 'credit' THEN amount ELSE -amount END), 0) as total").
			Scan(&result).Error

		if err == nil {
			openingBalance = int64(result.Total)
		}
	}

	// 3. Ambil List Mutasi
	transactions, total, err := savingsRepo.ListTransactionsWithFilters((page-1)*limit, limit, memberID, accountTypeID, fromDate, toDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"transactions":    transactions,
			"opening_balance": openingBalance, // Nilai 600.000 akan dikirim ke sini
		},
		"total": total,
		"page":  page,
		"limit": limit,
	})
}


// GetSavingsTransactionByID retrieves a single savings transaction by ID
func GetSavingsTransactionByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		fmt.Println("❌ Invalid transaction ID:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	fmt.Println("👤 Getting savings transaction by ID:", id)

	savingsRepo := repositories.NewSavingsRepository()
	transaction, err := savingsRepo.FindTransactionByID(uint(id))
	if err != nil {
		fmt.Println("❌ Transaction not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	fmt.Println("✅ Found transaction:", transaction.ID)

	c.JSON(http.StatusOK, transaction)
}

// CreateSavingsTransaction creates a new savings transaction
func CreateSavingsTransaction(c *gin.Context) {
    // --- PERUBAHAN 1: Ambil nama user dari middleware ---
    operatorName := c.GetString("current_user_name")
    if operatorName == "" {
        operatorName = "system"
    }

    var req CreateSavingsTransactionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        fmt.Println("❌ Error binding request:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    fmt.Println("📤 Creating savings transaction:")
    fmt.Printf("  MemberID: %s\n", req.MemberID)
    fmt.Printf("  MemberName: %s\n", req.MemberName)
    fmt.Printf("  AccountType: %s\n", req.AccountType)
    fmt.Printf("  TransactionType: %s\n", req.TransactionType)
    fmt.Printf("  Amount: %d\n", req.Amount)
    fmt.Printf("  Description: %s\n", req.Description)
    fmt.Printf("  TransactionDate: %s\n", req.TransactionDate)

    // Validate saving type code
    savingTypeRepo := repositories.NewSavingTypesRepository()
    savingType, err := savingTypeRepo.FindByCode(req.AccountType)
    if err != nil {
        fmt.Println("❌ Invalid account type:", req.AccountType)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account type"})
        return
    }

    // Validate transaction type
    if req.TransactionType != "credit" && req.TransactionType != "debit" {
        fmt.Println("❌ Invalid transaction type:", req.TransactionType)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type. Must be 'credit' or 'debit'"})
        return
    }

    // Validate amount
    if req.Amount <= 0 {
        fmt.Println("❌ Invalid amount:", req.Amount)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
        return
    }

    // Parse transaction date
    transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
    if err != nil {
        fmt.Println("❌ Error parsing transaction date:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction date format. Use YYYY-MM-DD"})
        return
    }

    // Verify member exists
    memberRepo := repositories.NewMemberRepository()
    memberID, err := strconv.ParseUint(req.MemberID, 10, 32)
    if err != nil {
        fmt.Println("❌ Invalid member ID:", req.MemberID)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
        return
    }

    _, err = memberRepo.FindByID(uint(memberID))
    if err != nil {
        fmt.Println("❌ Member not found:", req.MemberID)
        c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
        return
    }

    // For debit transactions, check if member has sufficient balance
    if req.TransactionType == "debit" {
        savingsRepo := repositories.NewSavingsRepository()
        currentBalance, err := savingsRepo.GetMemberBalanceByCode(uint(memberID), req.AccountType)
        if err != nil {
            fmt.Println("❌ Error getting member balance:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify member balance"})
            return
        }

        if currentBalance < float64(req.Amount) {
            fmt.Printf("❌ Insufficient balance. Current: %.2f, Requested: %d\n", currentBalance, req.Amount)
            c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance for withdrawal"})
            return
        }
    }

    // Find or create saving account
    savingsRepo := repositories.NewSavingsRepository()
    savingAccount, err := savingsRepo.FindSavingAccount(uint(memberID), savingType.ID)
    if err != nil {
        // Saving account doesn't exist, create one
        fmt.Println("📝 Creating new saving account for member:", uint(memberID), "type:", savingType.Code)
        savingAccount = &models.SavingAccount{
            MemberID:      uint(memberID),
            AccountTypeID: savingType.ID,
            AccountNumber: fmt.Sprintf("SAV%08d-%s", uint(memberID), savingType.Code),
            Balance:       0,
            IsActive:      true,
        }
        if err := savingsRepo.CreateSavingAccount(savingAccount); err != nil {
            fmt.Println("❌ Error creating saving account:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create saving account"})
            return
        }
    }

    // Create transaction object
    transaction := &models.SavingTransaction{
        SavingAccountID: savingAccount.ID,
        TransactionType: req.TransactionType,
        Amount:          float64(req.Amount),
        Description:     req.Description,
        BalanceBefore:    savingAccount.Balance,
        TransactionDate: transactionDate,
        // --- PERUBAHAN 2: Ganti "system" menjadi operatorName ---
        CreatedBy:       operatorName, 
        CreatedAt:       time.Now(),
    }

    // Calculate balance after
    if req.TransactionType == "credit" {
        transaction.BalanceAfter = savingAccount.Balance + float64(req.Amount)
    } else {
        transaction.BalanceAfter = savingAccount.Balance - float64(req.Amount)
    }

    // Insert transaction into database
    fmt.Println("💾 Inserting savings transaction into database...")
    if err := savingsRepo.CreateTransaction(transaction); err != nil {
        fmt.Println("❌ Error creating savings transaction:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create savings transaction"})
        return
    }

    // Update saving account balance
    savingAccount.Balance = transaction.BalanceAfter
    if err := savingsRepo.UpdateSavingAccount(savingAccount); err != nil {
        fmt.Println("❌ Error updating saving account balance:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update saving account balance"})
        return
    }

    fmt.Println("✅ Savings transaction created successfully with ID:", transaction.ID)
    fmt.Printf("📊 Saved Transaction Data: ID=%d, SavingAccountID=%d, Amount=%.2f, Type=%s\n",
        transaction.ID, transaction.SavingAccountID, transaction.Amount, transaction.TransactionType)

    c.JSON(http.StatusCreated, gin.H{
        "message": "Savings transaction created successfully",
        "data":    transaction,
    })
}


// UpdateSavingsTransaction updates an existing savings transaction
func UpdateSavingsTransaction(c *gin.Context) {
	// --- PERUBAHAN 1: Ambil nama user dari middleware (Ditambahkan) ---
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		fmt.Println("❌ Invalid transaction ID:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	fmt.Println("📝 Updating savings transaction ID:", id)

	var req CreateSavingsTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Error binding request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing transaction
	savingsRepo := repositories.NewSavingsRepository()
	transaction, err := savingsRepo.FindTransactionByID(uint(id))
	if err != nil {
		fmt.Println("❌ Transaction not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	fmt.Println("🔍 Found existing transaction:", transaction.ID)

	// Validate saving type code
	savingTypeRepo := repositories.NewSavingTypesRepository()
	savingType, err := savingTypeRepo.FindByCode(req.AccountType)
	if err != nil {
		fmt.Println("❌ Invalid account type:", req.AccountType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account type"})
		return
	}

	// Update transaction fields
	transaction.TransactionType = req.TransactionType
	transaction.Amount = float64(req.Amount)
	transaction.Description = req.Description
	transaction.TransactionDate, _ = time.Parse("2006-01-02", req.TransactionDate)

	// --- PERUBAHAN 2: Catat siapa yang mengubah data (Ditambahkan) ---
	transaction.UpdatedBy = operatorName

	// Update member info if member ID changed
	if req.MemberID != fmt.Sprintf("%d", transaction.SavingAccount.MemberID) {
		memberID, err := strconv.ParseUint(req.MemberID, 10, 32)
		if err != nil {
			fmt.Println("❌ Invalid member ID:", req.MemberID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
			return
		}

		memberRepo := repositories.NewMemberRepository()
		_, err = memberRepo.FindByID(uint(memberID))
		if err != nil {
			fmt.Println("❌ Member not found:", req.MemberID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}

		// Find or create new saving account
		newSavingAccount, err := savingsRepo.FindSavingAccount(uint(memberID), savingType.ID)
		if err != nil {
			newSavingAccount = &models.SavingAccount{
				MemberID:      uint(memberID),
				AccountTypeID: savingType.ID,
				AccountNumber: fmt.Sprintf("SAV%08d-%s", uint(memberID), savingType.Code),
				Balance:       0,
				IsActive:      true,
			}
			savingsRepo.CreateSavingAccount(newSavingAccount)
		}

		transaction.SavingAccountID = newSavingAccount.ID
	} else {
		// Update existing saving account if it was already loaded
		if transaction.SavingAccount.ID > 0 {
			savingsRepo.UpdateSavingAccount(&transaction.SavingAccount)
		}
	}

	// Save to database
	fmt.Println("💾 Updating savings transaction in database...")
	if err := savingsRepo.UpdateTransaction(transaction); err != nil {
		fmt.Println("❌ Error updating savings transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update savings transaction"})
		return
	}

	fmt.Println("✅ Savings transaction updated successfully with ID:", transaction.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Savings transaction updated successfully",
		"data":    transaction,
	})
}


// DeleteSavingsTransaction deletes a savings transaction
func DeleteSavingsTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		fmt.Println("❌ Invalid transaction ID:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	fmt.Println("🗑️ Deleting savings transaction ID:", id)

	savingsRepo := repositories.NewSavingsRepository()
	if err := savingsRepo.DeleteTransaction(uint(id)); err != nil {
		fmt.Println("❌ Error deleting savings transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete savings transaction"})
		return
	}

	fmt.Println("✅ Savings transaction deleted successfully:", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Savings transaction deleted successfully",
	})
}

// GetMemberSavingsBalance retrieves a member's savings balance for a specific account type
func GetMemberSavingsBalance(c *gin.Context) {
	memberIDParam := c.Param("memberId")
	memberID, err := strconv.ParseUint(memberIDParam, 10, 32)
	if err != nil {
		fmt.Println("❌ Invalid member ID:", memberIDParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	accountTypeCode := c.Query("account_type") // Optional: pokok, wajib, manasuka
	fmt.Println("💰 Getting member savings balance, MemberID:", memberID, "AccountType:", accountTypeCode)

	savingsRepo := repositories.NewSavingsRepository()

	var balance float64
	if accountTypeCode != "" {
		balance, err = savingsRepo.GetMemberBalanceByCode(uint(memberID), accountTypeCode)
	} else {
		// Get all balances
		balances, err := savingsRepo.GetMemberAllBalances(uint(memberID))
		if err != nil {
			fmt.Println("❌ Error getting member balances:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve member balances"})
			return
		}

		fmt.Println("✅ Member balances retrieved:", balances)

		c.JSON(http.StatusOK, gin.H{
			"member_id": memberID,
			"balances":  balances,
		})
		return
	}

	if err != nil {
		fmt.Println("❌ Error getting member balance:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve member balance"})
		return
	}

	fmt.Println("✅ Member balance retrieved:", balance)

	c.JSON(http.StatusOK, gin.H{
		"member_id":    memberID,
		"account_type": accountTypeCode,
		"balance":      int64(balance),
	})
}

// GetAllSavingsAccounts retrieves all member balances for all account types
func GetAllSavingsAccounts(c *gin.Context) {
	fmt.Println("📋 Getting all savings accounts...")

	savingsRepo := repositories.NewSavingsRepository()

	// Get all member balances
	memberBalances, total, err := savingsRepo.GetAllMemberBalances()
	if err != nil {
		fmt.Println("❌ Error getting member balances:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve member balances"})
		return
	}

	fmt.Println("✅ Retrieved", total, "member balances")

	c.JSON(http.StatusOK, gin.H{
		"data":  memberBalances,
		"total": total,
	})
}

// GetAllSavingTypesWithBalances retrieves all saving types with their total balances
func GetAllSavingTypesWithBalances(c *gin.Context) {
	fmt.Println("📋 Getting all saving types with balances...")

	savingTypeRepo := repositories.NewSavingTypesRepository()
	savingsRepo := repositories.NewSavingsRepository()

	// Get all saving types
	savingTypes, err := savingTypeRepo.List()
	if err != nil {
		fmt.Println("❌ Error getting saving types:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve saving types"})
		return
	}

	// Get balance for each saving type
	result := make([]map[string]interface{}, len(savingTypes))
	db := savingsRepo.GetDB()
	for _, st := range savingTypes {
		// Calculate total balance for this saving type
		var totalBalance float64
		err := db.Model(&models.SavingTransaction{}).
			Joins("JOIN saving_accounts ON saving_transactions.saving_account_id = saving_accounts.id").
			Where("saving_accounts.account_type_id = ?", st.ID).
			Select("COALESCE(SUM(CASE WHEN saving_transactions.transaction_type = 'credit' THEN saving_transactions.amount ELSE -saving_transactions.amount END), 0)").
			Scan(&totalBalance).Error

		if err != nil {
			fmt.Printf("❌ Error calculating balance for type %s: %v\n", st.Code, err)
			totalBalance = 0
		}

		typeData := map[string]interface{}{
			"id":          st.ID,
			"code":        st.Code,
			"name":        st.Name,
			"description":  st.Description,
			"min_balance":  st.MinBalance,
			"balance":     int64(totalBalance),
		}
		result = append(result, typeData)
	}

	fmt.Println("✅ Retrieved", len(result), "saving types with balances")

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
