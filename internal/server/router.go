package server

import (
	"back/internal/handler"
	"back/internal/middleware"
	"back/internal/service"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handlers contains all the handlers for the application
type Handlers struct {
	Auth    *handler.AuthHandler
	Profile *handler.ProfileHandler
	Casino  *handler.CasinoHandler
	Agent   *handler.AgentHandler
}

func SetupRouter(
	authService *service.AuthService,
	profileService *service.ProfileService,
	casinoService *service.CasinoService,
	agentService *service.AgentService,
	secretKey string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// Add logging middleware
	router.Use(middleware.RequestLogger())

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:4000", "http://localhost:5173", "http://localhost:7070", "https://www.fantasygaming.games", "https://fantasygaming.games"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoints
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Service is running"})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Service is healthy"})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// create handler instance if not injected
	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler(profileService)
	casinoHandler := handler.NewCasinoHandler(casinoService)
	agentHandler := handler.NewAgentHandler(agentService, authService)
	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/register", authHandler.Register)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.JWTAuth(secretKey))
	{
		protected.GET("/profile", profileHandler.GetProfileById)
		protected.GET("/profile/code/:code", profileHandler.GetProfileByCode)
		protected.PUT("/profile", profileHandler.UpdateProfile)
		protected.PUT("/profile/callback-url", profileHandler.UpdateCallbackURL)
		protected.PUT("/profile/change-password", profileHandler.ChangePassword)
		protected.POST("/casino/gameurl", casinoHandler.GetGameURL)
	}

	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.JWTAuth(secretKey), middleware.RequireRole("admin"))
	{
		admin.POST("/agents", agentHandler.CreateAgent)
		admin.GET("/agents", agentHandler.ListAgents)
		admin.PUT("/agents/:id", agentHandler.UpdateAgent)
		admin.DELETE("/agents/:id", agentHandler.DeleteAgent)
	}

	return router
}
