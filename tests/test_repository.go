package main

import (
	"fmt"
	"log"
	"os"
	
	"nusvakspps/config"
	"nusvakspps/repositories"
)

func main() {
	// Load environment variables from .env file
	setEnvVar("DB_HOST", "103.127.98.221")
	setEnvVar("DB_PORT", "1482")
	setEnvVar("DB_USER", "userkoperasi")
	setEnvVar("DB_PASSWORD", "nusva12345")
	setEnvVar("DB_NAME", "nusvakspps")

	// Initialize database connection
	if err := config.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Test repository List function
	fmt.Println("🧪 Testing UserRepository.List() function...")
	fmt.Println("==========================================")
	
	userRepo := repositories.NewUserRepository()
	
	// Test with different limits
	for _, limit := range []int{10, 100, 1000} {
		fmt.Printf("\n📊 Testing with limit=%d:\n", limit)
		users, total, err := userRepo.List(0, limit)
		if err != nil {
			log.Printf("❌ Error with limit=%d: %v\n", limit, err)
			continue
		}
		fmt.Printf("✅ Total records in database: %d\n", total)
		fmt.Printf("✅ Records returned: %d\n", len(users))
		
		if len(users) > 0 {
			fmt.Printf("First user: ID=%d, Username=%s\n", users[0].ID, users[0].Username)
			if len(users) > 1 {
				fmt.Printf("Last user: ID=%d, Username=%s\n", users[len(users)-1].ID, users[len(users)-1].Username)
			}
		}
	}
	
	// Test pagination
	fmt.Println("\n📄 Testing pagination:")
	for page := 1; page <= 3; page++ {
		offset := (page - 1) * 2
		fmt.Printf("Page %d (offset=%d, limit=2):\n", page, offset)
		users, total, err := userRepo.List(offset, 2)
		if err != nil {
			log.Printf("❌ Error on page %d: %v\n", page, err)
			continue
		}
		fmt.Printf("  Returned %d users (total: %d)\n", len(users), total)
		for _, user := range users {
			fmt.Printf("    - ID=%d, Username=%s\n", user.ID, user.Username)
		}
	}
}

func setEnvVar(key, value string) {
	os.Setenv(key, value)
}
