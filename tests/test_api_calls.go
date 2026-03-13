package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	
	"nusvakspps/config"
)

type CreateMemberRequest struct {
	MemberID   string  `json:"member_id" binding:"required"`
	FullName    string  `json:"full_name" binding:"required"`
	PhoneNumber string  `json:"phone_number" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	Address     string  `json:"address" binding:"required"`
	IsActive    *bool   `json:"is_active"`
}

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

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🧪 TESTING API CALLS TO LOCALHOST")
	fmt.Println(strings.Repeat("=", 80))

	// Test different API endpoints
	testCases := []struct {
		name   string
		method string
		url    string
		body   interface{}
	}{
		{
			name:   "GET Users",
			method: "GET",
			url:    "http://localhost:8080/api/v1/users",
			body:   nil,
		},
		{
			name:   "GET Members",
			method: "GET",
			url:    "http://localhost:8080/api/v1/members",
			body:   nil,
		},
		{
			name:   "POST Create Member",
			method: "POST",
			url:    "http://localhost:8080/api/v1/members",
			body: CreateMemberRequest{
				MemberID:   "TEST001",
				FullName:    "Test Member",
				PhoneNumber: "081234567890",
				Email:       "testmember@example.com",
				Address:     "Test Address 123",
				IsActive:    boolPtr(true),
			},
		},
	}

	for _, tc := range testCases {
		fmt.Printf("\n🧪 Test: %s\n", tc.name)
		fmt.Printf("   Method: %s\n", tc.method)
		fmt.Printf("   URL: %s\n", tc.url)
		
		if tc.body != nil {
			fmt.Printf("   Body: %v\n", tc.body)
		}

		var req *http.Request
		var err error
		
		if tc.method == "GET" {
			req, err = http.NewRequest(tc.method, tc.url, nil)
		} else {
			jsonBody, _ := json.Marshal(tc.body)
			req, err = http.NewRequest(tc.method, tc.url, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
		}

		if err != nil {
			fmt.Printf("   ❌ Error creating request: %v\n", err)
			continue
		}

		// Set timeout and send request
		client := &http.Client{Timeout: 10 * time.Second}
		startTime := time.Now()
		
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("   ❌ Error sending request: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		duration := time.Since(startTime)
		
		fmt.Printf("   📡 Response:\n")
		fmt.Printf("      Status: %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
		fmt.Printf("      Duration: %v\n", duration)
		fmt.Printf("      Headers:\n")
		
		for key, values := range resp.Header {
			fmt.Printf("         %s: %v\n", key, values)
		}

		// Read response body
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("      ❌ Error reading body: %v\n", err)
			continue
		}

		fmt.Printf("      Body Length: %d bytes\n", len(bodyBytes))
		
		// Try to parse as JSON
		var jsonResponse interface{}
		if err := json.Unmarshal(bodyBytes, &jsonResponse); err == nil {
			fmt.Printf("      Parsed Response:\n")
			formattedJSON, _ := json.MarshalIndent(jsonResponse, "         ", "   ")
			fmt.Printf("%s\n", formattedJSON)
		} else {
			fmt.Printf("      Raw Response:\n")
			fmt.Printf("      %s\n", string(bodyBytes))
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 FRONTEND CONFIGURATION CHECK")
	fmt.Println(strings.Repeat("=", 80))
	
	// Check what frontend .env is using
	fmt.Println("Frontend .env should contain:")
	fmt.Println("   VITE_API_URL=http://localhost:8080/api/v1")
	fmt.Println("\nFrontend .env currently contains:")
	fmt.Println("   VITE_API_URL=https://barrel-saver-copy-speak.trycloudflare.com/api/v1")
	fmt.Println("\n❌ PROBLEM: Frontend is calling WRONG API!")
	fmt.Println("   Frontend: https://barrel-saver-copy-speak.trycloudflare.com/api/v1")
	fmt.Println("   Backend:  http://localhost:8080/api/v1")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("✅ SOLUTION: Change frontend .env to use localhost API")
}

func boolPtr(b bool) *bool {
	return &b
}

func setEnvVar(key, value string) {
	os.Setenv(key, value)
}
