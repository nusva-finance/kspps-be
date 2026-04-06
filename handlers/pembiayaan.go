package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"

	"nusvakspps/models"
	"nusvakspps/repositories"
)

// RoundUp - Membulatkan nilai ke atas ke kelipatan tertentu
// Contoh: RoundUp(586700, 10000) = 590000
func RoundUp(val float64, roundTo float64) float64 {
	return math.Ceil(val/roundTo) * roundTo
}

// HitungJatuhTempo - Menghitung tanggal jatuh tempo berdasarkan tanggal pinjaman
// Jika tanggal 1-15: +1 bulan
// Jika tanggal 16-31: +2 bulan
func HitungJatuhTempo(tglPinjaman time.Time) time.Time {
	hari := tglPinjaman.Day()
	if hari <= 15 {
		return tglPinjaman.AddDate(0, 1, 0)
	}
	return tglPinjaman.AddDate(0, 2, 0)
}

type CreatePembiayaanRequest struct {
	IDMember         int     `json:"idmember" binding:"required"`
	TipePinjaman    string  `json:"tipepinjaman" binding:"required"`
	TanggalPinjaman string  `json:"tanggalpinjaman" binding:"required"`
	KategoriBarang  string  `json:"kategoribarang" binding:"required"`
	Tenor           int     `json:"tenor" binding:"required"`
	NominalPembelian float64 `json:"nominalpembelian" binding:"required,gt=0"`
	IDNusvaRekening  uint    `json:"idnusvarekening" binding:"required"`
}

type UpdatePembiayaanRequest struct {
	IDMember         int     `json:"idmember"`
	TipePinjaman    string  `json:"tipepinjaman"`
	TanggalPinjaman string  `json:"tanggalpinjaman"`
	KategoriBarang  string  `json:"kategoribarang"`
	Tenor           int     `json:"tenor"`
	NominalPembelian float64 `json:"nominalpembelian"`
}

// GetPembiayaan - Mengambil daftar pembiayaan
func GetPembiayaan(c *gin.Context) {
	fmt.Println("📋 Getting all pembiayaan...")

	memberID := c.Query("member_id")
	repo := repositories.NewPembiayaanRepository()

	var pembiayaans []models.PembiayaanWithMemberName
	var total int64
	var err error

	if memberID != "" {
		var id int
		if _, err := fmt.Sscanf(memberID, "%d", &id); err == nil {
			pembiayaans, total, err = repo.ListByMemberID(id, 0, 1000)
		} else {
			pembiayaans, total, err = repo.List(0, 1000)
		}
	} else {
		pembiayaans, total, err = repo.List(0, 1000)
	}

	if err != nil {
		fmt.Println("❌ Error getting pembiayaan:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pembiayaan"})
		return
	}

	fmt.Printf("✅ Retrieved %d pembiayaan\n", len(pembiayaans))
	c.JSON(http.StatusOK, gin.H{"data": pembiayaans, "total": total})
}

