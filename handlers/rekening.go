package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"nusvakspps/models"
	"nusvakspps/repositories"
)

type CreateRekeningRequest struct {
	NamaRekening string `json:"namarekening" binding:"required"`
	NoRekening    string `json:"norekening" binding:"required"`
	Aktif         bool   `json:"aktif"`
	Deskripsi     string `json:"deskripsi"`
}

type UpdateRekeningRequest struct {
	NamaRekening string `json:"namarekening"`
	NoRekening    string `json:"norekening"`
	Aktif         bool   `json:"aktif"`
	Deskripsi     string `json:"deskripsi"`
}

// GetRekenings - Mengambil daftar rekening
func GetRekenings(c *gin.Context) {
	fmt.Println("📋 Getting all rekenings...")

	repo := repositories.NewNusvaRekeningRepository()
	rekenings, _, err := repo.List(0, 1000)
	if err != nil {
		fmt.Println("❌ Error getting rekenings:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data rekening"})
		return
	}

	fmt.Printf("✅ Retrieved %d rekenings\n", len(rekenings))
	c.JSON(http.StatusOK, gin.H{"data": rekenings})
}

// GetRekeningByID - Mengambil rekening berdasarkan ID
func GetRekeningByID(c *gin.Context) {
	idParam := c.Param("id")
	fmt.Println("🔍 Getting rekening by ID:", idParam)

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	repo := repositories.NewNusvaRekeningRepository()
	rekening, err := repo.FindByID(id)
	if err != nil {
		fmt.Println("⚠️ Rekening not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Rekening tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rekening})
}

// CreateRekening - Membuat rekening baru
func CreateRekening(c *gin.Context) {
	// AMBIL OPERATOR DARI MIDDLEWARE
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	var req CreateRekeningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Validation error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data belum lengkap: " + err.Error()})
		return
	}

	fmt.Println("📤 Creating new rekening:")
	fmt.Println("   Nama Rekening:", req.NamaRekening)
	fmt.Println("   No Rekening:", req.NoRekening)
	fmt.Println("   Aktif:", req.Aktif)
	fmt.Println("   Deskripsi:", req.Deskripsi)
	fmt.Println("   Operator:", operatorName)

	repo := repositories.NewNusvaRekeningRepository()

	rekening := &models.NusvaRekening{
		NamaRekening: req.NamaRekening,
		NoRekening:    req.NoRekening,
		Aktif:         req.Aktif,
		Deskripsi:     req.Deskripsi,
		CreatedBy:     operatorName,
		UpdatedBy:     operatorName,
	}

	if !req.Aktif {
		rekening.Aktif = true
	}

	fmt.Println("📥 About to create rekening object:", rekening)

	if err := repo.Create(rekening); err != nil {
		fmt.Println("❌ Error creating rekening:")
		fmt.Println("   Error:", err.Error())
		fmt.Println("   Error Type:", fmt.Sprintf("%T", err))

		// Check if it's a duplicate constraint error
		errMsg := err.Error()
		if len(errMsg) > 0 && (strings.Contains(strings.ToLower(errMsg), "duplicate") || strings.Contains(strings.ToLower(errMsg), "unique") || strings.Contains(strings.ToLower(errMsg), "constraint")) {
			fmt.Println("⚠️ Duplicate constraint detected")
			c.JSON(http.StatusConflict, gin.H{"error": "Nomor rekening sudah ada"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data rekening: " + err.Error()})
		return
	}

	fmt.Println("✅ Rekening created successfully:")
	fmt.Println("   ID:", rekening.IDRekening)
	fmt.Println("   Nama Rekening:", rekening.NamaRekening)
	fmt.Println("   No Rekening:", rekening.NoRekening)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Rekening berhasil ditambahkan",
		"data":    rekening,
	})
}

// UpdateRekening - Mengubah data rekening
func UpdateRekening(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	idParam := c.Param("id")

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	repo := repositories.NewNusvaRekeningRepository()
	rekening, err := repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rekening tidak ditemukan"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Update fields from input
	if val, ok := input["namarekening"].(string); ok {
		rekening.NamaRekening = val
	}
	if val, ok := input["norekening"].(string); ok {
		rekening.NoRekening = val
	}
	if val, ok := input["aktif"].(bool); ok {
		rekening.Aktif = val
	}
	if val, ok := input["deskripsi"].(string); ok {
		rekening.Deskripsi = val
	}

	// Update audit fields
	rekening.UpdatedBy = operatorName

	if err := repo.Update(rekening); err != nil {
		fmt.Println("❌ Error updating rekening:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	fmt.Println("✅ Rekening updated successfully:")
	fmt.Println("   ID:", rekening.IDRekening)

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diperbarui",
		"data":    rekening,
	})
}

// DeleteRekening - Menghapus rekening
func DeleteRekening(c *gin.Context) {
	idParam := c.Param("id")
	fmt.Println("🗑️ Deleting rekening ID:", idParam)

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	repo := repositories.NewNusvaRekeningRepository()
	if err := repo.Delete(id); err != nil {
		fmt.Println("❌ Error deleting rekening:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data rekening"})
		return
	}

	fmt.Println("✅ Rekening deleted successfully, ID:", id)

	c.JSON(http.StatusOK, gin.H{"message": "Rekening berhasil dihapus"})
}
