package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"nusvakspps/config"
	"nusvakspps/models"
	"nusvakspps/repositories"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID          uint     `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	FullName    string   `json:"full_name"`
	IsActive    bool     `json:"is_active"`
	ForceChange bool     `json:"force_change"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

var (
	jwtSecret = []byte("nusvakspps-secret-key-2024-jwt-token")

	userRepoOnce  sync.Once
	userRepoValue *repositories.UserRepository

	auditRepoOnce  sync.Once
	auditRepoValue *repositories.AuditRepository
)

func getUserRepo() *repositories.UserRepository {
	userRepoOnce.Do(func() {
		userRepoValue = repositories.NewUserRepository()
	})
	return userRepoValue
}

func getAuditRepo() *repositories.AuditRepository {
	auditRepoOnce.Do(func() {
		auditRepoValue = repositories.NewAuditRepository()
	})
	return auditRepoValue
}

const maxFailedLogin = 5

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- LOG TRACE START ---
	fmt.Printf("\n--- 🛡️ DEBUG LOGIN ATTEMPT ---\n")
	fmt.Printf("👤 Username: [%s]\n", req.Username)

	user, err := getUserRepo().FindByUsername(req.Username)
	if err != nil {
		fmt.Printf("❌ User [%s] tidak ditemukan!\n", req.Username)
		logAudit(c, 0, "login", "auth", nil, nil, "failed", "User not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 1. CEK APAKAH AKUN TERKUNCI
	if user.IsLocked {
		fmt.Println("🚫 Login Ditolak: Akun dalam status TERKUNCI")
		logAudit(c, user.ID, "login", "auth", nil, nil, "failed", "Account is locked")
		c.JSON(http.StatusLocked, gin.H{"error": "Account is locked. Please contact administrator."})
		return
	}

	// 2. CEK PASSWORD
	if !checkPassword(req.Password, user.PasswordHash) {
		fmt.Println("❌ Password Salah!")
		user.FailedLogin++
		if user.FailedLogin >= maxFailedLogin {
			user.IsLocked = true
			fmt.Println("⚠️ AKUN OTOMATIS TERKUNCI (Max login reached)")
		}
		getUserRepo().Update(user)

		logAudit(c, user.ID, "login", "auth", nil, nil, "failed", "Invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 3. JIKA BERHASIL: Reset Failed Login & Update Last Login
	fmt.Println("✅ LOGIN BERHASIL! Menyiapkan Data...")
	if user.FailedLogin > 0 {
		user.FailedLogin = 0
		getUserRepo().Update(user)
	}

	now := time.Now()
	user.LastLogin = &now
	user.LastIP = c.ClientIP()
	getUserRepo().Update(user)

	// 4. AMBIL ROLES & PERMISSIONS (Dilakukan SEBELUM generate token)
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	permissions := getUserPermissions(user.ID)
	fmt.Printf("🎭 Roles: %v | Permissions Count: %d\n", roles, len(permissions))

	// 5. GENERATE TOKENS (Sekarang variabel roles sudah tersedia)
	token, err := generateToken(user.ID, user.Username, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	fmt.Printf("--- 🛡️ DEBUG LOGIN END ---\n\n")
	logAudit(c, user.ID, "login", "auth", nil, nil, "success", "")

	c.JSON(http.StatusOK, LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FullName:    user.FullName,
			IsActive:    user.IsActive,
			ForceChange: user.ForceChange,
			Roles:       roles,
			Permissions: permissions,
		},
	})
}

func Logout(c *gin.Context) {
	userID := c.GetUint("current_user_id")
	logAudit(c, userID, "logout", "auth", nil, nil, "success", "")
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func RefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"token": "new-mock-token",
	})
}

// Generate token dengan menyertakan roles ke dalam claims
func generateToken(userID uint, username string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"roles":    roles, // <--- Data roles masuk ke payload JWT
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func generateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserPermissions(userID uint) []string {
	var permissions []string
	// Menggunakan pluck untuk mengambil kolom code dari tabel permissions
	err := config.GetDB().Table("permissions p").
		Select("DISTINCT p.code").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN user_roles ur ON rp.role_id = ur.role_id").
		Where("ur.user_id = ? AND rp.is_allowed = ?", userID, true).
		Pluck("code", &permissions).Error

	if err != nil {
		fmt.Printf("Error fetching permissions for user %d: %v\n", userID, err)
		return []string{}
	}
	return permissions
}

func logAudit(c *gin.Context, userID uint, action, module string, oldData, newData interface{}, status, errorMsg string) {
	var uid *uint
	if userID > 0 {
		uid = &userID
	}

	audit := &models.AuditLog{
		UserID:    uid,
		Username:  c.GetString("current_user_name"),
		Action:    action,
		Module:    module,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Status:    status,
		ErrorMsg:  errorMsg,
	}

	err := getAuditRepo().Create(audit)
	if err != nil {
		log.Printf("Failed to log audit: %v", err)
	}
}