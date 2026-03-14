package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"nusvakspps/models"
	"nusvakspps/repositories"
)

type CreateMarginRequest struct {
	Category string  `json:"category" binding:"required"`
	Tenor    int     `json:"tenor" binding:"required"`
	Margin   float64 `json:"margin" binding:"required"`
}

type UpdateMarginRequest struct {
	Category string  `json:"category"`
	Tenor    int     `json:"tenor"`
	Margin   float64 `json:"margin"`
}

// GetMargins - Mengambil daftar pengaturan margin
func GetMargins(c *gin.Context) {
	fmt.Println("📋 Getting all margins...")

	repo := repositories.NewMarginRepository()
	margins, _, err := repo.List(0, 1000)
	if err != nil {
		fmt.Println("❌ Error getting margins:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data margin"})
		return
	}

	fmt.Printf("✅ Retrieved %d margins\n", len(margins))
	c.JSON(http.StatusOK, gin.H{"data": margins})
}

// GetMarginByID - Mengambil pengaturan margin berdasarkan ID
func GetMarginByID(c *gin.Context) {
	idParam := c.Param("id")
	fmt.Println("🔍 Getting margin by ID:", idParam)

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	repo := repositories.NewMarginRepository()
	margin, err := repo.FindByID(id)
	if err != nil {
		fmt.Println("⚠️ Margin not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Margin tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": margin})
}

// CreateMargin - Membuat pengaturan margin baru
func CreateMargin(c *gin.Context) {
	// AMBIL OPERATOR DARI MIDDLEWARE
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	var req CreateMarginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Validation error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data belum lengkap: " + err.Error()})
		return
	}

	fmt.Println("📤 Creating new margin:")
	fmt.Println("   Category:", req.Category)
	fmt.Println("   Tenor:", req.Tenor)
	fmt.Println("   Margin:", req.Margin)
	fmt.Println("   Operator:", operatorName)

	repo := repositories.NewMarginRepository()

	margin := &models.MarginSetup{
		Category:  req.Category,
		Tenor:     req.Tenor,
		Margin:    req.Margin,
		CreatedBy: operatorName,
		UpdatedBy: operatorName,
	}

	fmt.Println("📥 About to create margin object:", margin)

	if err := repo.Create(margin); err != nil {
		fmt.Println("❌ Error creating margin:")
		fmt.Println("   Error:", err.Error())
		fmt.Println("   Error Type:", fmt.Sprintf("%T", err))

		// Check if it's a duplicate constraint error
		errMsg := err.Error()
		if len(errMsg) > 0 && (strings.Contains(strings.ToLower(errMsg), "duplicate") || strings.Contains(strings.ToLower(errMsg), "unique") || strings.Contains(strings.ToLower(errMsg), "constraint")) {
			fmt.Println("⚠️ Duplicate constraint detected")
			c.JSON(http.StatusConflict, gin.H{"error": "Kombinasi kategori dan tenor sudah ada"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data margin: " + err.Error()})
		return
	}

	fmt.Println("✅ Margin created successfully:")
	fmt.Println("   ID:", margin.ID)
	fmt.Println("   Category:", margin.Category)
	fmt.Println("   Tenor:", margin.Tenor)
	fmt.Println("   Margin:", margin.Margin)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Margin berhasil ditambahkan",
		"data":    margin,
	})
}

// UpdateMargin - Mengubah pengaturan margin
func UpdateMargin(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	idParam := c.Param("id")

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	repo := repositories.NewMarginRepository()
	margin, err := repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Margin tidak ditemukan"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Update fields from input
	if val, ok := input["category"].(string); ok {
		margin.Category = val
	}
	if val, ok := input["tenor"].(float64); ok {
		margin.Tenor = int(val)
	}
	if val, ok := input["margin"].(float64); ok {
		margin.Margin = val
	}

	// Update audit fields
	margin.UpdatedBy = operatorName

	if err := repo.Update(margin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diperbarui",
		"data":    margin,
	})
}

// DeleteMargin - Menghapus pengaturan margin
func DeleteMargin(c *gin.Context) {
	idParam := c.Param("id")
	fmt.Println("🗑️ Deleting margin ID:", idParam)

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	repo := repositories.NewMarginRepository()
	if err := repo.Delete(id); err != nil {
		fmt.Println("❌ Error deleting margin:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data margin"})
		return
	}

	fmt.Println("✅ Margin deleted successfully, ID:", id)

	c.JSON(http.StatusOK, gin.H{"message": "Margin berhasil dihapus"})
}
