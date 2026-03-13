package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🧪 TESTING MEMBER API ENDPOINTS")
	fmt.Println(strings.Repeat("=", 80))

	// Test 1: GET Members
	fmt.Println("\n🧪 Test 1: GET Members")
	getRequest, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/members", nil)
	client := http.Client{Timeout: 10 * time.Second}
	startTime := time.Now()
	
	response, err := client.Do(getRequest)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer response.Body.Close()

	duration := time.Since(startTime)
	fmt.Printf("📡 Response: Status: %d %s, Duration: %v\n", response.StatusCode, http.StatusText(response.StatusCode), duration)
	
	bodyBytes, _ := io.ReadAll(response.Body)
	fmt.Printf("Body Length: %d bytes\n", len(bodyBytes))
	fmt.Printf("Raw Response: %s\n", string(bodyBytes))

	var getData interface{}
	json.Unmarshal(bodyBytes, &getData)
	formatted, _ := json.MarshalIndent(getData, "   ", "   ")
	fmt.Printf("Parsed Response:\n%s\n", formatted)

	// Test 2: POST Create Member
	fmt.Println("\n🧪 Test 2: POST Create Member")
	createData := make(map[string]interface{})
	createData["full_name"] = "API Test Member 2"
	createData["gender"] = "Laki-laki"
	createData["join_date"] = "2025-01-15"
	createData["birth_date"] = "1990-01-01"
	createData["birth_place"] = "Jakarta"
	createData["nik"] = "3201010101019988"
	createData["address"] = "Test Address 456"
	createData["city"] = "Jakarta"
	createData["province"] = "DKI Jakarta"
	createData["postal_code"] = "10110"
	createData["phone_number"] = "081234567892"
	createData["email"] = "apitest2@example.com"

	jsonBody, _ := json.Marshal(createData)
	postRequest, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/members", bytes.NewBuffer(jsonBody))
	postRequest.Header.Set("Content-Type", "application/json")
	
	startTime = time.Now()
	postResponse, err := client.Do(postRequest)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer postResponse.Body.Close()

	duration = time.Since(startTime)
	fmt.Printf("📡 Response: Status: %d %s, Duration: %v\n", postResponse.StatusCode, http.StatusText(postResponse.StatusCode), duration)
	
	bodyBytes, _ := io.ReadAll(postResponse.Body)
	fmt.Printf("Body Length: %d bytes\n", len(bodyBytes))
	fmt.Printf("Raw Response: %s\n", string(bodyBytes))

	var postData interface{}
	json.Unmarshal(bodyBytes, &postData)
	formatted, _ := json.MarshalIndent(postData, "   ", "   ")
	fmt.Printf("Parsed Response:\n%s\n", formatted)

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("📊 FINAL DIAGNOSIS")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("If backend is updated, these tests should show:")
	fmt.Println("1. GET /api/v1/members -> Returns members from database")
	fmt.Println("2. POST /api/v1/members -> Creates member and returns success")
	fmt.Println("\nIf frontend shows empty member list:")
	fmt.Println("Check browser network tab for API errors")
	fmt.Println("Check browser console for frontend errors")
	fmt.Println("Check backend console for handler logs")
}