// GetPembiayaanByID - Mengambil pembiayaan berdasarkan ID
func GetPembiayaanByID(c *gin.Context) {
	idParam := c.Param("id")
	fmt.Println("🔍 Getting pembiayaan by ID:", idParam)

	var id int
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	repo := repositories.NewPembiayaanRepository()
	pembiayaan, err := repo.FindByID(id)
	if err != nil {
		fmt.Println("⚠️ Pembiayaan not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Pembiayaan tidak ditemukan"})
		return
	}

	// Get idnusvarekening from rekening_transaction (original Insert record)
	rekeningTransactionRepo := repositories.NewRekeningTransactionRepository()
	idNusvaRekening, err := rekeningTransactionRepo.GetIDNusvaRekeningByPembiayaanID(uint(pembiayaan.IDPinjaman))
	if err != nil {
		fmt.Println("⚠️ Could not get idnusvarekening:", err)
		idNusvaRekening = 0
	}
	fmt.Println("✅ ID Nusva Rekening from transaction:", idNusvaRekening)

	// Return pembiayaan with idnusvarekening
	response := map[string]interface{}{
		"idpinjaman":        pembiayaan.IDPinjaman,
		"idmember":          pembiayaan.IDMember,
		"tipepinjaman":      pembiayaan.TipePinjaman,
		"tanggalpinjaman":   pembiayaan.TanggalPinjaman,
		"kategoribarang":    pembiayaan.KategoriBarang,
		"tenor":             pembiayaan.Tenor,
		"margin":            pembiayaan.Margin,
		"nominalpinjaman":   pembiayaan.NominalPinjaman,
		"nominalpembelian":  pembiayaan.NominalPembelian,
		"tgljtangsuran1":    pembiayaan.TglJtAngsuran1,
		"idnusvarekening":   idNusvaRekening,
		"created_by":        pembiayaan.CreatedBy,
		"updated_by":        pembiayaan.UpdatedBy,
		"created_at":        pembiayaan.CreatedAt,
		"updated_at":        pembiayaan.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// CreatePembiayaan - Membuat pembiayaan baru
func CreatePembiayaan(c *gin.Context) {
	// AMBIL OPERATOR DARI MIDDLEWARE
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	var req CreatePembiayaanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Validation error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data belum lengkap: " + err.Error()})
		return
	}

	fmt.Println("📤 Creating new pembiayaan:")
	fmt.Println("   ID Member:", req.IDMember)
	fmt.Println("   Tipe Pinjaman:", req.TipePinjaman)
	fmt.Println("   Tanggal Pinjaman:", req.TanggalPinjaman)
	fmt.Println("   Kategori Barang:", req.KategoriBarang)
	fmt.Println("   Tenor:", req.Tenor)
	fmt.Println("   Nominal Pembelian:", req.NominalPembelian)
	fmt.Println("   ID Nusva Rekening:", req.IDNusvaRekening)
	fmt.Println("   Operator:", operatorName)

	// Parse tanggal pinjaman
	tanggalPinjaman, err := time.Parse("2006-01-02", req.TanggalPinjaman)
	if err != nil {
		fmt.Println("❌ Error parsing tanggal pinjaman:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal pinjaman tidak valid (gunakan YYYY-MM-DD)"})
		return
	}

	// Lookup margin dari margin_setups berdasarkan kategori dan tenor
	marginRepo := repositories.NewMarginRepository()
	marginSetup, err := marginRepo.FindByCategoryAndTenor(req.KategoriBarang, req.Tenor)
	if err != nil {
		fmt.Println("❌ Error getting margin setup:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data margin"})
		return
	}

	if marginSetup == nil {
		fmt.Println("⚠️ Margin setup not found for category:", req.KategoriBarang, "tenor:", req.Tenor)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Margin tidak ditemukan untuk kategori '%s' dengan tenor %d bulan. Silakan setup margin terlebih dahulu.", req.KategoriBarang, req.Tenor)})
		return
	}

	margin := marginSetup.Margin

	// Kalkulasi nominal pinjaman (Nominal Pembiayaan) dari nominal pembelian
	// Nominal Pembiayaan = Nominal Pembelian + (Nominal Pembelian * Margin / 100)
	marginValue := margin / 100
	nominalPinjaman := req.NominalPembelian + (req.NominalPembelian * marginValue)
	nominalPinjamanDibulatkan := RoundUp(nominalPinjaman, 10000)

	// Kalkulasi tanggal jatuh tempo (jika tanggal 1-15: +1 bulan, jika 16-31: +2 bulan)
	tglJtAngsuran1 := HitungJatuhTempo(tanggalPinjaman)

	fmt.Println("📊 Calculation:")
	fmt.Println("   Margin:", margin, "(", (margin * 100), "% )")
	fmt.Println("   Margin Value (dibagi 100):", marginValue)
	fmt.Println("   Nominal Pembelian (input):", req.NominalPembelian)
	fmt.Println("   Nominal Pembiayaan (sebelum pembulatan):", nominalPinjaman)
	fmt.Println("   Nominal Pembiayaan (dibulatkan ke 10.000):", nominalPinjamanDibulatkan)
	fmt.Println("   Tgl Pinjaman:", tanggalPinjaman.Format("2006-01-02"))
	fmt.Println("   Tgl Jatuh Tempo:", tglJtAngsuran1.Format("2006-01-02"))

	pembiayaan := &models.Pembiayaan{
		IDMember:         req.IDMember,
		TipePinjaman:    req.TipePinjaman,
		TanggalPinjaman:  tanggalPinjaman,
		KategoriBarang:   req.KategoriBarang,
		Tenor:           req.Tenor,
		Margin:          margin,
		NominalPinjaman:  nominalPinjamanDibulatkan,
		NominalPembelian: req.NominalPembelian,
		TglJtAngsuran1:  tglJtAngsuran1,
		CreatedBy:        operatorName,
		UpdatedBy:        operatorName,
	}

	fmt.Println("📥 About to create pembiayaan object:", pembiayaan)

	pembiayaanRepo := repositories.NewPembiayaanRepository()
	if err := pembiayaanRepo.Create(pembiayaan); err != nil {
		fmt.Println("❌ Error creating pembiayaan:")
		fmt.Println("   Error:", err.Error())
		fmt.Println("   Error Type:", fmt.Sprintf("%T", err))

		// Check if it's a duplicate constraint error
		errMsg := err.Error()
		if len(errMsg) > 0 && (strings.Contains(strings.ToLower(errMsg), "duplicate") || strings.Contains(strings.ToLower(errMsg), "unique") || strings.Contains(strings.ToLower(errMsg), "constraint")) {
			fmt.Println("⚠️ Duplicate constraint detected")
			c.JSON(http.StatusConflict, gin.H{"error": "Data pembiayaan sudah ada"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data pembiayaan: " + err.Error()})
		return
	}

	fmt.Println("✅ Pembiayaan created successfully:")
	fmt.Println("   ID Pinjaman:", pembiayaan.IDPinjaman)
	fmt.Println("   ID Member:", pembiayaan.IDMember)
	fmt.Println("   Nominal Pembelian:", pembiayaan.NominalPembelian)

	// Create rekening_transaction record with nominal * -1
	fmt.Println("📝 About to create rekening_transaction...")
	fmt.Println("   ID Nusva Rekening:", req.IDNusvaRekening)
	fmt.Println("   Table Transaction: pembiayaan")
	fmt.Println("   ID Table Transaction:", pembiayaan.IDPinjaman)
	fmt.Println("   Nominal Transaksi:", -1 * req.NominalPembelian)
	fmt.Println("   Tanggal Transaksi:", time.Now())

	rekeningTransactionRepo := repositories.NewRekeningTransactionRepository()
	rekeningTransaction := &models.RekeningTransaction{
		TransactionType:    "Insert",
		IDNusvaRekening:    req.IDNusvaRekening,
		TableTransaction:   "pembiayaan",
		IDTableTransaction: uint(pembiayaan.IDPinjaman),
		TanggalTransaksi:   time.Now(),
		NominalTransaksi:   -1 * req.NominalPembelian,
		CreatedBy:          operatorName,
		UpdatedBy:          operatorName,
	}

	fmt.Println("📝 RekeningTransaction object to create:", rekeningTransaction)

	if err := rekeningTransactionRepo.Create(rekeningTransaction); err != nil {
		fmt.Println("❌ Error creating rekening_transaction:")
		fmt.Println("   Error:", err.Error())
		fmt.Println("   Error Type:", fmt.Sprintf("%T", err))
		// Continue anyway, the pembiayaan was created successfully
	} else {
		fmt.Println("✅ Rekening Transaction created successfully:")
		fmt.Println("   ID Rekening Transaction:", rekeningTransaction.IDRekeningTransaction)
		fmt.Println("   Transaction Type:", rekeningTransaction.TransactionType)
		fmt.Println("   Table Transaction:", rekeningTransaction.TableTransaction)
		fmt.Println("   ID Table Transaction:", rekeningTransaction.IDTableTransaction)
		fmt.Println("   ID Nusva Rekening:", rekeningTransaction.IDNusvaRekening)
		fmt.Println("   Nominal Transaksi:", rekeningTransaction.NominalTransaksi)
		fmt.Println("   Tanggal Transaksi:", rekeningTransaction.TanggalTransaksi)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pembiayaan berhasil ditambahkan",
		"data":    pembiayaan,
	})
}

// UpdatePembiayaan - Mengubah data pembiayaan
func UpdatePembiayaan(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	idParam := c.Param("id")

	var id int
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	repo := repositories.NewPembiayaanRepository()
	pembiayaan, err := repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pembiayaan tidak ditemukan"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Store old nominal pembelian for rekening_transaction
	oldNominalPembelian := pembiayaan.NominalPembelian

	// Get old idnusvarekening from rekening_transaction (original Insert record)
	rekeningTransactionRepo := repositories.NewRekeningTransactionRepository()
	oldIDNusvaRekening, err := rekeningTransactionRepo.GetIDNusvaRekeningByPembiayaanID(uint(pembiayaan.IDPinjaman))
	if err != nil {
		fmt.Println("⚠️ Could not get old idnusvarekening:", err)
		oldIDNusvaRekening = 0
	}

	// Get new idnusvarekening from input (required for Update)
	var newIDNusvaRekening uint
	if val, ok := input["idnusvarekening"].(float64); ok {
		newIDNusvaRekening = uint(val)
	} else {
		// If not provided, use the old one
		newIDNusvaRekening = oldIDNusvaRekening
	}

	// Update fields from input
	if val, ok := input["idmember"].(float64); ok {
		pembiayaan.IDMember = int(val)
	}
	if val, ok := input["tipepinjaman"].(string); ok {
		pembiayaan.TipePinjaman = val
	}
	if val, ok := input["tanggalpinjaman"].(string); ok && val != "" {
		tanggalPinjaman, err := time.Parse("2006-01-02", val)
		if err == nil {
			pembiayaan.TanggalPinjaman = tanggalPinjaman
		}
	}
	if val, ok := input["kategoribarang"].(string); ok {
		pembiayaan.KategoriBarang = val
	}
	if val, ok := input["tenor"].(float64); ok {
		pembiayaan.Tenor = int(val)
	}
	if val, ok := input["nominalpembelian"].(float64); ok {
		pembiayaan.NominalPembelian = val
	}

	// Recalculate margin, nominal pinjaman, and tgl jatuh tempo if category, tenor, or nominal pembelian changed
	nominalPembelianChanged := false
	if _, ok := input["nominalpembelian"]; ok {
		nominalPembelianChanged = true
	}
	categoryChanged := false
	if _, ok := input["kategoribarang"]; ok {
		categoryChanged = true
	}
	tenorChanged := false
	if _, ok := input["tenor"]; ok {
		tenorChanged = true
	}

	if nominalPembelianChanged || categoryChanged || tenorChanged {
		// Get margin from category and tenor
		marginRepo := repositories.NewMarginRepository()
		marginSetup, err := marginRepo.FindByCategoryAndTenor(pembiayaan.KategoriBarang, pembiayaan.Tenor)
		if err != nil {
			fmt.Println("❌ Error getting margin setup:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data margin"})
			return
		}

		if marginSetup == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Margin tidak ditemukan untuk kategori '%s' dengan tenor %d bulan", pembiayaan.KategoriBarang, pembiayaan.Tenor)})
			return
		}

		pembiayaan.Margin = marginSetup.Margin
		// Kalkulasi nominal pinjaman dari nominal pembelian dengan pembulatan ke 10.000
		marginValue := pembiayaan.Margin / 100
		nominalPinjaman := pembiayaan.NominalPembelian + (pembiayaan.NominalPembelian * marginValue)
		pembiayaan.NominalPinjaman = RoundUp(nominalPinjaman, 10000)
		// Kalkulasi tanggal jatuh tempo
		pembiayaan.TglJtAngsuran1 = HitungJatuhTempo(pembiayaan.TanggalPinjaman)
	}

	// Update audit fields
	pembiayaan.UpdatedBy = operatorName

	if err := repo.Update(pembiayaan); err != nil {
		fmt.Println("❌ Error updating pembiayaan:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	fmt.Println("✅ Pembiayaan updated successfully:")
	fmt.Println("   ID Pinjaman:", pembiayaan.IDPinjaman)

	// Create 2 rekening_transaction records: Reverse and Update
	// Record 1: Reverse (membalik transaksi lama) - positive to reverse negative
	fmt.Println("📝 Creating Reverse rekening_transaction...")
	reverseTransaction := &models.RekeningTransaction{
		TransactionType:    "Reverse",
		IDNusvaRekening:    oldIDNusvaRekening,
		TableTransaction:   "pembiayaan",
		IDTableTransaction: uint(pembiayaan.IDPinjaman),
		TanggalTransaksi:   time.Now(),
		NominalTransaksi:   oldNominalPembelian, // Positive to reverse the negative from Insert
		CreatedBy:          operatorName,
		UpdatedBy:          operatorName,
	}
	if err := rekeningTransactionRepo.Create(reverseTransaction); err != nil {
		fmt.Println("❌ Error creating reverse transaction:", err)
	} else {
		fmt.Println("✅ Reverse transaction created successfully")
	}

	// Record 2: Update (mencatat transaksi baru) - negative for new value
	fmt.Println("📝 Creating Update rekening_transaction...")
	updateTransaction := &models.RekeningTransaction{
		TransactionType:    "Update",
		IDNusvaRekening:    newIDNusvaRekening,
		TableTransaction:   "pembiayaan",
		IDTableTransaction: uint(pembiayaan.IDPinjaman),
		TanggalTransaksi:   time.Now(),
		NominalTransaksi:   -1 * pembiayaan.NominalPembelian, // Negative for new value
		CreatedBy:          operatorName,
		UpdatedBy:          operatorName,
	}
	if err := rekeningTransactionRepo.Create(updateTransaction); err != nil {
		fmt.Println("❌ Error creating update transaction:", err)
	} else {
		fmt.Println("✅ Update transaction created successfully")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diperbarui",
		"data":    pembiayaan,
	})
}

