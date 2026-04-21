package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2" // Library untuk handle Excel

	"nusvakspps/config"
	"nusvakspps/models"
	"nusvakspps/repositories"
)

// CreateMemberRequest - Struktur data dari Frontend
type CreateMemberRequest struct {
	FullName          string  `json:"full_name" binding:"required"`
	Gender            string  `json:"gender" binding:"required"`
	JoinDate          string  `json:"join_date" binding:"required"`
	BirthDate         string  `json:"birth_date" binding:"required"`
	BirthPlace        string  `json:"birth_place" binding:"required"`
	KtpNo             string  `json:"ktp_no" binding:"required"`
	PhoneNumber       string  `json:"phone_number" binding:"required"`
	Email             string  `json:"email" binding:"required"`
	CompanyName       string  `json:"company_name" binding:"required"`
	JobTitle          string  `json:"job_title" binding:"required"`
	BankAccountNo     string  `json:"bank_account_no" binding:"required"`
	BankName          string  `json:"bank_name" binding:"required"`
	NpwpNo            string  `json:"npwp_no"`
	AddressKtp        string  `json:"address_ktp"`
	City              string  `json:"city"`
	Province          string  `json:"province"`
	PostalCode        string  `json:"postal_code"`
	EmergencyName     string  `json:"emergency_name"`
	EmergencyPhone    string  `json:"emergency_phone"`
	EmergencyAddress  string  `json:"emergency_address"`
	EmergencyRelation string  `json:"emergency_relation"`
	QardhassanPlafon  float64 `json:"qardhassanplafon"`
}

// GetMembers - Mengambil semua data anggota
func GetMembers(c *gin.Context) {
	repo := repositories.NewMemberRepository()
	members, _, err := repo.List(0, 1000)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data anggota"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": members})
}

// CreateMember - Input anggota baru secara manual
func CreateMember(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	var req CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data belum lengkap: " + err.Error()})
		return
	}

	repo := repositories.NewMemberRepository()

	existing, _ := repo.FindByKtp(req.KtpNo)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Nomor KTP sudah terdaftar"})
		return
	}

	genderCode := "1"
	if req.Gender == "female" {
		genderCode = "2"
	}

	member := &models.Member{
		MemberNo:          generateMemberNumber(genderCode),
		FullName:          req.FullName,
		Gender:            req.Gender,
		JoinDate:          parseDate(req.JoinDate),
		BirthDate:         parseDate(req.BirthDate),
		BirthPlace:        req.BirthPlace,
		KtpNo:             req.KtpNo,
		PhoneNumber:       req.PhoneNumber,
		Email:             req.Email,
		CompanyName:       req.CompanyName,
		JobTitle:          req.JobTitle,
		BankAccountNo:     req.BankAccountNo,
		BankName:          req.BankName,
		NpwpNo:            req.NpwpNo,
		AddressKtp:        req.AddressKtp,
		City:              req.City,
		Province:          req.Province,
		PostalCode:        req.PostalCode,
		EmergencyName:     req.EmergencyName,
		EmergencyPhone:    req.EmergencyPhone,
		EmergencyAddress:  req.EmergencyAddress,
		EmergencyRelation: req.EmergencyRelation,
		IsActive:          true,
		QardhassanPlafon:  req.QardhassanPlafon,
		CreatedBy:         operatorName,
		UpdatedBy:         operatorName,
	}

	if err := repo.Create(member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data anggota"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Anggota berhasil didaftarkan", "data": member})
}

// ImportMembers - Logic Import Excel
func ImportMembers(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File Excel wajib diunggah"})
		return
	}

	openedFile, _ := file.Open()
	defer openedFile.Close()

	f, err := excelize.OpenReader(openedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca file excel"})
		return
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data sheet"})
		return
	}

	db := config.GetDB()
	var successCount, updateCount int

	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}

		mNo := row[0] // ID Anggota (Kunci untuk Update)
		name := row[1]

		var member models.Member
		res := db.Where("member_no = ?", mNo).First(&member)

		// Mapping Field
		member.MemberNo = mNo
		member.FullName = name
		if len(row) > 2 { member.Gender = row[2] }
		if len(row) > 3 { member.JoinDate = parseDate(row[3]) }
		if len(row) > 4 { member.BirthDate = parseDate(row[4]) }
		if len(row) > 6 { member.KtpNo = row[6] }
		if len(row) > 8 { member.AddressKtp = row[8] }
		if len(row) > 12 { member.PhoneNumber = row[12] }
		
		member.IsActive = true
		member.UpdatedBy = operatorName

		if res.Error != nil {
			// Jika belum ada, buat nomor anggota baru jika di excel kosong
			if member.MemberNo == "" {
				member.MemberNo = generateMemberNumber("1")
			}
			member.CreatedBy = operatorName
			if err := db.Create(&member).Error; err == nil {
				successCount++
			}
		} else {
			if err := db.Save(&member).Error; err == nil {
				updateCount++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Proses Import Selesai",
		"inserted": successCount,
		"updated":  updateCount,
	})
}

// GetMemberByID - Ambil satu anggota
func GetMemberByID(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.ParseUint(idParam, 10, 32)

	repo := repositories.NewMemberRepository()
	member, err := repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anggota tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": member})
}

// UpdateMember - Edit data anggota
func UpdateMember(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	idParam := c.Param("id")
	id, _ := strconv.ParseUint(idParam, 10, 32)

	repo := repositories.NewMemberRepository()
	member, err := repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anggota tidak ditemukan"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Update manual fields
	if val, ok := input["full_name"].(string); ok { member.FullName = val }
	if val, ok := input["gender"].(string); ok { member.Gender = val }
	if val, ok := input["ktp_no"].(string); ok { member.KtpNo = val }
	if val, ok := input["address_ktp"].(string); ok { member.AddressKtp = val }
	
	member.UpdatedBy = operatorName

	if err := repo.Update(member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui", "data": member})
}

// DeleteMember - Hapus anggota
func DeleteMember(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Fitur hapus dinonaktifkan untuk integritas data"})
}

// --- HELPERS ---

func generateMemberNumber(genderCode string) string {
	now := time.Now()
	prefix := fmt.Sprintf("%02d%02d", now.Year()%100, int(now.Month()))
	nextSeq := getNextMemberSequence()
	return fmt.Sprintf("%s%s%s", prefix, genderCode, nextSeq)
}

func getNextMemberSequence() string {
	var maxMemberNo string
	db := config.GetDB()
	query := `SELECT COALESCE(MAX(member_no), '000000') FROM members WHERE member_no LIKE ?`
	db.Raw(query, generateYearMonthPrefix()+"%").Scan(&maxMemberNo)

	if len(maxMemberNo) >= 6 {
		sequence := maxMemberNo[len(maxMemberNo)-6:]
		seqNum, _ := strconv.Atoi(sequence)
		seqNum++
		return fmt.Sprintf("%06d", seqNum)
	}
	return "000001"
}

func generateYearMonthPrefix() string {
	return time.Now().Format("0601")
}

func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}