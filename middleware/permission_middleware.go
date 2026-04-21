package middleware

import (
	"fmt"
	"nusvakspps/config"
	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil roles dari context (setelah AuthMiddleware)
		val, exists := c.Get("user_roles")
		if !exists || val == nil {
			c.JSON(403, gin.H{"error": "Forbidden", "message": "Sesi tidak memiliki role"})
			c.Abort()
			return
		}

		userRoles, ok := val.([]string)
		if !ok {
			c.JSON(500, gin.H{"error": "Internal Server Error", "message": "Gagal membaca format role"})
			c.Abort()
			return
		}

		// 2. LANGSUNG CEK KE DATABASE (Tanpa Bypass Admin)
		db := config.GetDB()
		var count int64
		
		// Query ini mengecek apakah salah satu role user punya permission_code terkait 
		// DAN is_allowed-nya harus TRUE
        err := db.Table("role_permissions rp").
            Joins("JOIN roles r ON r.id = rp.role_id").
            Joins("JOIN permissions p ON p.id = rp.permission_id").
            Where("r.name IN ?", userRoles).
            Where("rp.is_allowed = ?", true).
            // KUNCINYA DI SINI: Pakai OR untuk cek apakah dia punya kode spesifik ATAU punya '*'
            Where("(p.code = ? OR p.code = '*')", requiredPermission). 
            Count(&count).Error

		// Debug Log untuk memantau traffic akses
		fmt.Printf("🛡️ RBAC Audit: Roles=%v | Permission=%s | Result=%d\n", userRoles, requiredPermission, count)

		if err != nil {
			c.JSON(500, gin.H{"error": "Database Error", "message": err.Error()})
			c.Abort()
			return
		}

		// 3. JIKA TIDAK ADA IZIN (Count 0), TENDANG!
		if count == 0 {
			c.JSON(403, gin.H{
				"error": "Forbidden",
				"message": fmt.Sprintf("Akses ditolak! Anda tidak punya izin [%s]", requiredPermission),
			})
			c.Abort()
			return
		}

		// 4. JIKA LOLOS, LANJUTKAN
		c.Next()
	}
}