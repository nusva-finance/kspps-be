package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	
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

	// Test the correct API path that frontend would use
	fmt.Println("\n🧪 Testing frontend API path: /api/v1/users")
	fmt.Println("==========================================")
	
	testCases := []struct {
		name string
		url  string
	}{
		{"Default (limit=10)", "/api/v1/users"},
		{"Limit 1000", "/api/v1/users?limit=1000"},
		{"Page 1 Limit 3", "/api/v1/users?page=1&limit=3"},
		{"Page 2 Limit 3", "/api/v1/users?page=2&limit=3"},
	}

	for _, tc := range testCases {
		fmt.Printf("\n🧪 Testing: %s\n", tc.name)
		fmt.Println("==========================================")
		
		req, _ := http.NewRequest("GET", tc.url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		fmt.Printf("Status Code: %d\n", w.Code)
		
		if w.Code != http.StatusOK {
			fmt.Printf("Response Body: %s\n", w.Body.String())
			continue
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			fmt.Printf("❌ Error parsing response: %v\n", err)
			continue
		}

		total := int(response["total"].(float64))
		data := response["data"].([]interface{})
		limit := int(response["limit"].(float64))
		page := int(response["page"].(float64))
		
		fmt.Printf("Page: %d, Limit: %d, Total records: %d\n", page, limit, total)
		fmt.Printf("Records returned: %d\n", len(data))
		
		if len(data) > 0 {
			firstUser := data[0].(map[string]interface{})
			fmt.Printf("First user: ID=%v, Username=%v, Email=%v\n", 
				firstUser["id"], firstUser["username"], firstUser["email"])
			if len(data) > 1 {
				lastUser := data[len(data)-1].(map[string]interface{})
				fmt.Printf("Last user: ID=%v, Username=%v, Email=%v\n", 
					lastUser["id"], lastUser["username"], lastUser["email"])
			}
		}
	}
	
	// Test the old incorrect path that frontend might be using
	fmt.Println("\n⚠️  Testing incorrect path: /api/v1/users/ (with trailing slash)")
	fmt.Println("==========================================")
	req, _ := http.NewRequest("GET", "/api/v1/users/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	fmt.Printf("Status Code: %d\n", w.Code)
	fmt.Printf("Response Body: %s\n", w.Body.String())
}

func setEnvVar(key, value string) {
	os.Setenv(key, value)
}
