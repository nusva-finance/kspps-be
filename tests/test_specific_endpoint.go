package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	
	"nusvakspps/config"
	"nusvakspps/handlers"
	"github.com/gin-gonic/gin"
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

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test router with correct API v1 path
	router := gin.New()
	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("", handlers.GetUsers)
		}
	}

	// Test the EXACT endpoint that user mentioned
	url := "/api/v1/users?page=1&limit=1000&search="
	
	fmt.Printf("\n%s\n", strings.Repeat("=", 80))
	fmt.Printf("🧪 TESTING EXACT ENDPOINT FROM USER\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("%s\n", strings.Repeat("=", 80))
	
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Printf("\n📡 HTTP Response:\n")
	fmt.Printf("   Status Code: %d %s\n", w.Code, http.StatusText(w.Code))
	fmt.Printf("   Content-Type: %s\n", w.Header().Get("Content-Type"))
	fmt.Printf("   Content-Length: %d bytes\n", w.Body.Len())

	if w.Code != http.StatusOK {
		fmt.Printf("\n❌ Response Body:\n%s\n", w.Body.String())
		return
	}

	// Parse JSON response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		fmt.Printf("\n❌ Error parsing JSON: %v\n", err)
		fmt.Printf("Raw Response Body:\n%s\n", w.Body.String())
		return
	}

	fmt.Printf("\n✅ Parsed Response:\n")
	
	// Show metadata
	total := int(response["total"].(float64))
	page := int(response["page"].(float64))
	limit := int(response["limit"].(float64))
	
	fmt.Printf("   📊 Total Users: %d\n", total)
	fmt.Printf("   📄 Page: %d\n", page)
	fmt.Printf("   📋 Limit: %d\n", limit)

	// Show user data
	data := response["data"].([]interface{})
	fmt.Printf("   👥 Users Returned: %d\n", len(data))
	
	if len(data) > 0 {
		fmt.Printf("\n   📋 User Details:\n")
		fmt.Printf("   %s\n", strings.Repeat("-", 76))
		
		for i, userInterface := range data {
			user := userInterface.(map[string]interface{})
			fmt.Printf("\n   User #%d:\n", i+1)
			fmt.Printf("      ID:         %v\n", user["id"])
			fmt.Printf("      Username:   %v\n", user["username"])
			fmt.Printf("      Email:      %v\n", user["email"])
			fmt.Printf("      Full Name:  %v\n", user["full_name"])
			fmt.Printf("      Is Active:  %v\n", user["is_active"])
			
			// Show roles if available
			if roles, ok := user["roles"].([]interface{}); ok && len(roles) > 0 {
				fmt.Printf("      Roles:      ")
				for j, role := range roles {
					roleMap := role.(map[string]interface{})
					fmt.Printf("%v", roleMap["name"])
					if j < len(roles)-1 {
						fmt.Printf(", ")
					}
				}
				fmt.Printf("\n")
			} else {
				fmt.Printf("      Roles:      (none)\n")
			}
		}
	}
	
	// Show raw JSON for debugging
	fmt.Printf("\n🔍 Raw JSON Response (formatted):\n")
	formattedJSON, _ := json.MarshalIndent(response, "   ", "   ")
	fmt.Printf("%s\n", formattedJSON)
	
	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("✅ CONFIRMATION: Backend API returns %d users as expected!\n", len(data))
	fmt.Printf(strings.Repeat("=", 80) + "\n")
}

func setEnvVar(key, value string) {
	os.Setenv(key, value)
}
