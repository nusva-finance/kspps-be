package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"nusvakspps/models"
	"nusvakspps/repositories"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID          uint     `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	FullName    string   `json:"full_name"`
	IsActive    bool     `json:"is_active"`
	ForceChange  bool     `json:"force_change"`
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

	// Find user
	user, err := getUserRepo().FindByUsername(req.Username)
	if err != nil {
		logAudit(c, 0, "login", "auth", nil, nil, "failed", "User not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check if user is locked
	if user.IsLocked {
		logAudit(c, user.ID, "login", "auth", nil, nil, "failed", "Account is locked")
		c.JSON(http.StatusLocked, gin.H{"error": "Account is locked. Please contact administrator."})
		return
	}

	// Check password
	if !checkPassword(req.Password, user.PasswordHash) {
		// Increment failed login count
		user.FailedLogin++
		if user.FailedLogin >= maxFailedLogin {
			user.IsLocked = true
		}
		getUserRepo().Update(user)

		logAudit(c, user.ID, "login", "auth", nil, nil, "failed", "Invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Reset failed login count on successful login
	if user.FailedLogin > 0 {
		user.FailedLogin = 0
		getUserRepo().Update(user)
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	user.LastIP = c.ClientIP()
	getUserRepo().Update(user)

	// Generate token
	token, err := generateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Generate refresh token
	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Get roles and permissions
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	permissions := getUserPermissions(user.ID)

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
	userID := c.GetUint("user_id")
	logAudit(c, userID, "logout", "auth", nil, nil, "success", "")

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func RefreshToken(c *gin.Context) {
	// For simplicity, just generate a new token
	// In production, validate refresh token
	c.JSON(http.StatusOK, gin.H{
		"token": "new-mock-token",
	})
}

func generateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserPermissions(userID uint) []string {
	// TODO: Implement permission retrieval from database
	// For now, return all permissions for admin
	return []string{
		"dashboard.view",
		"user_mgmt.view",
		"user_mgmt.create",
		"user_mgmt.edit",
		"user_mgmt.delete",
		"member_mgmt.view",
		"member_mgmt.create",
		"member_mgmt.edit",
		"member_mgmt.delete",
		"saving_mgmt.view",
		"saving_mgmt.create",
		"saving_mgmt.edit",
		"saving_mgmt.delete",
		"loan_mgmt.view",
		"loan_mgmt.create",
		"loan_mgmt.edit",
		"loan_mgmt.delete",
		"loan_mgmt.approve",
		"role_mgmt.view",
		"audit.view",
	}
}

func logAudit(c *gin.Context, userID uint, action, module string, oldData, newData interface{}, status, errorMsg string) {
	var uid *uint
	if userID > 0 {
		uid = &userID
	}

	audit := &models.AuditLog{
		UserID:    uid,
		Username:  c.GetString("username"),
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
