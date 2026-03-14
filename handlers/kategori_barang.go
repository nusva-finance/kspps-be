package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"nusvakspps/models"
	"nusvakspps/repositories"
)

type CreateKategoriBarangRequest struct {
	Kategori string `json:"kategori" binding:"required"`
	Aktif   bool   `json:"aktif"`
}

type UpdateKategoriBarangRequest struct {
	Kategori string `json:"kategori"`
	Aktif   bool   `json:"aktif"`
}

// GetKategoriBarangs - Mengambil daftar kategori barang
func GetKategoriBarangs(c *gin.Context) {
	fmt.Println("📋 Getting all kategori barang...")

	repo := repositories.NewKategoriBarangRepository()
	kategoriBarangs, _, err := repo.List(0, 1000)
	if err != nil {
		fmt.Println("❌ Error getting kategori barang:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kategori barang"})
		return
	}

	fmt.Printf("✅ Retrieved %d kategori barang\n", len(kategoriBarangs))
	c.JSON(http.StatusOK, gin.H{"data": kategoriBarangs})
}

// GetKategoriBarangByID - Mengambil kategori barang berdasarkan ID
func GetKategoriBarangByID(c *gin.Context) {
	idParam := c.Param("id")
	fmt.Println("🔍 Getting kategori barang by ID:", idParam)

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	repo := repositories.NewKategoriBarangRepository()
	kategoriBarang, err := repo.FindByID(id)
	if err != nil {
		fmt.Println("⚠️ Kategori barang not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori barang tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kategoriBarang})
}

// CreateKategoriBarang - Membuat kategori barang baru
func CreateKategoriBarang(c *gin.Context) {
	// AMBIL OPERATOR DARI MIDDLEWARE
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	var req CreateKategoriBarangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Validation error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data belum lengkap: " + err.Error()})
		return
	}

	fmt.Println("📤 Creating new kategori barang:", req.Kategori)

	repo := repositories.NewKategoriBarangRepository()

	// Check if category already exists
	existing, _ := repo.FindByKategori(req.Kategori)
	if existing != nil {
		fmt.Println("⚠️ Kategori already exists:", req.Kategori)
		c.JSON(http.StatusConflict, gin.H{"error": "Kategori sudah ada"})
		return
	}

	kategoriBarang := &models.KategoriBarang{
		Kategori:  req.Kategori,
		Aktif:     req.Aktif,
		CreatedBy: operatorName,
		UpdatedBy: operatorName,
	}

	if err := repo.Create(kategoriBarang); err != nil {
		fmt.Println("❌ Error creating kategori barang:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data kategori barang"})
		return
	}

	fmt.Println("✅ Kategori barang created successfully:", kategoriBarang.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kategori barang berhasil ditambahkan",
		"data":    kategoriBarang,
	})
}

// UpdateKategoriBarang - Mengubah kategori barang
func UpdateKategoriBarang(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	idParam := c.Param("id")

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	repo := repositories.NewKategoriBarangRepository()
	kategoriBarang, err := repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori barang tidak ditemukan"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Update fields from input
	if val, ok := input["kategori"].(string); ok {
		kategoriBarang.Kategori = val
	}
	if val, ok := input["aktif"].(bool); ok {
		kategoriBarang.Aktif = val
	}

	// Update audit fields
	kategoriBarang.UpdatedBy = operatorName

	if err := repo.Update(kategoriBarang); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diperbarui",
		"data":    kategoriBarang,
	})
}

// DeleteKategoriBarang - Menghapus kategori barang
func DeleteKategoriBarang(c *gin.Context) {
	idParam := c.Param("id")
	fmt.Println("🗑️ Deleting kategori barang ID:", idParam)

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	repo := repositories.NewKategoriBarangRepository()
	if err := repo.Delete(id); err != nil {
		fmt.Println("❌ Error deleting kategori barang:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data kategori barang"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kategori barang berhasil dihapus"})
}
