package handlers

import (
	"fmt"
	"net/http"
	"time" 

	"github.com/gin-gonic/gin"
	
)

type LoanApplicationRequest struct {
	MemberID        uint    `json:"member_id" binding:"required"`
	PrincipalAmount float64 `json:"principal_amount" binding:"required"`
	MarginRate      float64 `json:"margin_rate" binding:"required"`
	TermMonths      int     `json:"term_months" binding:"required"`
	StartDate       string  `json:"start_date" binding:"required"`
	Purpose         string  `json:"purpose"`
	ContractType    string  `json:"contract_type" binding:"required"`
}

func GetLoanApplications(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []map[string]interface{}{
			{
				"id":               1,
				"application_no":   "LA2503001",
				"member_name":      "Ahmad Rizki",
				"principal_amount": 5000000,
				"margin_rate":      15,
				"term_months":      12,
				"status":           "pending",
			},
		},
		"total": 1,
	})
}

func CreateLoanApplication(c *gin.Context) {
	// Ambil operator dari middleware
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	var req LoanApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Logging untuk memastikan data masuk
	fmt.Printf("📝 Loan Application created by: %s for Member ID: %d\n", operatorName, req.MemberID)

	// Note: Sementara response sukses, Kakak bisa tambahkan logic save ke DB di sini nanti
	c.JSON(http.StatusCreated, gin.H{
		"message": "Loan application created successfully",
		"created_by": operatorName,
	})
}

func ApproveLoan(c *gin.Context) {
	// Ambil operator dari middleware
	approverName := c.GetString("current_user_name")
	if approverName == "" {
		approverName = "system"
	}

	id := c.Param("id") // Ambil ID dari param rute

	// Perbaikan: Gunakan variabel 'id' yang sudah didefinisikan di atas, bukan 'idParam'
	c.JSON(http.StatusOK, gin.H{
		"message":     "Loan approved successfully",
		"id":          id, 
		"approved_by": approverName,
		"approved_at": time.Now().Format("2006-01-02 15:04:05"),
	})
}

func GetLoanSchedule(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"application_id": id,
		"schedules": []map[string]interface{}{
			{
				"sequence_number": 1,
				"due_date":        "2026-04-01",
				"principal":       416666.67,
				"margin":          62500,
				"total_amount":    479166.67,
				"is_paid":         false,
			},
		},
	})
}

func CreateLoanTransaction(c *gin.Context) {
	// Ambil operator dari middleware
	operatorName := c.GetString("current_user_name")
	if operatorName == "" {
		operatorName = "system"
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Loan transaction created successfully",
		"created_by": operatorName,
	})
}