// handlers/pembiayaan.go

func DeletePembiayaan(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	idParam := c.Param("id")
	fmt.Println("🗑️ Memproses penghapusan pembiayaan ID:", idParam)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println("❌ ID tidak valid:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pembiayaan tidak valid"})
		return
	}

	repo := repositories.NewPembiayaanRepository()

	// Ambil data untuk memastikan record ada sebelum dihapus
	pembiayaan, err := repo.FindByID(id)
	if err != nil {
		fmt.Println("⚠️ Data pembiayaan tidak ditemukan untuk ID:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	// GUNAKAN VARIABEL: Kita cetak info pembiayaan ke terminal agar variabel 'pembiayaan' terpakai
	fmt.Printf("✅ Menghapus data: [%s] untuk Member ID: %d dengan Nominal: %.2f\n",
		pembiayaan.TipePinjaman, pembiayaan.IDMember, pembiayaan.NominalPinjaman)

	// Get idnusvarekening from rekening_transaction (original Insert record)
	rekeningTransactionRepo := repositories.NewRekeningTransactionRepository()
	idNusvaRekening, err := rekeningTransactionRepo.GetIDNusvaRekeningByPembiayaanID(uint(pembiayaan.IDPinjaman))
	if err != nil {
		fmt.Println("⚠️ Could not get idnusvarekening:", err)
		idNusvaRekening = 0
	}

	// Eksekusi penghapusan di database
	if err := repo.Delete(id); err != nil {
		fmt.Println("❌ Gagal menghapus dari database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data: " + err.Error()})
		return
	}

	// Create rekening_transaction record for Delete
	fmt.Println("📝 Creating Delete rekening_transaction...")
	deleteTransaction := &models.RekeningTransaction{
		TransactionType:    "Delete",
		IDNusvaRekening:    idNusvaRekening,
		TableTransaction:   "pembiayaan",
		IDTableTransaction: uint(pembiayaan.IDPinjaman),
		TanggalTransaksi:   time.Now(),
		NominalTransaksi:   pembiayaan.NominalPembelian, // Positive to reverse the negative from Insert
		CreatedBy:          operatorName,
		UpdatedBy:          operatorName,
	}
	if err := rekeningTransactionRepo.Create(deleteTransaction); err != nil {
		fmt.Println("❌ Error creating delete transaction:", err)
	} else {
		fmt.Println("✅ Delete transaction created successfully")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembiayaan berhasil dihapus"})
}


// GetMarginByCategoryAndTenor - Helper endpoint untuk mendapatkan margin berdasarkan kategori dan tenor
func GetMarginByCategoryAndTenor(c *gin.Context) {
	category := c.Query("category")
	tenorParam := c.Query("tenor")

	if category == "" || tenorParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter category dan tenor wajib diisi"})
		return
	}

	var tenor int
	if _, err := fmt.Sscanf(tenorParam, "%d", &tenor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter tenor tidak valid"})
		return
	}

	fmt.Println("🔍 Getting margin for category:", category, "tenor:", tenor)

	repo := repositories.NewMarginRepository()
	marginSetup, err := repo.FindByCategoryAndTenor(category, tenor)
	if err != nil {
		fmt.Println("❌ Error getting margin setup:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data margin"})
		return
	}

	if marginSetup == nil {
		fmt.Println("⚠️ Margin setup not found")
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Margin tidak ditemukan untuk kategori '%s' dengan tenor %d bulan", category, tenor)})
		return
	}

	fmt.Println("✅ Margin found:", marginSetup.Margin)

	c.JSON(http.StatusOK, gin.H{"data": marginSetup})
}
