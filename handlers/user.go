package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"nusvakspps/models"
	"nusvakspps/repositories"
)

type CreateUserRequest struct {
	Username        string   `json:"username" binding:"required"`
	Email           string   `json:"email" binding:"required,email"`
	FullName        string   `json:"full_name" binding:"required"`
	Password        string   `json:"password" binding:"required"`
	ConfirmPassword string   `json:"confirm_password" binding:"required"`
	Roles           []string `json:"roles" binding:"required,min=1"`
	IsActive        bool     `json:"is_active"`
}

type UpdateUserRequest struct {
	FullName *string  `json:"full_name"`
	Email    *string  `json:"email"`
	IsActive *bool    `json:"is_active"`
	Roles   []string `json:"roles"`
	Password *string  `json:"password"`
}

func GetUsers(c *gin.Context) {
	page := 1
	limit := 10
	search := ""

	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if s := c.Query("search"); s != "" {
		search = s
	}

	fmt.Println("📋 Getting users, page:", page, "limit:", limit, "search:", search)

	userRepo := repositories.NewUserRepository()
	offset := (page - 1) * limit

	users, total, err := userRepo.List(offset, limit, search)
	if err != nil {
		fmt.Println("❌ Error getting users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	fmt.Println("✅ Retrieved", len(users), "users, total:", total)

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	fmt.Println("👤 Getting user by ID:", id)

	userRepo := repositories.NewUserRepository()
	user, err := userRepo.FindByID(id)
	if err != nil {
		fmt.Println("❌ User not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	fmt.Println("✅ Found user:", user.Username)

	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {

	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Validation error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("📤 Creating user:", req.Username)

	// Validate passwords match
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// Check if username already exists
	userRepo := repositories.NewUserRepository()
	if _, err := userRepo.FindByUsername(req.Username); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Check if email already exists
	if _, err := userRepo.FindByEmail(req.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Start transaction for atomic user creation
	tx := repositories.GetDB().Begin()
	if tx.Error != nil {
		fmt.Println("❌ Error starting transaction:", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Create user
	user := &models.User{
		Username:    req.Username,
		Email:       req.Email,
		FullName:    req.FullName,
		PasswordHash:    string(hashedPassword),
		IsActive:    req.IsActive,
		ForceChange: true,
		// --- TAMBAHKAN INI: Isi audit trail ---
		CreatedBy:   operatorName,
		UpdatedBy:   operatorName,
	}

	if err := userRepo.Create(user); err != nil {
		fmt.Println("❌ Error creating user:", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Assign roles AFTER user is created
	if len(req.Roles) > 0 {
		roleRepo := repositories.NewRoleRepository()

		for _, roleName := range req.Roles {
			role, err := roleRepo.FindByName(roleName)
			if err != nil {
				fmt.Printf("❌ Role not found: %s, error: %v\n", roleName, err)
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Role '%s' not found", roleName)})
				return
			}

			fmt.Printf("🔗 Assigning role %s (ID: %d) to user %d\n", role.Name, role.ID, user.ID)

			if err := userRepo.AssignRole(user.ID, role.ID); err != nil {
				fmt.Printf("❌ Error assigning role: %v\n", err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
				return
			} else {
				fmt.Printf("✅ Role %s assigned successfully\n", role.Name)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		fmt.Println("❌ Error committing transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit changes"})
		return
	}

	fmt.Println("✅ Transaction committed successfully")

	// Reload user with roles for response
	fmt.Println("🔄 Reloading user with roles for response")
	user, _ = userRepo.FindByID(user.ID)

	fmt.Println("✅ User created successfully with ID:", user.ID)
	fmt.Printf("📊 Final user state: ID=%d, Username=%s, FullName=%s, Email=%s, IsActive=%v, Roles=%d\n",
		user.ID, user.Username, user.FullName, user.Email, user.IsActive, len(user.Roles))

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

func UpdateUser(c *gin.Context) {
	operatorName := c.GetString("current_user_name")
	
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("📝 Updating user:", id, "with data:", req)

	userRepo := repositories.NewUserRepository()
	user, err := userRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Start transaction for atomic updates
	tx := repositories.GetDB().Begin()
	if tx.Error != nil {
		fmt.Println("❌ Error starting transaction:", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// =============================================================
	// 1. HANDLE ROLE UPDATES
	// =============================================================
	if len(req.Roles) > 0 {
		roleRepo := repositories.NewRoleRepository()
		
		// Remove all old roles
		for _, role := range user.Roles {
			if err := userRepo.RemoveRole(id, role.ID); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove old roles"})
				return
			}
		}

		// Add new roles
		for _, roleName := range req.Roles {
			role, err := roleRepo.FindByName(roleName)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Role '%s' not found", roleName)})
				return
			}
			if err := userRepo.AssignRole(id, role.ID); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
				return
			}
		}
	}

	// =============================================================
	// 2. HANDLE PASSWORD UPDATE (THE MISSING PIECE)
	// =============================================================
	if req.Password != nil && *req.Password != "" {
		fmt.Println("🔐 Password baru terdeteksi, melakukan hashing...")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), 10)
		if err != nil {
			fmt.Println("❌ Error hashing password:", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			return
		}
		// Update hash dan otomatis buka kunci akun jika sedang locked
		user.PasswordHash = string(hashedPassword)
		user.IsLocked = false
		user.FailedLogin = 0
		fmt.Println("✅ Password hash updated successfully")
	}

	// =============================================================
	// 3. UPDATE BASIC FIELDS
	// =============================================================
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	
	if req.Email != nil {
		// Check if email already exists (excluding current user)
		if existing, _ := userRepo.FindByEmail(*req.Email); existing != nil && existing.ID != id {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		user.Email = *req.Email
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	user.UpdatedBy = operatorName

	// Save the updated user object
	if err := userRepo.Update(user); err != nil {
		fmt.Println("❌ Error updating user basic fields:", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		fmt.Println("❌ Error committing transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit changes"})
		return
	}

	// Reload for final response
	user, _ = userRepo.FindByID(id)

	fmt.Println("✅ User updated successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"id":      id,
		"data":    user,
	})

	if req.Password != nil {
    	fmt.Printf("🔑 PASSWORD DITERIMA DARI FRONTEND: [%s]\n", *req.Password)
	} else {
    	fmt.Println("⚠️ PASSWORD TIDAK DITERIMA (NIL). Cek nama field di Frontend!")
	}

}


func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	fmt.Println("🗑️ Deleting user:", id)

	userRepo := repositories.NewUserRepository()
	if err := userRepo.Delete(id); err != nil {
		fmt.Println("❌ Error deleting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	fmt.Println("✅ User deleted successfully:", id)

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
