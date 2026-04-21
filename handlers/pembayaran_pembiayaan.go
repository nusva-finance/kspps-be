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

type CreatePembayaranRequest struct {
	IDPinjaman            int     `json:"idpinjaman" binding:"required"`
	IDRekening            int     `json:"idrekening" binding:"required"`
	TglPembayaran         string  `json:"tglpembayaran" binding:"required"`
	NominalPembayaran     float64 `json:"nominalpembayaran" binding:"required"`
	NominalAngsuran       float64 `json:"nominalangsuran" binding:"required"`
	NominalPendapatanLain float64 `json:"nominalpendapatanlainlain"`
	AngsuranKe            int     `json:"angsuranke" binding:"required"`
	TglJtAngsuran         string  `json:"tgljtangsuran" binding:"required"`
	Keterangan            string  `json:"keterangan"`
}

// GetPembayaranByID retrieves a single payment by ID
func GetPembayaranByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println("❌ Invalid pembayaran ID:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pembayaran ID"})
		return
	}

	fmt.Println("📋 Getting pembayaran by ID:", id)

	repo := repositories.NewPembayaranPembiayaanRepository()
	pembayaran, err := repo.FindByID(id)
	if err != nil {
		fmt.Println("❌ Error getting pembayaran:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Pembayaran not found"})
		return
	}

	fmt.Println("✅ Retrieved pembayaran ID:", id)

	c.JSON(http.StatusOK, gin.H{
		"data": pembayaran,
	})
}

// GetPembayaranByPinjamanID retrieves all payments for a specific loan
func GetPembayaranByPinjamanID(c *gin.Context) {
	idParam := c.Param("id")
	idPinjaman, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println("❌ Invalid pinjaman ID:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pinjaman ID"})
		return
	}

	fmt.Println("📋 Getting pembayaran for pinjaman ID:", idPinjaman)

	repo := repositories.NewPembayaranPembiayaanRepository()
	pembayaran, err := repo.ListByPinjamanID(idPinjaman)
	if err != nil {
		fmt.Println("❌ Error getting pembayaran:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pembayaran"})
		return
	}

	fmt.Println("✅ Retrieved", len(pembayaran), "pembayaran records")

	c.JSON(http.StatusOK, gin.H{
		"data": pembayaran,
	})
}

// GetAngsuranKe calculates the next installment number for a loan
func GetAngsuranKe(c *gin.Context) {
	idParam := c.Param("id")
	idPinjaman, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println("❌ Invalid pinjaman ID:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pinjaman ID"})
		return
	}

	repo := repositories.NewPembayaranPembiayaanRepository()
	count, err := repo.CountByPinjamanID(idPinjaman)
	if err != nil {
		fmt.Println("❌ Error counting pembayaran:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count pembayaran"})
		return
	}

	angsuranKe := int(count) + 1

	fmt.Printf("📊 AngsuranKe for pinjaman %d: %d (existing records: %d)\n", idPinjaman, angsuranKe, count)

	c.JSON(http.StatusOK, gin.H{
		"angsuranke":   angsuranKe,
		"total_record": count,
	})
}

