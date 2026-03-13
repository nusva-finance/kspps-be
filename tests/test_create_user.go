package main

import (
	"fmt"
	"os"
	"strings"
	
	"nusvakspps/config"
	"nusvakspps/models"
	"nusvakspps/repositories"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load environment variables
	setEnvVar("DB_HOST", "103.127.98.221")
	setEnvVar("DB_PORT", "1482")
	setEnvVar("DB_USER", "userkoperasi")
	setEnvVar("DB_PASSWORD", "nusva12345")
	setEnvVar("DB_NAME", "nusvakspps")

	// Initialize database
	if err := config.InitDB(); err != nil {
		fmt.Println("❌ Failed to initialize database:", err)
		return
	}
	fmt.Println("✅ Database connected")

	userRepo := repositories.NewUserRepository()
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 CURRENT STATE BEFORE TEST")
	fmt.Println(strings.Repeat("=", 80))
	
	// Count users before
	var totalBefore int64
	err := config.GetDB().Model(&models.User{}).Count(&totalBefore).Error
	if err != nil {
		fmt.Println("❌ Error counting users:", err)
		return
	}
	fmt.Printf("📊 Total users before test: %d\n", totalBefore)
	
	// List users before
	usersBefore, _, err := userRepo.List(0, 100, "")
	if err != nil {
		fmt.Println("❌ Error listing users:", err)
		return
	}
	fmt.Printf("👥 Users returned before test: %d\n", len(usersBefore))
	
	if len(usersBefore) > 0 {
		fmt.Println("   Sample users:")
		for i, user := range usersBefore {
			if i >= 3 {
				fmt.Println("   ... (and more)")
				break
			}
			fmt.Printf("   - ID: %d, Username: %s, Email: %s\n", user.ID, user.Username, user.Email)
		}
	}
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📤 CREATING NEW USER")
	fmt.Println(strings.Repeat("=", 80))
	
	// Create a test user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("TestPassword123!"), 14)
	if err != nil {
		fmt.Println("❌ Error hashing password:", err)
		return
	}
	
	newUser := &models.User{
		Username:     "testuser999",
		Email:        "testuser999@example.com",
		FullName:     "Test User 999",
		PasswordHash: string(hashedPassword),
		IsActive:     true,
	}
	
	fmt.Printf("📝 Creating user: %s (%s)\n", newUser.Username, newUser.Email)
	fmt.Printf("   Full Name: %s\n", newUser.FullName)
	
	err = userRepo.Create(newUser)
	if err != nil {
		fmt.Println("❌ ❌ Error creating user:", err)
		fmt.Println("🔍 SQL Error Details:")
		fmt.Println(err)
		return
	}
	
	fmt.Printf("✅ User created successfully with ID: %d\n", newUser.ID)
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 CURRENT STATE AFTER CREATE")
	fmt.Println(strings.Repeat("=", 80))
	
	// Count users after
	var totalAfter int64
	err = config.GetDB().Model(&models.User{}).Count(&totalAfter).Error
	if err != nil {
		fmt.Println("❌ Error counting users:", err)
		return
	}
	fmt.Printf("📊 Total users after test: %d\n", totalAfter)
	fmt.Printf("📈 Users increased by: %d\n", totalAfter-totalBefore)
	
	// List users after create
	usersAfter, _, err := userRepo.List(0, 100, "")
	if err != nil {
		fmt.Println("❌ Error listing users:", err)
		return
	}
	fmt.Printf("👥 Users returned after test: %d\n", len(usersAfter))
	
	// Check if new user is in the list
	found := false
	for _, user := range usersAfter {
		if user.ID == newUser.ID {
			found = true
			fmt.Printf("✅ New user found in list! ID: %d, Username: %s\n", user.ID, user.Username)
			break
		}
	}
	
	if !found {
		fmt.Println("❌ ❌ New user NOT found in list!")
		fmt.Println("This is the bug - user was created but not returned by List()")
	}
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 TESTING SEARCH FUNCTIONALITY")
	fmt.Println(strings.Repeat("=", 80))
	
	// Test search for the new user
	searchResults, _, err := userRepo.List(0, 100, "testuser999")
	if err != nil {
		fmt.Println("❌ Error searching users:", err)
		return
	}
	fmt.Printf("🔍 Search results for 'testuser999': %d users\n", len(searchResults))
	
	if len(searchResults) > 0 {
		for _, user := range searchResults {
			fmt.Printf("   - Found: %s (%s)\n", user.Username, user.Email)
		}
	}
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✅ TEST SUMMARY")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Before: %d users | After: %d users | Increase: %d\n", totalBefore, totalAfter, totalAfter-totalBefore)
	fmt.Printf("List returns: %d users | Search returns: %d users\n", len(usersAfter), len(searchResults))
	
	if totalAfter > totalBefore && len(usersAfter) > len(usersBefore) {
		fmt.Println("✅ SUCCESS: User creation and listing working correctly!")
	} else {
		fmt.Println("❌ PROBLEM DETECTED: User creation or listing has issues!")
	}
}

func setEnvVar(key, value string) {
	os.Setenv(key, value)
}
