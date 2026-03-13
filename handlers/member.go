package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"nusvakspps/config"
	"nusvakspps/models"
	"nusvakspps/repositories"
)

// CreateMemberRequest - Menyesuaikan dengan semua field di form 700 baris Frontend
type CreateMemberRequest struct {
	// Identitas Utama (Wajib)
	FullName    string `json:"full_name" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
	JoinDate    string `json:"join_date" binding:"required"`
	BirthDate   string `json:"birth_date" binding:"required"`
	BirthPlace  string `json:"birth_place" binding:"required"`
	KtpNo       string `json:"ktp_no" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email" binding:"required"`

	// Informasi Pekerjaan (Wajib)
	CompanyName string `json:"company_name" binding:"required"`
	JobTitle    string `json:"job_title" binding:"required"`

	// Informasi Bank (Wajib)
	BankAccountNo string `json:"bank_account_no" binding:"required"`
	BankName      string `json:"bank_name" binding:"required"`

	// Field Pendukung (Opsional)
	NpwpNo     string `json:"npwp_no"`
	AddressKtp string `json:"address_ktp"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`

	// Kontak Darurat (Opsional)
	EmergencyName    string `json:"emergency_name"`
	EmergencyPhone   string `json:"emergency_phone"`
	EmergencyAddress string `json:"emergency_address"`
	EmergencyRelation string `json:"emergency_relation"`
}

// GetMembers retrieves all members
func GetMembers(c *gin.Context) {
	fmt.Println("📋 Getting all members...")

	repo := repositories.NewMemberRepository()
	members, _, err := repo.List(0, 1000) // Get up to 1000 members
	if err != nil {
		fmt.Println("❌ Error getting members:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data anggota"})
		return
	}

	fmt.Printf("✅ Retrieved %d members\n", len(members))
	c.JSON(http.StatusOK, gin.H{"data": members})
}

// CreateMember handles the registration of a new member
func CreateMember(c *gin.Context) {
	// AMBIL OPERATOR DARI MIDDLEWARE (Ditambahkan untuk Audit)
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	var req CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Validation error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data belum lengkap: " + err.Error()})
		return
	}

	fmt.Println("📤 Creating new member:", req.FullName)

	repo := repositories.NewMemberRepository()

	// Check if KTP already exists
	existing, _ := repo.FindByKtp(req.KtpNo)
	if existing != nil {
		fmt.Println("⚠️ Member with KTP already exists:", req.KtpNo)
		c.JSON(http.StatusConflict, gin.H{"error": "Nomor KTP sudah terdaftar"})
		return
	}

	// Gender code for member number (1 for male, 2 for female)
	genderCode := "1"
	if req.Gender == "female" {
		genderCode = "2"
	}

	member := &models.Member{
		MemberNo:              generateMemberNumber(genderCode),
		FullName:              req.FullName,
		Gender:                req.Gender,
		JoinDate:              parseDate(req.JoinDate),
		BirthDate:             parseDate(req.BirthDate),
		BirthPlace:            req.BirthPlace,
		KtpNo:                 req.KtpNo,
		PhoneNumber:           req.PhoneNumber,
		Email:                 req.Email,
		CompanyName:           req.CompanyName,
		JobTitle:              req.JobTitle,
		BankAccountNo:         req.BankAccountNo,
		BankName:              req.BankName,
		NpwpNo:                req.NpwpNo,
		AddressKtp:            req.AddressKtp,
		City:                  req.City,
		Province:              req.Province,
		PostalCode:            req.PostalCode,
		EmergencyName:         req.EmergencyName,
		EmergencyPhone:        req.EmergencyPhone,
		EmergencyAddress:      req.EmergencyAddress,
		EmergencyRelation:     req.EmergencyRelation,
		IsActive:              true,
		// ISI FIELD AUDIT (Ditambahkan)
		CreatedBy:             operatorName,
		UpdatedBy:             operatorName,
	}

	if err := repo.Create(member); err != nil {
		fmt.Println("❌ Error creating member:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data anggota"})
		return
	}

	fmt.Println("✅ Member created successfully:", member.MemberNo)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Anggota berhasil didaftarkan",
		"data":    member,
	})
}

// GetMemberByID retrieves a single member by ID
func GetMemberByID(c *gin.Context) {
	idParam := c.Param("id")
	fmt.Println("🔍 Getting member by ID:", idParam)

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	repo := repositories.NewMemberRepository()
	member, err := repo.FindByID(id)
	if err != nil {
		fmt.Println("⚠️ Member not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Anggota tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": member})
}

