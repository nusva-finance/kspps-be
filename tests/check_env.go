package main

import (
	"fmt"
	"os"
)

func main() {
	envVars := []string{
		"DB_HOST",
		"DB_PORT", 
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"SERVER_PORT",
	}

	fmt.Println("🔍 Checking Environment Variables:")
	fmt.Println("==========================================")
	
	for _, envVar := range envVars {
		value := os.Getenv(envVar)
		if value == "" {
			fmt.Printf("%s: (not set)\n", envVar)
		} else {
			// Hide password
			if envVar == "DB_PASSWORD" {
				fmt.Printf("%s: *** (hidden)\n", envVar)
			} else {
				fmt.Printf("%s: %s\n", envVar, value)
			}
		}
	}
	
	fmt.Println("\n📋 Current Working Directory:")
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("%s\n", cwd)
	}
	
	fmt.Println("\n🔍 Checking .env file existence:")
	envFile := ".env"
	if _, err := os.Stat(envFile); err == nil {
		fmt.Printf("✅ .env file exists at: %s/%s\n", cwd, envFile)
	} else {
		fmt.Printf("❌ .env file not found at: %s/%s\n", cwd, envFile)
		fmt.Println("This is the problem! Backend is using default config instead of .env")
	}
}
