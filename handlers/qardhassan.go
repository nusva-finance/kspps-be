package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"nusvakspps/config"
	"nusvakspps/models"
	"nusvakspps/repositories"
)

type CreateQardHassanRequest struct {
	IDMember        int     `json:"idmember" binding:"required"`
	IDRekening      int     `json:"idrekening" binding:"required"`
	TanggalPinjaman string  `json:"tanggalpinjaman" binding:"required"`
	BiayaAdmin      float64 `json:"biayaadmin"`
	NominalPinjaman float64 `json:"nominalpinjaman" binding:"required"`
	TglJtTempo      string  `json:"tgljttempo" binding:"required"`
	Keterangan      string  `json:"keterangan"`
}

// GetQardHassan retrieves all qardhassan records
func GetQardHassan(c *gin.Context) {
	fmt.Println("📋 Getting all qardhassan")

	repo := repositories.NewQardHassanRepository()
	qardhassan, total, err := repo.List(0, 1000)
	if err != nil {
		fmt.Println("❌ Error getting qardhassan:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve qardhassan"})
		return
	}

	fmt.Println("✅ Retrieved", len(qardhassan), "qardhassan records")

	c.JSON(http.StatusOK, gin.H{
		"data":  qardhassan,
		"total": total,
	})
}

// GetQardHassanByID retrieves a single qardhassan by ID
func GetQardHassanByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid qardhassan ID"})
		return
	}

	repo := repositories.NewQardHassanRepository()
	qardhassan, err := repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QardHassan not found"})
		return
	}

	// Cari idnusvarekening terakhir untuk transaksi ini
	var oldTx models.RekeningTransaction
	config.GetDB().Where("tabletransaction = 'qardhassan' AND idtabletransaction = ?", id).
		Order("created_at DESC").First(&oldTx)

	// Cari rekening yang digunakan untuk pembayaran (transaksi dengan nominal positif = uang masuk)
	var payTx models.RekeningTransaction
	config.GetDB().Where("tabletransaction = 'qardhassan' AND idtabletransaction = ? AND nominaltransaction > 0", id).
		Order("created_at DESC").First(&payTx)

	response := map[string]interface{}{
		"idqardhassan":       qardhassan.IDQardHassan,
		"idmember":           qardhassan.IDMember,
		"idnusvarekening":    oldTx.IDNusvaRekening,
		"tanggalpinjaman":    qardhassan.TanggalPinjaman,
		"biayaadmin":         qardhassan.BiayaAdmin,
		"nominalpinjaman":    qardhassan.NominalPinjaman,
		"tgljttempo":         qardhassan.TglJtTempo,
		"keterangan":         qardhassan.Keterangan,
		"nominalpembayaran":  qardhassan.NominalPembayaran,
		"tanggalpembayaran":  qardhassan.TanggalPembayaran,
		"idrekeningbayar":    payTx.IDNusvaRekening,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// CreateQardHassan creates a new qardhassan record
// CreateQardHassan creates a new qardhassan record
func CreateQardHassan(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	var req CreateQardHassanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates
	tanggalPinjaman, err := time.Parse("2006-01-02", req.TanggalPinjaman)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tanggalpinjaman format. Use YYYY-MM-DD"})
		return
	}

	tglJtTempo, err := time.Parse("2006-01-02", req.TglJtTempo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tgljttempo format. Use YYYY-MM-DD"})
		return
	}

	// ==========================================
	// TAMBAHAN: BLOK VALIDASI PLAFON QARD HASSAN
	// ==========================================
	
	// 1. Ambil data Plafon Member dari tabel Members
	memberRepo := repositories.NewMemberRepository()
	memberData, err := memberRepo.FindByID(uint(req.IDMember))
	if err != nil || memberData == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data anggota tidak ditemukan"})
		return
	}

	// 2. Ambil total sisa hutang (Outstanding) Qard Hassan anggota ini
	qardRepo := repositories.NewQardHassanRepository()
	outstandingLama, err := qardRepo.GetOutstandingByMemberID(req.IDMember)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung outstanding pinjaman"})
		return
	}

	// 3. Hitung total pinjaman jika pengajuan baru ini disetujui
	totalPengajuan := outstandingLama + req.NominalPinjaman

	// 4. Validasi: Apakah melebihi Plafon? (Dilewati jika Plafon 0 / Unlimited)
	if memberData.QardhassanPlafon > 0 {
		if totalPengajuan > memberData.QardhassanPlafon {
			pesanError := fmt.Sprintf(
				"Limit Plafon tidak mencukupi. Plafon: Rp %.0f | Sisa Hutang: Rp %.0f | Pengajuan Baru: Rp %.0f", 
				memberData.QardhassanPlafon, outstandingLama, req.NominalPinjaman,
			)
			c.JSON(http.StatusBadRequest, gin.H{"error": pesanError})
			return
		}
	}
	// ==========================================


	// Create qardhassan record (Jika lolos validasi, lanjut simpan ke DB)
	qardhassan := &models.QardHassan{
		IDMember:        req.IDMember,
		TanggalPinjaman: tanggalPinjaman,
		BiayaAdmin:      req.BiayaAdmin,
		NominalPinjaman: req.NominalPinjaman,
		TglJtTempo:      tglJtTempo,
		Keterangan:      req.Keterangan,
		SysRevID:        1,
		CreatedBy:       operatorName,
		UpdatedBy:       operatorName,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	repo := repositories.NewQardHassanRepository() // catatan: Kita inisialisasi ulang repo karena variabel qardRepo di atas beda
	if err := repo.Create(qardhassan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create qardhassan"})
		return
	}

	// Create rekening_transaction record
	rekeningTransaction := &models.RekeningTransaction{
		TransactionType:    "Insert",
		IDNusvaRekening:    uint(req.IDRekening),
		TableTransaction:   "qardhassan",
		IDTableTransaction: uint(qardhassan.IDQardHassan),
		TanggalTransaksi:   tanggalPinjaman,
		NominalTransaksi:   -1 * req.NominalPinjaman, // Uang Keluar (Pinjaman) = Minus
		CreatedBy:          operatorName,
		UpdatedBy:          operatorName,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := repo.CreateRekeningTransaction(rekeningTransaction); err != nil {
		fmt.Println("❌ Error creating rekening_transaction:", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "QardHassan created successfully",
		"data":    qardhassan,
	})
}

// UpdateQardHassan updates an existing qardhassan record
func UpdateQardHassan(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid qardhassan ID"})
		return
	}

	var req CreateQardHassanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates
	tanggalPinjaman, err := time.Parse("2006-01-02", req.TanggalPinjaman)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tanggalpinjaman format"})
		return
	}

	tglJtTempo, err := time.Parse("2006-01-02", req.TglJtTempo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tgljttempo format"})
		return
	}

	repo := repositories.NewQardHassanRepository()
	qardhassan, err := repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QardHassan not found"})
		return
	}

	// Simpan nominal lama untuk Reversal
	oldNominalPinjaman := qardhassan.NominalPinjaman

	// Cari rekening lama untuk Reversal
	var oldTx models.RekeningTransaction
	config.GetDB().Where("tabletransaction = 'qardhassan' AND idtabletransaction = ?", id).
		Order("created_at DESC").First(&oldTx)
	oldIDNusvaRekening := oldTx.IDNusvaRekening

	// Update fields
	qardhassan.IDMember = req.IDMember
	qardhassan.TanggalPinjaman = tanggalPinjaman
	qardhassan.BiayaAdmin = req.BiayaAdmin
	qardhassan.NominalPinjaman = req.NominalPinjaman
	qardhassan.TglJtTempo = tglJtTempo
	qardhassan.Keterangan = req.Keterangan
	qardhassan.UpdatedBy = operatorName
	qardhassan.UpdatedAt = time.Now()

	if err := repo.Update(qardhassan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update qardhassan"})
		return
	}

	// 1. Buat record REVERSAL (Membatalkan yang lama)
	if oldIDNusvaRekening > 0 {
		reverseTx := &models.RekeningTransaction{
			TransactionType:    "Reverse",
			IDNusvaRekening:    oldIDNusvaRekening,
			TableTransaction:   "qardhassan",
			IDTableTransaction: uint(id),
			TanggalTransaksi:   time.Now(),
			NominalTransaksi:   oldNominalPinjaman, // Positif untuk membalikkan nilai negatif sebelumnya
			CreatedBy:          operatorName,
			UpdatedBy:          operatorName,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		repo.CreateRekeningTransaction(reverseTx)
	}

	// 2. Buat record UPDATE (Mencatat yang baru)
	updateTx := &models.RekeningTransaction{
		TransactionType:    "Update",
		IDNusvaRekening:    uint(req.IDRekening),
		TableTransaction:   "qardhassan",
		IDTableTransaction: uint(id),
		TanggalTransaksi:   tanggalPinjaman,
		NominalTransaksi:   -1 * req.NominalPinjaman, // Negatif karena uang keluar
		CreatedBy:          operatorName,
		UpdatedBy:          operatorName,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	repo.CreateRekeningTransaction(updateTx)

	c.JSON(http.StatusOK, gin.H{
		"message": "QardHassan updated successfully",
		"data":    qardhassan,
	})
}

// DeleteQardHassan deletes a qardhassan record
func DeleteQardHassan(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid qardhassan ID"})
		return
	}

	repo := repositories.NewQardHassanRepository()
	qardhassan, err := repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QardHassan not found"})
		return
	}

	// Cari rekening terakhir untuk mencatat jurnal Delete
	var oldTx models.RekeningTransaction
	config.GetDB().Where("tabletransaction = 'qardhassan' AND idtabletransaction = ?", id).
		Order("created_at DESC").First(&oldTx)
	idNusvaRekening := oldTx.IDNusvaRekening

	if err := repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete qardhassan"})
		return
	}

	// Buat record DELETE di rekening_transaction
	if idNusvaRekening > 0 {
		deleteTx := &models.RekeningTransaction{
			TransactionType:    "Delete",
			IDNusvaRekening:    idNusvaRekening,
			TableTransaction:   "qardhassan",
			IDTableTransaction: uint(id),
			TanggalTransaksi:   time.Now(),
			NominalTransaksi:   qardhassan.NominalPinjaman, // Positif untuk mengembalikan saldo kas yang keluar
			CreatedBy:          operatorName,
			UpdatedBy:          operatorName,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		repo.CreateRekeningTransaction(deleteTx)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "QardHassan deleted successfully",
	})
}

