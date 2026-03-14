package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

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
	IDMember        int     `json:"idmember" binding:"required"`
	TipePinjaman    string  `json:"tipepinjaman" binding:"required"`
	TanggalPinjaman string  `json:"tanggalpinjaman" binding:"required"`
	KategoriBarang  string  `json:"kategoribarang" binding:"required"`
	Tenor           int     `json:"tenor" binding:"required"`
	NominalPinjaman float64 `json:"nominalpinjaman" binding:"required,gt=0"`
}

type UpdatePembiayaanRequest struct {
	IDMember        int     `json:"idmember"`
	TipePinjaman    string  `json:"tipepinjaman"`
	TanggalPinjaman string  `json:"tanggalpinjaman"`
	KategoriBarang  string  `json:"kategoribarang"`
	Tenor           int     `json:"tenor"`
	NominalPinjaman float64 `json:"nominalpinjaman"`
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

	c.JSON(http.StatusOK, gin.H{"data": pembiayaan})
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
	fmt.Println("   Nominal Pinjaman:", req.NominalPinjaman)
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

	// Kalkulasi nominal pembelian (dengan pembulatan ke 10.000)
	// Margin adalah persentase, jadi perlu dibagi 100 dulu
	marginValue := margin / 100
	totalPembelian := req.NominalPinjaman + (req.NominalPinjaman * marginValue)
	nominalPembelian := RoundUp(totalPembelian, 10000)

	// Kalkulasi tanggal jatuh tempo (jika tanggal 1-15: +1 bulan, jika 16-31: +2 bulan)
	tglJtAngsuran1 := HitungJatuhTempo(tanggalPinjaman)

	fmt.Println("📊 Calculation:")
	fmt.Println("   Margin:", margin, "(", (margin * 100), "% )")
	fmt.Println("   Margin Value (dibagi 100):", marginValue)
	fmt.Println("   Total Pembelian (sebelum pembulatan):", totalPembelian)
	fmt.Println("   Nominal Pembelian (dibulatkan ke 10.000):", nominalPembelian)
	fmt.Println("   Tgl Pinjaman:", tanggalPinjaman.Format("2006-01-02"))
	fmt.Println("   Tgl Jatuh Tempo:", tglJtAngsuran1.Format("2006-01-02"))

	pembiayaan := &models.Pembiayaan{
		IDMember:         req.IDMember,
		TipePinjaman:    req.TipePinjaman,
		TanggalPinjaman:  tanggalPinjaman,
		KategoriBarang:   req.KategoriBarang,
		Tenor:           req.Tenor,
		Margin:          margin,
		NominalPinjaman:  req.NominalPinjaman,
		NominalPembelian: nominalPembelian,
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
	if val, ok := input["nominalpinjaman"].(float64); ok {
		pembiayaan.NominalPinjaman = val
	}

	// Recalculate margin, nominal pembelian, and tgl jatuh tempo if category or tenor changed
	categoryChanged := false
	if _, ok := input["kategoribarang"]; ok {
		categoryChanged = true
	}
	tenorChanged := false
	if _, ok := input["tenor"]; ok {
		tenorChanged = true
	}

	if categoryChanged || tenorChanged {
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
		// Kalkulasi nominal pembelian dengan pembulatan ke 10.000
		// Margin adalah persentase, jadi perlu dibagi 100 dulu
		marginValue := pembiayaan.Margin / 100
		totalPembelian := pembiayaan.NominalPinjaman + (pembiayaan.NominalPinjaman * marginValue)
		pembiayaan.NominalPembelian = RoundUp(totalPembelian, 10000)
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diperbarui",
		"data":    pembiayaan,
	})
}

// DeletePembiayaan - Menghapus pembiayaan
func DeletePembiayaan(c *gin.Context) {
	idParam := c.Param("id")
	fmt.Println("🗑️ Deleting pembiayaan ID:", idParam)

	var id int
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		fmt.Println("❌ ID tidak valid:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pembiayaan tidak valid"})
		return
	}

	repo := repositories.NewPembiayaanRepository()

	// Cek apakah pembiayaan ini punya angsuran terkait
	// Cek tabel pembayaran (jika ada) untuk referensi
	pembiayaan, err := repo.FindByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			fmt.Println("⚠️ Pembiayaan tidak ditemukan:", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Data pembiayaan tidak ditemukan"})
		} else {
			fmt.Println("❌ Error saat mencari pembiayaan:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pembiayaan"})
		}
		return
	}

	fmt.Println("📊 Checking related records before delete...")
	fmt.Println("   ID Pinjaman:", pembiayaan.IDPinjaman)
	fmt.Println("   ID Member:", pembiayaan.IDMember)

	// Periksa apakah ada referensi dari tabel lain
	// Catatan: Asumsikan tabel lain menggunakan nama yang sama
	// Di Go kita bisa cek constraint error dari database

	if err := repo.Delete(id); err != nil {
		fmt.Println("❌ Error deleting pembiayaan:", err)
		// Cek apakah error adalah constraint violation
		errMsg := err.Error()
		if strings.Contains(strings.ToLower(errMsg), "duplicate") ||
		   strings.Contains(strings.ToLower(errMsg), "unique") ||
		   strings.Contains(strings.ToLower(errMsg), "constraint") ||
		   strings.Contains(strings.ToLower(errMsg), "foreign key") ||
		   strings.Contains(strings.ToLower(errMsg), "violates") {
			fmt.Println("⚠️ Constraint violation detected - member/pembayaran terkait")
			c.JSON(http.StatusConflict, gin.H{
				"error": "Tidak dapat menghapus data ini. Data ini memiliki referensi ke tabel lain (anggota, pembayaran, dll).",
			})
			return
		}

		fmt.Println("   Error type:", fmt.Sprintf("%T", err))
		fmt.Println("   Error message:", errMsg)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data pembiayaan: " + errMsg})
		return
	}

	fmt.Println("✅ Pembiayaan deleted successfully, ID:", id)
	fmt.Println("   ID Pinjaman:", pembiayaan.IDPinjaman)

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