// CreatePembayaran creates a new payment record
func CreatePembayaran(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	var req CreatePembayaranRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Error binding request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("📤 Creating pembayaran:")
	fmt.Printf("  IDPinjaman: %d\n", req.IDPinjaman)
	fmt.Printf("  IDRekening: %d\n", req.IDRekening)
	fmt.Printf("  TglPembayaran: %s\n", req.TglPembayaran)
	fmt.Printf("  NominalPembayaran: %.2f\n", req.NominalPembayaran)
	fmt.Printf("  NominalAngsuran: %.2f\n", req.NominalAngsuran)
	fmt.Printf("  NominalPendapatanLain: %.2f\n", req.NominalPendapatanLain)
	fmt.Printf("  AngsuranKe: %d\n", req.AngsuranKe)
	fmt.Printf("  TglJtAngsuran: %s\n", req.TglJtAngsuran)
	fmt.Printf("  Keterangan: %s\n", req.Keterangan)

	// Parse dates
	tglPembayaran, err := time.Parse("2006-01-02", req.TglPembayaran)
	if err != nil {
		fmt.Println("❌ Invalid tglpembayaran format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tglpembayaran format. Use YYYY-MM-DD"})
		return
	}

	tglJtAngsuran, err := time.Parse("2006-01-02", req.TglJtAngsuran)
	if err != nil {
		fmt.Println("❌ Invalid tgljtangsuran format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tgljtangsuran format. Use YYYY-MM-DD"})
		return
	}

	// Verify pinjaman exists
	pembiayaanRepo := repositories.NewPembiayaanRepository()
	_, err = pembiayaanRepo.FindByID(req.IDPinjaman)
	if err != nil {
		fmt.Println("❌ Pinjaman not found:", req.IDPinjaman)
		c.JSON(http.StatusNotFound, gin.H{"error": "Pinjaman not found"})
		return
	}

	// Create pembayaran record
	pembayaran := &models.PembayaranPembiayaan{
		IDPinjaman:            req.IDPinjaman,
		TglPembayaran:         tglPembayaran,
		NominalPembayaran:     req.NominalPembayaran,
		NominalAngsuran:       req.NominalAngsuran,
		NominalPendapatanLain: req.NominalPendapatanLain,
		AngsuranKe:            req.AngsuranKe,
		TglJtAngsuran:         tglJtAngsuran,
		Keterangan:            req.Keterangan,
		SysRevID:              1,
		CreatedBy:             operatorName,
		UpdatedBy:             operatorName,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	repo := repositories.NewPembayaranPembiayaanRepository()
	if err := repo.Create(pembayaran); err != nil {
		fmt.Println("❌ Error creating pembayaran:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pembayaran"})
		return
	}

	fmt.Println("✅ Pembayaran created successfully with ID:", pembayaran.IDPembayaranPembiayaan)

	// Create rekening_transaction record
	rekeningTransaction := &models.RekeningTransaction{
		TransactionType:     "INSERT",
		IDNusvaRekening:     uint(req.IDRekening),
		TableTransaction:    "pembayaran_pembiayaan",
		IDTableTransaction:  uint(pembayaran.IDPembayaranPembiayaan),
		TanggalTransaksi:    tglPembayaran,
		NominalTransaksi:    req.NominalPembayaran,
		CreatedBy:           operatorName,
		UpdatedBy:           operatorName,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := repo.CreateRekeningTransaction(rekeningTransaction); err != nil {
		fmt.Println("❌ Error creating rekening_transaction:", err)
		// Don't fail the request, just log the error
	} else {
		fmt.Println("✅ Rekening transaction created successfully")
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pembayaran created successfully",
		"data":    pembayaran,
	})
}

