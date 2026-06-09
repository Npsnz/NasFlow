package main

import (
	"log"
	"net/http"

	"backend/config"
	"backend/database"
	"backend/handlers"
	"backend/middleware"
	"backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load configuration
	config.LoadConfig()

	// 2. Initialize database
	database.InitDB()

	// 3. Initialize SSE Broker
	handlers.InitSSEBroker()

	// 4. Connect Services with Handlers via callback to avoid circular imports
	services.BroadcastCallback = func(userID uint, event string, payload interface{}) {
		handlers.BroadcastEvent(userID, event, payload)
	}

	// 5. Start background overdue tasks checker
	services.StartOverdueCheckService(database.DB)

	// 6. Setup Router
	r := gin.Default()

	// 7. Global Middleware
	r.Use(middleware.CORSMiddleware())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Public Auth endpoints
	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)
	r.POST("/api/auth/refresh", handlers.RefreshToken)

	// Protected endpoints (require JWT auth)
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Profile
		api.GET("/auth/me", handlers.GetMe)
		api.POST("/auth/logout", handlers.Logout)
		api.PUT("/auth/profile", handlers.UpdateProfile)
		api.DELETE("/auth/delete", handlers.DeleteAccount)

		// Workspaces
		api.GET("/workspaces", handlers.ListWorkspaces)
		api.POST("/workspaces", handlers.CreateWorkspace)
		api.PUT("/workspaces/:id", handlers.UpdateWorkspace)
		api.DELETE("/workspaces/:id", handlers.DeleteWorkspace)
		api.PUT("/workspaces/reorder", handlers.ReorderWorkspaces)

		// Tasks
		api.GET("/tasks", handlers.ListTasks)
		api.POST("/tasks", handlers.CreateTask)
		api.GET("/tasks/:id", handlers.GetTask)
		api.PUT("/tasks/:id", handlers.UpdateTask)
		api.DELETE("/tasks/:id", handlers.DeleteTask)
		api.PUT("/tasks/:id/status", handlers.UpdateTaskStatus)
		api.PUT("/tasks/reorder", handlers.ReorderTasks)
		api.POST("/tasks/:id/complete", handlers.CompleteTask)
		api.GET("/tasks/overdue", handlers.GetOverdue)

		// Tags
		api.GET("/tags", handlers.ListTags)
		api.POST("/tags", handlers.CreateTag)
		api.PUT("/tags/:id", handlers.UpdateTag)
		api.DELETE("/tags/:id", handlers.DeleteTag)

		// Comments
		api.POST("/tasks/:id/comments", handlers.AddComment)
		api.DELETE("/comments/:id", handlers.DeleteComment)

		// Statistics
		api.GET("/stats", handlers.GetStats)

		// SSE Real-time Feed
		api.GET("/sse", handlers.HandleSSE)
	}

	// 8. Start Web Server
	port := config.AppConfig.Port
	log.Printf("TaskFlow backend running on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
