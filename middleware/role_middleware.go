package middleware // <--- BARIS INI WAJIB ADA DI PALING ATAS

import (
	"github.com/gin-gonic/gin" // <--- WAJIB ADA INI
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Ambil roles user yang disimpan saat login (biasanya di AuthMiddleware)
        // Pastikan AuthMiddleware kamu menyimpan roles ke context
        userRolesInterface, exists := c.Get("user_roles") 
        if !exists {
            c.JSON(403, gin.H{"error": "Hak akses tidak ditemukan"})
            c.Abort()
            return
        }

        userRoles := userRolesInterface.([]string)
        
        isAllowed := false
        for _, r := range userRoles {
            for _, allowed := range allowedRoles {
                if r == allowed {
                    isAllowed = true
                    break
                }
            }
        }

        if !isAllowed {
            c.JSON(403, gin.H{"error": "Role anda tidak diizinkan mengakses fitur ini"})
            c.Abort()
            return
        }
        c.Next()
    }
}