// ===============================================
// FITUR PELUNASAN QARD HASSAN
// ===============================================

// Struct untuk request pembayaran
type PayQardHassanRequest struct {
	IDRekening        int     `json:"idrekening" binding:"required"`
	TanggalPembayaran string  `json:"tanggalpembayaran" binding:"required"`
	NominalPembayaran float64 `json:"nominalpembayaran" binding:"required,gt=0"`
}

// PayQardHassan mencatat pelunasan Qard Hassan
func PayQardHassan(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req PayQardHassanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tanggalPembayaran, err := time.Parse("2006-01-02", req.TanggalPembayaran)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal salah"})
		return
	}

	repo := repositories.NewQardHassanRepository()
	qardhassan, err := repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data Qard Hassan tidak ditemukan"})
		return
	}

	// Validasi jika sudah lunas
	if qardhassan.NominalPembayaran > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Qard Hassan ini sudah dibayar/lunas"})
		return
	}

	// 1. Update data Qard Hassan
	qardhassan.NominalPembayaran = req.NominalPembayaran
	qardhassan.TanggalPembayaran = &tanggalPembayaran
	qardhassan.UpdatedBy = operatorName
	qardhassan.UpdatedAt = time.Now()

	if err := repo.Update(qardhassan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pembayaran"})
		return
	}

	// 2. Insert ke Rekening Transaction (UANG MASUK -> NOMINAL POSITIF)
	rekeningTransaction := &models.RekeningTransaction{
		TransactionType:    "Insert",
		IDNusvaRekening:    uint(req.IDRekening),
		TableTransaction:   "qardhassan", // Tetap di-link ke qardhassan agar tracking mudah
		IDTableTransaction: uint(id),
		TanggalTransaksi:   tanggalPembayaran,
		NominalTransaksi:   req.NominalPembayaran, // Positif karena kas koperasi bertambah
		CreatedBy:          operatorName,
		UpdatedBy:          operatorName,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	repo.CreateRekeningTransaction(rekeningTransaction)

	c.JSON(http.StatusOK, gin.H{
		"message": "Pembayaran berhasil dicatat",
		"data":    qardhassan,
	})
}