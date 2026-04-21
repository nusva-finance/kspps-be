package middleware

import (
	"net/http"
	"nusvakspps/config"
	"github.com/gin-gonic/gin"
)

func RequirePermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil ID User dari context
		userID, exists := c.Get("current_user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
			c.Abort()
			return
		}

		// 2. Eksekusi Query
		var count int64
		// Kita pisahkan DB ke variabel agar tidak memicu error 'cannot slice'
		db := config.GetDB()
		err := db.Table("users u").
			Joins("JOIN user_roles ur ON u.id = ur.user_id").
			Joins("JOIN role_permissions rp ON ur.role_id = rp.role_id").
			Joins("JOIN permissions p ON rp.permission_id = p.id").
			Where("u.id = ? AND rp.is_allowed = ? AND (p.code = ? OR p.code = '*')", userID, true, requiredPermission).
			Count(&count).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan verifikasi akses"})
			c.Abort()
			return
		}

		// 3. Jika tidak punya akses (dan bukan super_admin)
		if count == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Akses Ditolak: Anda tidak memiliki hak akses (" + requiredPermission + ")",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}