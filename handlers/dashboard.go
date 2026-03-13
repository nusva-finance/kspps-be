package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardStats struct {
	TotalMembers      int64   `json:"total_members"`
	ActiveMembers     int64   `json:"active_members"`
	TotalSavings      float64 `json:"total_savings"`
	TotalLoans        float64 `json:"total_loans"`
	TotalLoansActive  int64   `json:"total_loans_active"`
	ThisMonthDeposits float64 `json:"this_month_deposits"`
	ThisMonthWithdraw float64 `json:"this_month_withdraw"`
}

type RecentTransaction struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Date      string    `json:"date"`
	Description string  `json:"description"`
}

type DashboardResponse struct {
	Stats              DashboardStats       `json:"stats"`
	RecentTransactions []RecentTransaction `json:"recent_transactions"`
}

func GetDashboard(c *gin.Context) {
	// TODO: Get actual data from database
	response := DashboardResponse{
		Stats: DashboardStats{
			TotalMembers:      125,
			ActiveMembers:     120,
			TotalSavings:      525000000,
			TotalLoans:        320000000,
			TotalLoansActive:  45,
			ThisMonthDeposits: 8500000,
			ThisMonthWithdraw: 1200000,
		},
		RecentTransactions: []RecentTransaction{
			{ID: 1, Type: "deposit", Amount: 500000, Date: "2026-03-03", Description: "Setoran Simpanan Wajib"},
			{ID: 2, Type: "credit", Amount: 2500000, Date: "2026-03-02", Description: "Pembayaran Angsuran #001"},
			{ID: 3, Type: "deposit", Amount: 1000000, Date: "2026-03-01", Description: "Setoran Simpanan Pokok"},
		},
	}

	c.JSON(http.StatusOK, response)
}
