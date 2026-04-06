package routes

import (
	"nusvakspps/handlers"
	"nusvakspps/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
			auth.POST("/refresh", handlers.RefreshToken)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Dashboard
			protected.GET("/dashboard", handlers.GetDashboard)

			// Users
			users := protected.Group("/users")
			{
				users.GET("", handlers.GetUsers)
				users.GET("/:id", handlers.GetUserByID)
				users.POST("", handlers.CreateUser)
				users.PUT("/:id", handlers.UpdateUser)
				users.DELETE("/:id", handlers.DeleteUser)
			}

			// Members
			members := protected.Group("/members")
			{
				members.GET("", handlers.GetMembers)
				members.GET("/:id", handlers.GetMemberByID)
				members.POST("", handlers.CreateMember)
				members.PUT("/:id", handlers.UpdateMember)
				members.DELETE("/:id", handlers.DeleteMember)
			}

			// Savings
			savings := protected.Group("/savings")
			{
				// Saving Types (dynamic configuration)
				savingTypes := savings.Group("/types")
				{
					savingTypes.GET("", handlers.GetSavingTypes)
					savingTypes.GET("/:id", handlers.GetSavingTypeByID)
					savingTypes.POST("", handlers.CreateSavingType)
					savingTypes.PUT("/:id", handlers.UpdateSavingType)
					savingTypes.DELETE("/:id", handlers.DeleteSavingType)
					savingTypes.POST("/initialize", handlers.InitializeSavingTypes)
				}

				// Saving Accounts
				savings.GET("/accounts", handlers.GetAllSavingsAccounts)

			// Saving Types with Total Balances
				savingsTypesGroup := savings.Group("/types")
				{
					savingsTypesGroup.GET("/balances", handlers.GetAllSavingTypesWithBalances)
				}

			// Saving Transactions
				savings.GET("/transactions", handlers.GetSavingsTransactions)
				savings.GET("/transactions/all", handlers.GetAllTransactionsList)
				savings.GET("/transactions/:id", handlers.GetSavingsTransactionByID)
				savings.POST("/transactions", handlers.CreateSavingsTransaction)
				savings.PUT("/transactions/:id", handlers.UpdateSavingsTransaction)
				savings.DELETE("/transactions/:id", handlers.DeleteSavingsTransaction)

				// Member Balances
				savings.GET("/member/:memberId/balance", handlers.GetMemberSavingsBalance)

				// Legacy endpoints (for backward compatibility)
				savings.GET("", handlers.GetSavingsTransactions)
				savings.GET("/:id", handlers.GetSavingsTransactionByID)
				savings.POST("", handlers.CreateSavingsTransaction)
				savings.PUT("/:id", handlers.UpdateSavingsTransaction)
				savings.DELETE("/:id", handlers.DeleteSavingsTransaction)
			}


			// Security
			security := protected.Group("/security")
			{
				security.GET("/roles", handlers.GetRoles)
				security.POST("/roles", handlers.CreateRole)
				security.GET("/menus", handlers.GetMenus)
				security.GET("/permissions", handlers.GetPermissions)
				security.POST("/roles/:role-id/permissions", handlers.AssignPermissions)
				security.GET("/audit-logs", handlers.GetAuditLogs)
			}

			margin := protected.Group("/margin-setups")
			{
				margin.GET("", handlers.GetMargins)
				margin.GET("/:id", handlers.GetMarginByID)
				margin.POST("", handlers.CreateMargin)
				margin.PUT("/:id", handlers.UpdateMargin)
				margin.DELETE("/:id", handlers.DeleteMargin)
			}

			// Kategori Barang
			kategoriBarang := protected.Group("/kategori-barangs")
			{
				kategoriBarang.GET("", handlers.GetKategoriBarangs)
				kategoriBarang.GET("/:id", handlers.GetKategoriBarangByID)
				kategoriBarang.POST("", handlers.CreateKategoriBarang)
				kategoriBarang.PUT("/:id", handlers.UpdateKategoriBarang)
				kategoriBarang.DELETE("/:id", handlers.DeleteKategoriBarang)
			}

			// Rekening
			rekening := protected.Group("/rekening")
			{
				rekening.GET("", handlers.GetRekenings)
				rekening.GET("/:id", handlers.GetRekeningByID)
				rekening.GET("/:id/mutasi", handlers.GetMutasiRekening)
				rekening.POST("", handlers.CreateRekening)
				rekening.PUT("/:id", handlers.UpdateRekening)
				rekening.DELETE("/:id", handlers.DeleteRekening)
			}

			// Pembiayaan
			pembiayaan := protected.Group("/pembiayaan")
			{
				pembiayaan.GET("", handlers.GetPembiayaan)
				pembiayaan.GET("/margin", handlers.GetMarginByCategoryAndTenor)
				pembiayaan.GET("/:id", handlers.GetPembiayaanByID)
				pembiayaan.GET("/:id/pembayaran", handlers.GetPembayaranByPinjamanID)
				pembiayaan.GET("/:id/angsuranke", handlers.GetAngsuranKe)
				pembiayaan.POST("", handlers.CreatePembiayaan)
				pembiayaan.PUT("/:id", handlers.UpdatePembiayaan)
				pembiayaan.DELETE("/:id", handlers.DeletePembiayaan)
			}

			// Pembayaran Pembiayaan (standalone endpoints)
			pembayaran := protected.Group("/pembayaran-pembiayaan")
			{
				pembayaran.GET("/:id", handlers.GetPembayaranByID)
				pembayaran.POST("", handlers.CreatePembayaran)
				pembayaran.PUT("/:id", handlers.UpdatePembayaran)
				pembayaran.DELETE("/:id", handlers.DeletePembayaran)
			}

		}
	}
}
