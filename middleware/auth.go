package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Gunakan secret key yang sama dengan yang ada di handlers/auth.go
var jwtSecret = []byte("nusvakspps-secret-key-2024-jwt-token")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil token dari header "Authorization"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Otorisasi diperlukan (Header kosong)"})
			c.Abort()
			return
		}

		// 2. Format header harus "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format otorisasi tidak valid"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Parse dan Validasi Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing-nya HMAC (sesuai dengan yang kita buat di Login)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak terduga: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		// 4. Jika token tidak valid atau expired
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi berakhir atau token tidak valid. Silakan login kembali."})
			c.Abort()
			return
		}

		// 5. Ekstrak data dari Claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// --- Simpan user_id ---
			if id, ok := claims["user_id"].(float64); ok {
				c.Set("current_user_id", uint(id))
			}

			// --- Simpan username ---
			if username, ok := claims["username"].(string); ok {
				c.Set("current_user_name", username)
			}

			// --- SETUP ROLES (TAMBAHKAN BAGIAN INI) ---
			if rolesInterface, ok := claims["roles"]; ok {
				// JWT biasanya menyimpan array sebagai interface slice
				var roles []string
				if rolesSlice, ok := rolesInterface.([]interface{}); ok {
					for _, r := range rolesSlice {
						if roleStr, ok := r.(string); ok {
							roles = append(roles, roleStr)
						}
					}
				}
				// Simpan ke context dengan nama yang dicari PermissionMiddleware
				c.Set("user_roles", roles) 
				fmt.Println("🔑 Roles found in token:", roles) // Debugging
			} else {
				fmt.Println("⚠️ No roles found in token claims!")
			}
		}

		// 6. Lanjutkan ke Handler berikutnya
		c.Next()
	}
}