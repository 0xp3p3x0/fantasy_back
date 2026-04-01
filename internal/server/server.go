package server

import (
	"back/internal/config"
	"back/internal/handler"
	"back/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	engine   *gin.Engine
	handlers *Handlers
	db       *gorm.DB
}

func NewServer(db *gorm.DB, cfg *config.Config) (*Server, error) {
	config.SetDB(db)

	authService := service.NewAuthService(db, cfg)
	profileService := service.NewProfileService(db, cfg)
	casinoService := service.NewCasinoService(db, cfg)
	agentService := service.NewAgentService(db)
	// previous Handlers struct is kept for compatibility but not used
	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler(profileService)
	casinoHandler := handler.NewCasinoHandler(casinoService)
	agentHandler := handler.NewAgentHandler(agentService, authService)
	handlers := &Handlers{
		Auth:    authHandler,
		Profile: profileHandler,
		Casino:  casinoHandler,
		Agent:   agentHandler,
	}

	server := &Server{
		engine:   gin.Default(),
		handlers: handlers,
		db:       db,
	}

	server.engine = SetupRouter(authService, profileService, casinoService, agentService, cfg.SecretKey)
	return server, nil
}

func (s *Server) Run(port string) error {
	return s.engine.Run(":" + port)
}
