package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jaydeep/go-n8n/configs"
	"github.com/jaydeep/go-n8n/internal/interfaces/http/middleware"
	"github.com/jaydeep/go-n8n/pkg/database"
	"github.com/jaydeep/go-n8n/pkg/logger"
)

// NewRouter creates and configures the main router
func NewRouter(cfg *configs.Config, db *database.DB, log *logger.Logger) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(log))
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS(cfg.CORS))
	
	// Rate limiting
	if cfg.RateLimit.Enabled {
		router.Use(middleware.RateLimit(cfg.RateLimit))
	}

	// Health check endpoints
	router.GET("/health", healthCheck)
	router.GET("/ready", readinessCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", registerHandler)
			auth.POST("/login", loginHandler)
			auth.POST("/refresh", refreshTokenHandler)
			auth.POST("/forgot-password", forgotPasswordHandler)
			auth.POST("/reset-password", resetPasswordHandler)
			auth.POST("/verify-email", verifyEmailHandler)
		}

		// Webhook endpoints (public but validated)
		v1.Any("/webhook/:path", webhookHandler)

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.Auth(cfg.JWT))
		{
			// User routes
			protected.GET("/auth/me", getCurrentUser)
			protected.PUT("/auth/me", updateCurrentUser)
			protected.POST("/auth/logout", logoutHandler)
			protected.POST("/auth/change-password", changePasswordHandler)
			protected.POST("/auth/2fa/enable", enable2FAHandler)
			protected.POST("/auth/2fa/disable", disable2FAHandler)
			protected.POST("/auth/2fa/verify", verify2FAHandler)

			// Workflow routes
			workflows := protected.Group("/workflows")
			{
				workflows.GET("", listWorkflows)
				workflows.POST("", createWorkflow)
				workflows.GET("/:id", getWorkflow)
				workflows.PUT("/:id", updateWorkflow)
				workflows.DELETE("/:id", deleteWorkflow)
				workflows.POST("/:id/activate", activateWorkflow)
				workflows.POST("/:id/deactivate", deactivateWorkflow)
				workflows.POST("/:id/execute", executeWorkflow)
				workflows.POST("/:id/duplicate", duplicateWorkflow)
				workflows.GET("/:id/executions", getWorkflowExecutions)
				workflows.POST("/:id/share", shareWorkflow)
				workflows.GET("/:id/versions", getWorkflowVersions)
				workflows.POST("/:id/test", testWorkflow)
				workflows.GET("/:id/nodes", getWorkflowNodes)
				workflows.PUT("/:id/nodes", updateWorkflowNodes)
				workflows.GET("/:id/export", exportWorkflow)
				workflows.POST("/import", importWorkflow)
				workflows.GET("/:id/statistics", getWorkflowStatistics)
				workflows.GET("/:id/metrics", getWorkflowMetrics)
				workflows.POST("/:id/versions/:versionId/restore", restoreWorkflowVersion)
				workflows.POST("/batch", batchWorkflowOperations)
			}

			// Node routes
			nodes := protected.Group("/nodes")
			{
				nodes.GET("/types", listNodeTypes)
				nodes.GET("/types/:type", getNodeType)
				nodes.GET("/types/:type/schema", getNodeSchema)
				nodes.POST("/test", testNode)
				nodes.PUT("/:id", updateNode)
				nodes.DELETE("/:id", deleteNode)
				nodes.POST("/:id/test", testNodeById)
				nodes.GET("/:id/executions/:executionId/data", getNodeExecutionData)
				nodes.POST("/:id/pin", pinNodeData)
				nodes.DELETE("/:id/pin", unpinNodeData)
			}

			// Execution routes
			executions := protected.Group("/executions")
			{
				executions.GET("", listExecutions)
				executions.GET("/:id", getExecution)
				executions.POST("/:id/stop", stopExecution)
				executions.POST("/:id/retry", retryExecution)
				executions.DELETE("/:id", deleteExecution)
				executions.GET("/:id/data", getExecutionData)
				executions.POST("/delete", deleteMultipleExecutions)
				executions.GET("/:id/logs", getExecutionLogs)
				executions.GET("/:id/timeline", getExecutionTimeline)
			}

			// Credential routes
			credentials := protected.Group("/credentials")
			{
				credentials.GET("", listCredentials)
				credentials.POST("", createCredential)
				credentials.GET("/:id", getCredential)
				credentials.PUT("/:id", updateCredential)
				credentials.DELETE("/:id", deleteCredential)
				credentials.POST("/:id/test", testCredential)
				credentials.GET("/oauth2/:credentialType/auth", getOAuth2URL)
				credentials.GET("/oauth2/callback", oAuth2Callback)
				credentials.POST("/:id/share", shareCredential)
			}

			// Variable routes
			variables := protected.Group("/variables")
			{
				variables.GET("", listVariables)
				variables.POST("", createVariable)
				variables.GET("/:key", getVariable)
				variables.PUT("/:key", updateVariable)
				variables.DELETE("/:key", deleteVariable)
			}

			// Tag routes
			tags := protected.Group("/tags")
			{
				tags.GET("", listTags)
				tags.POST("", createTag)
				tags.PUT("/:id", updateTag)
				tags.DELETE("/:id", deleteTag)
			}

			// Settings routes
			settings := protected.Group("/settings")
			{
				settings.GET("", getSettings)
				settings.PUT("", updateSettings)
				settings.GET("/smtp", getSMTPSettings)
				settings.PUT("/smtp", updateSMTPSettings)
				settings.POST("/smtp/test", testSMTPSettings)
			}

			// Stats routes
			stats := protected.Group("/stats")
			{
				stats.GET("/workflows", getWorkflowStats)
				stats.GET("/executions", getExecutionStats)
				stats.GET("/usage", getUsageStats)
			}

			// User management routes
			users := protected.Group("/users")
			{
				users.GET("/:id", getUser)
				users.PUT("/:id", updateUser)
				users.PUT("/:id/settings", updateUserSettings)
				users.GET("/:id/permissions", getUserPermissions)
				users.PUT("/:id/permissions", updateUserPermissions)
			}

			// Templates routes
			templates := protected.Group("/templates")
			{
				templates.GET("", listTemplates)
				templates.GET("/:id", getTemplate)
				templates.POST("", createTemplate)
				templates.PUT("/:id", updateTemplate)
				templates.DELETE("/:id", deleteTemplate)
				templates.POST("/:id/use", useTemplate)
				templates.GET("/categories", getTemplateCategories)
			}

			// API Keys routes
			apiKeys := protected.Group("/api-keys")
			{
				apiKeys.GET("", listAPIKeys)
				apiKeys.POST("", createAPIKey)
				apiKeys.GET("/:id", getAPIKey)
				apiKeys.DELETE("/:id", revokeAPIKey)
			}

			// Webhooks routes
			webhooks := protected.Group("/webhooks")
			{
				webhooks.GET("", listWebhooks)
				webhooks.POST("", createWebhook)
				webhooks.GET("/:id", getWebhook)
				webhooks.PUT("/:id", updateWebhook)
				webhooks.DELETE("/:id", deleteWebhook)
				webhooks.POST("/:id/test", testWebhook)
				webhooks.GET("/:id/url", getWebhookURL)
			}

			// Schedules routes
			schedules := protected.Group("/schedules")
			{
				schedules.GET("", listSchedules)
				schedules.POST("", createSchedule)
				schedules.GET("/:id", getSchedule)
				schedules.PUT("/:id", updateSchedule)
				schedules.DELETE("/:id", deleteSchedule)
				schedules.POST("/:id/activate", activateSchedule)
				schedules.POST("/:id/deactivate", deactivateSchedule)
			}

			// Notifications routes
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", getNotifications)
				notifications.PUT("/:id/read", markNotificationRead)
				notifications.PUT("/read-all", markAllNotificationsRead)
				notifications.DELETE("/:id", deleteNotification)
				notifications.GET("/settings", getNotificationSettings)
				notifications.PUT("/settings", updateNotificationSettings)
			}

			// Search routes
			search := protected.Group("/search")
			{
				search.GET("", globalSearch)
				search.GET("/workflows", searchWorkflows)
				search.GET("/executions", searchExecutions)
			}

			// Audit logs routes
			auditLogs := protected.Group("/audit-logs")
			{
				auditLogs.GET("", listAuditLogs)
				auditLogs.GET("/:id", getAuditLog)
			}

			// Metrics routes
			metrics := protected.Group("/metrics")
			{
				metrics.GET("", getMetrics)
				metrics.GET("/queue", getQueueStatus)
				metrics.GET("/executions", getExecutionStatistics)
				metrics.GET("/workers", getWorkerStatus)
				metrics.GET("/performance", getPerformanceMetrics)
			}

			// Import/Export routes
			protected.GET("/export/workflows", exportAllWorkflows)
			protected.GET("/export/credentials", exportAllCredentials)
			protected.GET("/export/all", exportAllData)
			protected.POST("/import", importData)

			// Community routes
			community := protected.Group("/community")
			{
				community.GET("/workflows", getCommunityWorkflows)
				community.POST("/workflows", publishWorkflowToCommunity)
				community.GET("/workflows/:id/reviews", getWorkflowReviews)
				community.POST("/workflows/:id/reviews", addWorkflowReview)
				community.POST("/workflows/:id/report", reportWorkflow)
			}

			// Integrations routes
			integrations := protected.Group("/integrations")
			{
				integrations.GET("", listIntegrations)
				integrations.GET("/:name", getIntegrationDetails)
				integrations.POST("/:name/install", installIntegration)
				integrations.POST("/:name/uninstall", uninstallIntegration)
				integrations.PUT("/:name", updateIntegration)
			}

			// Teams routes
			teams := protected.Group("/teams")
			{
				teams.GET("", listTeams)
				teams.POST("", createTeam)
				teams.GET("/:id", getTeam)
				teams.PUT("/:id", updateTeam)
				teams.DELETE("/:id", deleteTeam)
				teams.POST("/:id/members", addTeamMember)
				teams.DELETE("/:id/members/:userId", removeTeamMember)
				teams.PUT("/:id/members/:userId", updateTeamMemberRole)
			}

			// Billing routes (Enterprise)
			billing := protected.Group("/billing")
			{
				billing.GET("/usage", getUsageStatistics)
				billing.GET("/info", getBillingInfo)
				billing.GET("/invoices", getInvoices)
				billing.GET("/subscription", getSubscription)
				billing.PUT("/subscription", updateSubscription)
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/users", listUsers)
				admin.GET("/users/:id", getUser)
				admin.PUT("/users/:id", updateUser)
				admin.DELETE("/users/:id", deleteUser)
				admin.POST("/users/:id/activate", activateUser)
				admin.POST("/users/:id/deactivate", deactivateUser)
			}
		}
	}

	// WebSocket endpoint
	router.GET("/ws", websocketHandler)

	// Static files (if needed)
	router.Static("/assets", "./assets")

	return router
}

