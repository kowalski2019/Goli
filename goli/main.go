package main

import (
	aux "goli/auxiliary"
	"goli/database"
	"goli/handler"
	"goli/middlewares"
	"goli/queue"
	response_util "goli/utils"
	"goli/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

var port = aux.GetFromConfig("constants.port")
var host = aux.GetFromConfig("constants.host")

func main() {
	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Initialize and start job queue
	jobQueue := queue.GetQueue()
	jobQueue.SetWebSocketHub(wsHub) // Pass hub to queue for broadcasting
	jobQueue.Start()
	defer jobQueue.Stop()

	// Authenticate with GitHub Container Registry if credentials are configured
	if err := response_util.AuthenticateGitHubContainerRegistry(); err != nil {
		log.Printf("Warning: Failed to authenticate with GitHub Container Registry at startup: %v", err)
		log.Println("You can configure GitHub credentials in the Settings page")
	} else {
		log.Println("GitHub Container Registry authentication successful (if configured)")
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		jobQueue.Stop()
		os.Exit(0)
	}()

	// Create Gin router
	r := gin.Default()

	// Add logging middleware
	r.Use(middlewares.RequestLogger())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public API routes (no auth required)
	public := r.Group("/api/v1")
	{
		// Setup endpoints
		public.POST("/setup/verify", handler.VerifySetupPasswordHandler)
		public.GET("/setup/status", handler.GetSetupStatusHandler)

		// Auth endpoints
		public.POST("/auth/login", handler.LoginHandler)
		public.POST("/auth/2fa/verify", handler.Verify2FAHandler)
		public.POST("/auth/logout", handler.LogoutHandler)
	}

	// Protected API routes (auth required)
	api := r.Group("/api/v1")
	api.Use(middlewares.AuthMiddleware())
	{
		// Job management endpoints
		api.GET("/jobs", handler.ListJobsHandler)
		api.POST("/jobs", handler.CreateJobHandler)
		api.GET("/jobs/:id", handler.GetJobHandler)
		api.POST("/jobs/:id/cancel", handler.CancelJobHandler)

		// Pipeline management endpoints
		api.GET("/pipelines", handler.ListPipelinesHandler)
		api.POST("/pipelines", handler.CreatePipelineHandler)
		api.POST("/pipelines/upload", handler.UploadPipelineHandler)
		api.GET("/pipelines/:id", handler.GetPipelineHandler)
		api.PUT("/pipelines/:id", handler.UpdatePipelineHandler)
		api.POST("/pipelines/:id/run", handler.RunPipelineHandler)
		api.DELETE("/pipelines/:id", handler.DeletePipelineHandler)

		// Config management endpoints
		api.GET("/config", handler.GetConfigHandler)
		api.POST("/config", handler.UpdateConfigHandler)

		// User management endpoints
		api.GET("/users", handler.ListUsersHandler)
		api.POST("/users", handler.CreateUserHandler)
		api.PUT("/users/:id", handler.UpdateUserHandler)
		api.DELETE("/users/:id", handler.DeleteUserHandler)

		// Docker endpoints
		api.POST("/docker/container/start", handler.StartADocker)
		api.POST("/docker/container/stop", handler.StopADocker)
		api.POST("/docker/container/rm", handler.RemoveADocker)
		api.POST("/docker/container/run", handler.RunDockerContainer)
		api.POST("/docker/image/pull", handler.PullAnDockerImage)
		api.POST("/docker/image/rm", handler.RemoveAnDockerImage)
		api.POST("/docker/ps", handler.GetDockerPS)
		api.POST("/docker/images", handler.GetDockerImages)
		api.POST("/docker/compose/up", handler.StartADockerOrchestra)
		api.POST("/docker/compose/down", handler.StopADockerOrchestra)
	}

	// WebSocket endpoint (no auth middleware, but can check auth in handler if needed)
	r.GET("/ws", func(c *gin.Context) {
		handler.ServeWebSocket(wsHub, c)
	})

	// Serve static files (UI) - use NoRoute to handle all non-API routes
	// This allows the frontend SPA to handle client-side routing
	r.NoRoute(func(c *gin.Context) {
		// Check if the request is for a static file (has extension)
		path := c.Request.URL.Path

		// If it's an API or WebSocket route, return 404
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/ws") {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}

		// Try to serve the file, if not found, serve index.html (for SPA routing)
		filePath := "./web" + path
		if _, err := os.Stat(filePath); err == nil && !strings.HasSuffix(path, "/") {
			c.File(filePath)
			return
		}

		// Serve index.html for all other routes (SPA fallback)
		c.File("./web/index.html")
	})

	log.Println("Goli CI/CD Backend is running on " + host)
	log.Fatal(http.ListenAndServe(host+":"+port, r))
}
