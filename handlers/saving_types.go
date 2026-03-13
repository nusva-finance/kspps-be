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

type CreateSavingTypeRequest struct {
	Code         string  `json:"code" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	IsRequired   bool    `json:"is_required"`
	MinBalance   float64 `json:"min_balance"`
	IsActive     bool    `json:"is_active"`
	DisplayOrder int     `json:"display_order"`
}

// GetSavingTypes retrieves all saving types
func GetSavingTypes(c *gin.Context) {
	includeInactive := c.Query("include_inactive") == "true"
	fmt.Println("📋 Getting saving types, include_inactive:", includeInactive)

	repo := repositories.NewSavingTypesRepository()

	var savingTypes []models.SavingType
	var err error

	if includeInactive {
		savingTypes, err = repo.ListAll()
	} else {
		savingTypes, err = repo.List()
	}

	if err != nil {
		fmt.Println("❌ Error getting saving types:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve saving types"})
		return
	}

	fmt.Println("✅ Retrieved", len(savingTypes), "saving types")

	c.JSON(http.StatusOK, gin.H{
		"data":  savingTypes,
		"total": len(savingTypes),
	})
}

// GetSavingTypeByID retrieves a single saving type by ID
func GetSavingTypeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		fmt.Println("❌ Invalid saving type ID:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid saving type ID"})
		return
	}

	fmt.Println("🔍 Getting saving type by ID:", id)

	repo := repositories.NewSavingTypesRepository()
	savingType, err := repo.FindByID(uint(id))
	if err != nil {
		fmt.Println("❌ Saving type not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Saving type not found"})
		return
	}

	fmt.Println("✅ Found saving type:", savingType.Code)

	c.JSON(http.StatusOK, savingType)
}

// CreateSavingType creates a new saving type
func CreateSavingType(c *gin.Context) {
	var req CreateSavingTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Error binding request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("📤 Creating saving type:")
	fmt.Printf("  Code: %s\n", req.Code)
	fmt.Printf("  Name: %s\n", req.Name)

	// Check if code already exists
	repo := repositories.NewSavingTypesRepository()
	existingType, err := repo.FindByCode(req.Code)
	if err == nil && existingType != nil {
		fmt.Println("❌ Saving type with code already exists:", req.Code)
		c.JSON(http.StatusConflict, gin.H{"error": "Saving type with this code already exists"})
		return
	}

	// Create saving type
	savingType := &models.SavingType{
		Code:         req.Code,
		Name:         req.Name,
		Description:   req.Description,
		IsRequired:    req.IsRequired,
		MinBalance:    req.MinBalance,
		IsActive:      req.IsActive,
		DisplayOrder:  req.DisplayOrder,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := repo.Create(savingType); err != nil {
		fmt.Println("❌ Error creating saving type:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create saving type"})
		return
	}

	fmt.Println("✅ Saving type created successfully with ID:", savingType.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Saving type created successfully",
		"data":    savingType,
	})
}

// UpdateSavingType updates an existing saving type
func UpdateSavingType(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		fmt.Println("❌ Invalid saving type ID:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid saving type ID"})
		return
	}

	fmt.Println("📝 Updating saving type ID:", id)

	var req CreateSavingTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Error binding request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo := repositories.NewSavingTypesRepository()
	savingType, err := repo.FindByID(uint(id))
	if err != nil {
		fmt.Println("❌ Saving type not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Saving type not found"})
		return
	}

	fmt.Println("🔍 Found existing saving type:", savingType.Code)

	// Check if code conflict with other saving types
	existingType, _ := repo.FindByCode(req.Code)
	if existingType != nil && existingType.ID != uint(id) {
		fmt.Println("❌ Saving type code already used by another type:", req.Code)
		c.JSON(http.StatusConflict, gin.H{"error": "Saving type code already used by another type"})
		return
	}

	// Update saving type fields
	savingType.Code = req.Code
	savingType.Name = req.Name
	savingType.Description = req.Description
	savingType.IsRequired = req.IsRequired
	savingType.MinBalance = req.MinBalance
	savingType.IsActive = req.IsActive
	savingType.DisplayOrder = req.DisplayOrder
	savingType.UpdatedAt = time.Now()

	// Save to database
	fmt.Println("💾 Updating saving type in database...")
	if err := repo.Update(savingType); err != nil {
		fmt.Println("❌ Error updating saving type:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update saving type"})
		return
	}

	fmt.Println("✅ Saving type updated successfully with ID:", savingType.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Saving type updated successfully",
		"data":    savingType,
	})
}

// DeleteSavingType deletes a saving type
func DeleteSavingType(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		fmt.Println("❌ Invalid saving type ID:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid saving type ID"})
		return
	}

	fmt.Println("🗑️ Deleting saving type ID:", id)

	repo := repositories.NewSavingTypesRepository()

	// Check if saving type exists
	_, err = repo.FindByID(uint(id))
	if err != nil {
		fmt.Println("❌ Saving type not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Saving type not found"})
		return
	}

	// Check if there are any saving accounts using this type
	savingsRepo := repositories.NewSavingsRepository()
	accounts, err := savingsRepo.FindSavingAccount(0, uint(id))
	if err == nil && accounts != nil {
		fmt.Println("❌ Cannot delete saving type, accounts exist:", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete saving type with existing accounts"})
		return
	}

	if err := repo.Delete(uint(id)); err != nil {
		fmt.Println("❌ Error deleting saving type:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete saving type"})
		return
	}

	fmt.Println("✅ Saving type deleted successfully:", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Saving type deleted successfully",
	})
}

// InitializeSavingTypes initializes default saving types
func InitializeSavingTypes(c *gin.Context) {
	fmt.Println("🚀 Initializing default saving types...")

	repo := repositories.NewSavingTypesRepository()
	if err := repo.InitializeDefaultTypes(); err != nil {
		fmt.Println("❌ Error initializing saving types:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize saving types"})
		return
	}

	fmt.Println("✅ Default saving types initialized successfully")

	// Return the initialized types
	savingTypes, err := repo.List()
	if err != nil {
		fmt.Println("❌ Error getting saving types:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve saving types"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Default saving types initialized successfully",
		"data":    savingTypes,
	})
}