// UpdatePembayaran updates an existing payment record
func UpdatePembayaran(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println("❌ Invalid pembayaran ID:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pembayaran ID"})
		return
	}

	var req CreatePembayaranRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Error binding request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates
	tglPembayaran, err := time.Parse("2006-01-02", req.TglPembayaran)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tglpembayaran format"})
		return
	}

	tglJtAngsuran, err := time.Parse("2006-01-02", req.TglJtAngsuran)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tgljtangsuran format"})
		return
	}

	repo := repositories.NewPembayaranPembiayaanRepository()
	pembayaran, err := repo.FindByID(id)
	if err != nil {
		fmt.Println("❌ Pembayaran not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Pembayaran not found"})
		return
	}

	// Get old rekening_transaction record for REVERSE
	oldTransaction, err := repo.FindRekeningTransactionByReference("pembayaran_pembiayaan", uint(id))
	if err != nil {
		fmt.Println("⚠️ No existing rekening_transaction found for REVERSE:", err)
		// Continue without REVERSE record if not found
	} else {
		// Create REVERSE record (copy of old data with negative nominal)
		reverseTransaction := &models.RekeningTransaction{
			TransactionType:     "REVERSE",
			IDNusvaRekening:     oldTransaction.IDNusvaRekening,
			TableTransaction:    "pembayaran_pembiayaan",
			IDTableTransaction:  uint(id),
			TanggalTransaksi:    oldTransaction.TanggalTransaksi,
			NominalTransaksi:    -oldTransaction.NominalTransaksi, // Negative for reversal
			CreatedBy:           operatorName,
			UpdatedBy:           operatorName,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		if err := repo.CreateRekeningTransaction(reverseTransaction); err != nil {
			fmt.Println("❌ Error creating REVERSE rekening_transaction:", err)
		} else {
			fmt.Println("✅ REVERSE rekening transaction created successfully")
		}
	}

	// Update pembayaran fields
	pembayaran.TglPembayaran = tglPembayaran
	pembayaran.NominalPembayaran = req.NominalPembayaran
	pembayaran.NominalAngsuran = req.NominalAngsuran
	pembayaran.NominalPendapatanLain = req.NominalPendapatanLain
	pembayaran.AngsuranKe = req.AngsuranKe
	pembayaran.TglJtAngsuran = tglJtAngsuran
	pembayaran.Keterangan = req.Keterangan
	pembayaran.UpdatedBy = operatorName
	pembayaran.UpdatedAt = time.Now()

	if err := repo.Update(pembayaran); err != nil {
		fmt.Println("❌ Error updating pembayaran:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pembayaran"})
		return
	}

	fmt.Println("✅ Pembayaran updated successfully:", id)

	// Create UPDATE record with new data
	updateTransaction := &models.RekeningTransaction{
		TransactionType:     "UPDATE",
		IDNusvaRekening:     uint(req.IDRekening),
		TableTransaction:    "pembayaran_pembiayaan",
		IDTableTransaction:  uint(id),
		TanggalTransaksi:    tglPembayaran,
		NominalTransaksi:    req.NominalPembayaran,
		CreatedBy:           operatorName,
		UpdatedBy:           operatorName,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := repo.CreateRekeningTransaction(updateTransaction); err != nil {
		fmt.Println("❌ Error creating UPDATE rekening_transaction:", err)
	} else {
		fmt.Println("✅ UPDATE rekening transaction created successfully")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pembayaran updated successfully",
		"data":    pembayaran,
	})
}


// DeletePembayaran deletes a payment record
func DeletePembayaran(c *gin.Context) {
	// Ambil operator untuk audit
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println("❌ Invalid pembayaran ID:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pembayaran ID"})
		return
	}

	fmt.Println("🗑️ Deleting pembayaran ID:", id)

	repo := repositories.NewPembayaranPembiayaanRepository()

	// 1. Ambil data pembayaran SEBELUM dihapus untuk mendapatkan Nominal dan data aslinya
	pembayaran, err := repo.FindByID(id)
	if err != nil {
		fmt.Println("❌ Pembayaran not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Pembayaran not found"})
		return
	}

	// 2. Cari transaksi rekening awal untuk mendapatkan IDNusvaRekening
	oldTransaction, err := repo.FindRekeningTransactionByReference("pembayaran_pembiayaan", uint(id))
	var idNusvaRekening uint = 0
	if err == nil && oldTransaction != nil {
		idNusvaRekening = oldTransaction.IDNusvaRekening
	}

	// 3. Hapus pembayaran dari tabel pembayaran_pembiayaan
	if err := repo.Delete(id); err != nil {
		fmt.Println("❌ Error deleting pembayaran:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pembayaran"})
		return
	}

	fmt.Println("✅ Pembayaran deleted successfully:", id)

	// 4. Buat record rekening_transaction dengan tipe "Delete"
	if idNusvaRekening > 0 {
		deleteTransaction := &models.RekeningTransaction{
			TransactionType:    "Delete",
			IDNusvaRekening:    idNusvaRekening,
			TableTransaction:   "pembayaran_pembiayaan",
			IDTableTransaction: uint(id),
			TanggalTransaksi:   time.Now(),
			// Nominal kita buat negatif (-) untuk mengurangi/membalikkan uang yang sebelumnya masuk (Insert)
			NominalTransaksi:   -pembayaran.NominalPembayaran,
			CreatedBy:          operatorName,
			UpdatedBy:          operatorName,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}

		if err := repo.CreateRekeningTransaction(deleteTransaction); err != nil {
			fmt.Println("❌ Error creating DELETE rekening_transaction:", err)
		} else {
			fmt.Println("✅ DELETE rekening transaction created successfully")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pembayaran deleted successfully",
	})
}