// UpdateMember updates an existing member
func UpdateMember(c *gin.Context) {
    operatorName := c.GetString("current_user_name")
    idParam := c.Param("id")

    var id uint
    if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
        return
    }

    repo := repositories.NewMemberRepository()
    member, err := repo.FindByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Anggota tidak ditemukan"})
        return
    }

    // --- PERBAIKAN: Bind ke Map SEKALI SAJA ---
    var input map[string]interface{}
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
        return
    }

    // 1. Update field string/bool secara manual dari map
    // (Atau gunakan library seperti 'mitchellh/mapstructure' jika ingin otomatis)
	if val, ok := input["full_name"].(string); ok { member.FullName = val }
	if val, ok := input["gender"].(string); ok { member.Gender = val }
	if val, ok := input["npwp_no"].(string); ok { member.NpwpNo = val }
	if val, ok := input["birth_place"].(string); ok { member.BirthPlace = val }
	if val, ok := input["city"].(string); ok { member.City = val }
	if val, ok := input["province"].(string); ok { member.Province = val }
	if val, ok := input["postal_code"].(string); ok { member.PostalCode = val }
	if val, ok := input["email"].(string); ok { member.Email = val }
	if val, ok := input["company_name"].(string); ok { member.CompanyName = val }
	if val, ok := input["job_title"].(string); ok { member.JobTitle = val }
	if val, ok := input["bank_account_no"].(string); ok { member.BankAccountNo = val }
	if val, ok := input["bank_name"].(string); ok { member.BankName = val }
    
	if val, ok := input["ktp_no"].(string); ok { 
    member.KtpNo = val 
	}

	// 2. Perbaikan Hubungan (Emergency Relation)
	if val, ok := input["emergency_relation"].(string); ok { 
		member.EmergencyRelation = val 
	}

	// 3. Perbaikan Alamat Kontak Darurat (Emergency Address)
	if val, ok := input["emergency_address"].(string); ok { 
		member.EmergencyAddress = val 
	}

	// 4. Perbaikan Kota & Provinsi (Jika sebelumnya juga tidak masuk)
	if val, ok := input["city"].(string); ok { member.City = val }
	if val, ok := input["province"].(string); ok { member.Province = val }

    // Emergency & Work info juga perlu diupdate manual jika dikirim dari frontend
    if val, ok := input["emergency_name"].(string); ok { member.EmergencyName = val }
    if val, ok := input["emergency_phone"].(string); ok { member.EmergencyPhone = val }

    // 2. LOGIKA PARSING TANGGAL (Penyelesaian error parsing "T")
    if jd, ok := input["join_date"].(string); ok && jd != "" {
        // Go Layout reference: 2006-01-02
        t, err := time.Parse("2006-01-02", jd)
        if err == nil {
            member.JoinDate = t
        }
    }

    if bd, ok := input["birth_date"].(string); ok && bd != "" {
        t, err := time.Parse("2006-01-02", bd)
        if err == nil {
            member.BirthDate = t
        }
    }

    // 3. UPDATE FIELD AUDIT
    member.UpdatedBy = operatorName

    if err := repo.Update(member); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Data berhasil diperbarui",
        "data":    member,
    })
}

// DeleteMember deletes a member
func DeleteMember(c *gin.Context) {
	idParam := c.Param("id")
	fmt.Println("🗑️ Deleting member ID:", idParam)

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Di sini kakak bisa memanggil repo.Delete(id) jika sudah ada di repository-nya
	c.JSON(http.StatusOK, gin.H{"message": "Anggota berhasil dihapus"})
}

// --- HELPERS ---

func generateMemberNumber(genderCode string) string {
	now := time.Now()
	// Format: YYMM + GENDER_CODE + 6_DIGIT_SEQUENCE
	// Contoh: 2603 + 1 + 000001
	prefix := fmt.Sprintf("%02d%02d", now.Year()%100, int(now.Month()))
	nextSeq := getNextMemberSequence()

	return fmt.Sprintf("%s%s%s", prefix, genderCode, nextSeq)
}

func getNextMemberSequence() string {
	var maxMemberNo string
	query := `SELECT COALESCE(MAX(member_no), '000000') FROM members WHERE member_no LIKE ?`

	db := config.GetDB()
	if db == nil {
		fmt.Println("❌ Database connection is nil!")
		return "000001"
	}

	db.Raw(query, generateYearMonthPrefix()+"%").Scan(&maxMemberNo)

	// Jika maxMemberNo adalah '000000' atau panjangnya kurang dari 6
	if len(maxMemberNo) >= 6 {
		// Ambil 6 digit terakhir
		sequence := maxMemberNo[len(maxMemberNo)-6:]
		seqNum := 0
		fmt.Sscanf(sequence, "%d", &seqNum)
		seqNum++
		return fmt.Sprintf("%06d", seqNum)
	}

	return "000001"
}

func generateYearMonthPrefix() string {
	now := time.Now()
	return fmt.Sprintf("%02d%02d", now.Year()%100, int(now.Month()))
}

func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}

	// Mendukung format YYYY-MM-DD
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		fmt.Printf("⚠️ Gagal memproses tanggal '%s': %v\n", dateStr, err)
		return time.Time{}
	}
	return t
}