// Placeholder handlers - to be implemented
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "healthy"})
}

func readinessCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ready"})
}

func registerHandler(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func loginHandler(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func refreshTokenHandler(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func forgotPasswordHandler(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func resetPasswordHandler(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func webhookHandler(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getCurrentUser(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func updateCurrentUser(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func logoutHandler(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func changePasswordHandler(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func listWorkflows(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func createWorkflow(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getWorkflow(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func updateWorkflow(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func deleteWorkflow(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func activateWorkflow(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func deactivateWorkflow(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func executeWorkflow(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func duplicateWorkflow(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getWorkflowExecutions(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func shareWorkflow(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getWorkflowVersions(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func listNodeTypes(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getNodeType(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func testNode(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func listExecutions(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getExecution(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func stopExecution(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func retryExecution(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func deleteExecution(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getExecutionData(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func listCredentials(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func createCredential(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getCredential(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func updateCredential(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func deleteCredential(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func testCredential(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func listVariables(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func createVariable(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getVariable(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func updateVariable(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func deleteVariable(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func listTags(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func createTag(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func updateTag(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func deleteTag(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getSettings(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func updateSettings(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getWorkflowStats(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getExecutionStats(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getUsageStats(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func listUsers(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func getUser(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func updateUser(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func deleteUser(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func activateUser(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func deactivateUser(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func websocketHandler(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
