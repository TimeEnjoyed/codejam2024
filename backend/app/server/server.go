package server

import (
	"codejam.io/config"
	"codejam.io/logging"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"os"
)

var logger = logging.NewLogger(logging.Options{Name: "Server", Level: logging.DEBUG})

var SessionCookieName string = "session"

type Server struct {
	Config config.Config
	OAuth  *oauth2.Config
	Gin    *gin.Engine
}

func (server *Server) SetupSessionStore() {
	logger.Info("Setting up Session Store")

	store, err := redis.NewStore(
		server.Config.Redis.Size,
		server.Config.Redis.Protocol,
		server.Config.Redis.Address,
		server.Config.Redis.Password,
		[]byte(server.Config.Redis.CookieStoreSecret))
	if err != nil {
		logger.Critical("Error initializing Redis session store: %v", err)
		os.Exit(1)
	}

	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 90,
		Secure:   true,
		HttpOnly: true,
	})
	server.Gin.Use(sessions.Sessions(SessionCookieName, store))
}

func (server *Server) StartServer() {
	logger.Info("Starting server...")

	server.Gin = gin.Default()

	server.SetupSessionStore()
	server.SetupOAuth()

	// Setup routes...
	server.SetupOAuthRoutes()
	server.SetupUserRoutes()
	server.SetupEventRoutes()
	server.SetupStaticRoutes()

	// Start the server...
	logger.Info("Server Started")
	server.Gin.Run(server.Config.Server.Listen)
}
