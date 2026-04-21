package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"nusvakspps/config"
	"nusvakspps/models"
	"nusvakspps/repositories"
)

func GetRoles(c *gin.Context) {
	roleRepo := repositories.NewRoleRepository()
	roles, err := roleRepo.List()
	if err != nil {
		fmt.Println("❌ Error getting roles:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve roles"})
		return
	}

	fmt.Println("✅ Retrieved", len(roles), "roles")

	c.JSON(http.StatusOK, gin.H{
		"data":  roles,
		"total": len(roles),
	})
}

func CreateRole(c *gin.Context) {
	var role struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("📤 Creating role:", role.Name)

	roleRepo := repositories.NewRoleRepository()

	// Check if role already exists
	if _, err := roleRepo.FindByName(role.Name); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Role already exists"})
		return
	}

	newRole := &models.Role{
		Name:        role.Name,
		Description: role.Description,
		IsActive:    true,
	}

	if err := roleRepo.Create(newRole).Error; err != nil {
		fmt.Println("❌ Error creating role:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	fmt.Println("✅ Role created successfully with ID:", newRole.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role created successfully",
		"data":    newRole,
	})
}

func GetMenus(c *gin.Context) {
    fmt.Println("🚀 API GetMenus Berhasil Dipanggil!") 
    
    var menus []models.Menu
    // Ambil semua data tanpa filter dulu untuk memastikan koneksi DB ok
    err := config.GetDB().Find(&menus).Error
    
    if err != nil {
        fmt.Println("❌ Gagal tarik data menu:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    fmt.Printf("✅ Berhasil menarik %d data menu\n", len(menus))
    
    // Kirim response dengan struktur yang diharapkan Frontend
    c.JSON(http.StatusOK, gin.H{
        "data": menus,
    })
}

// GetPermissions - Mengambil daftar semua hak akses
func GetPermissions(c *gin.Context) {
	var permissions []models.Permission
	// Urutkan berdasarkan menu_id agar rapi di tampilan
	if err := config.GetDB().Order("menu_id ASC, id ASC").Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data permission"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": permissions})
}

func GetRolePermissions(c *gin.Context) {
    roleID := c.Param("role-id")
    
    // LOG 1: Cek apakah ID role yang dikirim Frontend beneran masuk
    fmt.Printf("\n--- DEBUG GET PERMISSIONS ---\n")
    fmt.Printf("🔍 Request masuk untuk Role ID: [%s]\n", roleID)

    var rolePerms []models.RolePermission
    
    // Ambil data dari database
    result := config.GetDB().Where("role_id = ? AND is_allowed = ?", roleID, true).Find(&rolePerms)
    
    if result.Error != nil {
        fmt.Printf("❌ Database Error: %v\n", result.Error)
        c.JSON(500, gin.H{"data": []uint{}})
        return
    }

    // LOG 2: Cek berapa banyak record yang ditemukan di DB
    fmt.Printf("📊 Record ditemukan di database: %d baris\n", len(rolePerms))

    var permIDs []uint
    for _, rp := range rolePerms {
        permIDs = append(permIDs, rp.PermissionID)
    }

    // LOG 3: Lihat isi array ID yang bakal dikirim ke Frontend
    fmt.Printf("📦 ID yang dikirim ke Frontend: %v\n", permIDs)
    fmt.Printf("-----------------------------\n")

    c.JSON(http.StatusOK, gin.H{
        "data": permIDs,
    })
}

// Struct untuk menerima data array checkbox dari frontend
type AssignPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids"`
}

// AssignPermissions - Menyimpan centangan hak akses dari frontend
func AssignPermissions(c *gin.Context) {
	roleIDParam := c.Param("role-id")
	roleID, err := strconv.Atoi(roleIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Role tidak valid"})
		return
	}

	var req AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	db := config.GetDB()
	tx := db.Begin() // Mulai transaksi database agar aman

	// 1. Bersihkan (Hapus) semua permission lama untuk role ini
	if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membersihkan hak akses lama"})
		return
	}

	// 2. Insert permission baru sesuai yang dicentang di frontend
	for _, permID := range req.PermissionIDs {
		rp := models.RolePermission{
			RoleID:       uint(roleID),
			PermissionID: permID,
			IsAllowed:    true,
		}
		if err := tx.Create(&rp).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan hak akses baru"})
			return
		}
	}

	// Jika semua lancar, commit transaksinya
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Hak akses berhasil diperbarui"})
}

func GetAuditLogs(c *gin.Context) {
	page := 1
	limit := 20

	// Tangkap parameter query untuk paginasi (jika ada dari frontend)
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	offset := (page - 1) * limit
	repo := repositories.NewAuditRepository()
	
	// Tarik data dari repository
	logs, total, err := repo.List(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data audit log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
