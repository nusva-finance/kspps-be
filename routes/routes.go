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
		// --- AUTH ROUTES (Public) ---
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
			auth.POST("/refresh", handlers.RefreshToken)
		}

		// --- PROTECTED ROUTES (Harus Login) ---
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Dashboard (Asal login boleh lihat)
			protected.GET("/dashboard", handlers.GetDashboard)

			// ---------------------------------------------------------
			// 1. USERS & SECURITY
			// ---------------------------------------------------------
			users := protected.Group("/users")
			{
				users.GET("", middleware.PermissionMiddleware("users.view"), handlers.GetUsers)
				users.GET("/:id", middleware.PermissionMiddleware("users.view"), handlers.GetUserByID)
				users.POST("", middleware.PermissionMiddleware("users.create"), handlers.CreateUser)
				users.PUT("/:id", middleware.PermissionMiddleware("users.update"), handlers.UpdateUser)
				users.DELETE("/:id", middleware.PermissionMiddleware("users.delete"), handlers.DeleteUser)
			}

			security := protected.Group("/security")
			{
				security.GET("/roles", middleware.PermissionMiddleware("security.view"), handlers.GetRoles)
				security.POST("/roles", middleware.PermissionMiddleware("security.manage"), handlers.CreateRole)
				security.GET("/menus", middleware.PermissionMiddleware("security.view"), handlers.GetMenus)
				security.GET("/permissions", middleware.PermissionMiddleware("security.view"), handlers.GetPermissions)
				security.POST("/roles/:role-id/permissions", middleware.PermissionMiddleware("security.manage"), handlers.AssignPermissions)
				security.GET("/roles/:role-id/permissions", middleware.PermissionMiddleware("security.view"), handlers.GetRolePermissions)
				security.GET("/audit-logs", middleware.PermissionMiddleware("security.audit"), handlers.GetAuditLogs)
			}

			// ---------------------------------------------------------
			// 2. MEMBERS
			// ---------------------------------------------------------
			members := protected.Group("/members")
			{
				members.GET("", middleware.PermissionMiddleware("members.view"), handlers.GetMembers)
				members.GET("/:id", middleware.PermissionMiddleware("members.view"), handlers.GetMemberByID)
				members.POST("", middleware.PermissionMiddleware("members.create"), handlers.CreateMember)
				members.PUT("/:id", middleware.PermissionMiddleware("members.update"), handlers.UpdateMember)
				members.DELETE("/:id", middleware.PermissionMiddleware("members.delete"), handlers.DeleteMember)
				
				// Rute Import ditaruh di sini (di dalam group members)
				members.POST("/import", middleware.PermissionMiddleware("members.import"), handlers.ImportMembers)
			}

			// ---------------------------------------------------------
			// 3. SAVINGS (SIMPANAN)
			// ---------------------------------------------------------
			savings := protected.Group("/savings")
			{
				// Saving Types
				savingTypes := savings.Group("/types")
				{
					savingTypes.GET("", middleware.PermissionMiddleware("savings.config"), handlers.GetSavingTypes)
					savingTypes.POST("", middleware.PermissionMiddleware("savings.config"), handlers.CreateSavingType)
					savingTypes.PUT("/:id", middleware.PermissionMiddleware("savings.config"), handlers.UpdateSavingType)
					savingTypes.DELETE("/:id", middleware.PermissionMiddleware("savings.config"), handlers.DeleteSavingType)
					savingTypes.POST("/initialize", middleware.PermissionMiddleware("savings.config"), handlers.InitializeSavingTypes)
					savingTypes.GET("/balances", middleware.PermissionMiddleware("savings.view"), handlers.GetAllSavingTypesWithBalances)
				}

				// Accounts & Transactions
				savings.GET("/accounts", middleware.PermissionMiddleware("savings.view"), handlers.GetAllSavingsAccounts)
				savings.GET("/transactions", middleware.PermissionMiddleware("savings.view"), handlers.GetSavingsTransactions)
				savings.GET("/transactions/all", middleware.PermissionMiddleware("savings.view"), handlers.GetAllTransactionsList)
				savings.POST("/transactions", middleware.PermissionMiddleware("savings.transaction"), handlers.CreateSavingsTransaction)
				savings.PUT("/transactions/:id", middleware.PermissionMiddleware("savings.update"), handlers.UpdateSavingsTransaction)
				savings.DELETE("/transactions/:id", middleware.PermissionMiddleware("savings.delete"), handlers.DeleteSavingsTransaction)
				savings.GET("/member/:memberId/balance", middleware.PermissionMiddleware("savings.view"), handlers.GetMemberSavingsBalance)
			}

			// ---------------------------------------------------------
			// 4. PEMBIAYAAN & MARGIN
			// ---------------------------------------------------------
			margin := protected.Group("/margin-setups")
			{
				margin.GET("", middleware.PermissionMiddleware("config.view"), handlers.GetMargins)
				margin.POST("", middleware.PermissionMiddleware("config.manage"), handlers.CreateMargin)
				margin.PUT("/:id", middleware.PermissionMiddleware("config.manage"), handlers.UpdateMargin)
				margin.DELETE("/:id", middleware.PermissionMiddleware("config.manage"), handlers.DeleteMargin)
			}

			pembiayaan := protected.Group("/pembiayaan")
			{
				pembiayaan.GET("", middleware.PermissionMiddleware("pembiayaan.view"), handlers.GetPembiayaan)
				pembiayaan.POST("", middleware.PermissionMiddleware("pembiayaan.create"), handlers.CreatePembiayaan)
				pembiayaan.PUT("/:id", middleware.PermissionMiddleware("pembiayaan.update"), handlers.UpdatePembiayaan)
				pembiayaan.DELETE("/:id", middleware.PermissionMiddleware("pembiayaan.delete"), handlers.DeletePembiayaan)
				pembiayaan.GET("/:id/pembayaran", middleware.PermissionMiddleware("pembiayaan.view"), handlers.GetPembayaranByPinjamanID)
				pembiayaan.GET("/margin", middleware.PermissionMiddleware("pembiayaan.view"), handlers.GetMarginByCategoryAndTenor)
			}

			pembayaran := protected.Group("/pembayaran-pembiayaan")
			{
				pembayaran.POST("", middleware.PermissionMiddleware("pembiayaan.pay"), handlers.CreatePembayaran)
				pembayaran.DELETE("/:id", middleware.PermissionMiddleware("pembiayaan.delete"), handlers.DeletePembayaran)
			}

			// ---------------------------------------------------------
			// 5. QARD HASSAN
			// ---------------------------------------------------------
			qardhassan := protected.Group("/qardhassan")
			{
				qardhassan.GET("", middleware.PermissionMiddleware("qardhassan.view"), handlers.GetQardHassan)
				qardhassan.GET("/:id", middleware.PermissionMiddleware("qardhassan.view"), handlers.GetQardHassanByID)
				qardhassan.POST("", middleware.PermissionMiddleware("qardhassan.create"), handlers.CreateQardHassan)
				qardhassan.PUT("/:id", middleware.PermissionMiddleware("qardhassan.update"), handlers.UpdateQardHassan)
				qardhassan.DELETE("/:id", middleware.PermissionMiddleware("qardhassan.delete"), handlers.DeleteQardHassan)
				qardhassan.POST("/:id/pay", middleware.PermissionMiddleware("qardhassan.pay"), handlers.PayQardHassan)
			}

			// ---------------------------------------------------------
			// 6. MASTER DATA (REKENING, KATEGORI)
			// ---------------------------------------------------------
			rekening := protected.Group("/rekening")
			{
				rekening.GET("", middleware.PermissionMiddleware("rekening.view"), handlers.GetRekenings)
				rekening.POST("", middleware.PermissionMiddleware("rekening.manage"), handlers.CreateRekening)
				rekening.GET("/:id/mutasi", middleware.PermissionMiddleware("rekening.view"), handlers.GetMutasiRekening)
			}

			kategoriBarang := protected.Group("/kategori-barangs")
			{
				kategoriBarang.GET("", middleware.PermissionMiddleware("config.view"), handlers.GetKategoriBarangs)
				kategoriBarang.POST("", middleware.PermissionMiddleware("config.manage"), handlers.CreateKategoriBarang)
			}
		}
	}
}