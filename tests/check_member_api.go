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
	getReq, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/members", nil)
	client := &http.Client{Timeout: 10 * time.Second}
	startTime := time.Now()
	
	resp, err := client.Do(getReq)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)
	fmt.Printf("📡 Response: Status: %d %s, Duration: %v\n", resp.StatusCode, http.StatusText(resp.StatusCode), duration)
	
	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Printf("Body Length: %d bytes\n", len(bodyBytes))
	fmt.Printf("Raw Response: %s\n", string(bodyBytes))

	var getResp interface{}
	json.Unmarshal(bodyBytes, &getResp)
	formatted, _ := json.MarshalIndent(getResp, "   ", "   ")
	fmt.Printf("Parsed Response:\n%s\n", formatted)

	// Test 2: POST Create Member
	fmt.Println("\n🧪 Test 2: POST Create Member")
	createData := map[string]interface{}{
		"full_name": "API Test Member 2",
		"gender": "Laki-laki",
		"join_date": "2025-01-15",
		"birth_date": "1990-01-01",
		"birth_place": "Jakarta",
		"nik": "3201010101019988",
		"address": "Test Address 456",
		"city": "Jakarta",
		"province": "DKI Jakarta",
		"postal_code": "10110",
		"phone_number": "081234567892",
		"email": "apitest2@example.com",
	}
	
	jsonBody, _ := json.Marshal(createData)
	postReq, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/members", bytes.NewBuffer(jsonBody))
	postReq.Header.Set("Content-Type", "application/json")
	
	startTime = time.Now()
	resp, err = client.Do(postReq)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	duration = time.Since(startTime)
	fmt.Printf("📡 Response: Status: %d %s, Duration: %v\n", resp.StatusCode, http.StatusText(resp.StatusCode), duration)
	
	bodyBytes, _ = io.ReadAll(resp.Body)
	fmt.Printf("Body Length: %d bytes\n", len(bodyBytes))
	fmt.Printf("Raw Response: %s\n", string(bodyBytes))

	var postResp interface{}
	json.Unmarshal(bodyBytes, &postResp)
	formatted, _ := json.MarshalIndent(postResp, "   ", "   ")
	fmt.Printf("Parsed Response:\n%s\n", formatted)

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("📊 DIAGNOSIS")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("Check the backend console logs for any errors or issues")
}
