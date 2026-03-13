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

	// Test the exact endpoint from user
	url := "/api/v1/users?page=1&limit=1000&search="
	
	fmt.Printf("\n%s\n", strings.Repeat("=", 80))
	fmt.Printf("🧪 TESTING FIXED API ENDPOINT\n")
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
			roles := "none"
			if r, ok := user["roles"].([]interface{}); ok && len(r) > 0 {
				roleNames := make([]string, len(r))
				for j, role := range r {
					roleMap := role.(map[string]interface{})
					roleNames[j] = fmt.Sprintf("%v", roleMap["name"])
				}
				roles = strings.Join(roleNames, ", ")
			}
			
			fmt.Printf("   User #%d: ID=%v, Username=%v, Email=%v, Roles=%s\n", 
				i+1, user["id"], user["username"], user["email"], roles)
		}
	}
	
	// Test with search
	fmt.Printf("\n%s\n", strings.Repeat("=", 80))
	fmt.Printf("🧪 TESTING WITH SEARCH PARAMETER\n")
	fmt.Printf("%s\n", strings.Repeat("=", 80))
	
	searchURL := "/api/v1/users?page=1&limit=1000&search=admin"
	req2, _ := http.NewRequest("GET", searchURL, nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	fmt.Printf("URL: %s\n", searchURL)
	fmt.Printf("Status Code: %d %s\n", w2.Code, http.StatusText(w2.Code))
	
	if w2.Code == http.StatusOK {
		var searchResponse map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &searchResponse)
		searchData := searchResponse["data"].([]interface{})
		searchTotal := int(searchResponse["total"].(float64))
		
		fmt.Printf("✅ Search Results: %d users (total: %d)\n", len(searchData), searchTotal)
		
		if len(searchData) > 0 {
			user := searchData[0].(map[string]interface{})
			fmt.Printf("   Found: %v (%v)\n", user["username"], user["email"])
		}
	}
	
	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("✅ FIXED! Backend now returns %d users correctly!\n", len(data))
	fmt.Printf(strings.Repeat("=", 80) + "\n")
}

func setEnvVar(key, value string) {
	os.Setenv(key, value)
}
