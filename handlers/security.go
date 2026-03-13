package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

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
	// TODO: Implement menu retrieval
	c.JSON(http.StatusOK, gin.H{
		"data":  []string{},
		"total": 0,
	})
}

func GetPermissions(c *gin.Context) {
	// TODO: Implement permissions retrieval
	c.JSON(http.StatusOK, gin.H{
		"data":  []string{},
		"total": 0,
	})
}

func AssignPermissions(c *gin.Context) {
	// TODO: Implement permission assignment
	c.JSON(http.StatusOK, gin.H{
		"message": "Permissions will be implemented",
	})
}

func GetAuditLogs(c *gin.Context) {
	// TODO: Implement audit logs retrieval
	c.JSON(http.StatusOK, gin.H{
		"data":  []string{},
		"total": 0,
	})
}
