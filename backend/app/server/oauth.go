package server

import (
	"net/http"
	"os"
	"strings"

	"codejam.io/database"
	"codejam.io/integrations"
	"github.com/emicklei/pgtalk/convert"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"
)

// SetupOAuth initializes the OAuth provider specified in the application config.
func (server *Server) SetupOAuth() {
	if server.Debug {
		logger.Warn("Debug mode is set. No OAuth Providers are set!")
		return
	}
	
	var endpoint oauth2.Endpoint

	switch strings.ToLower(server.Config.OAuth.Provider) {
	case "github":
		endpoint = githubOAuth.Endpoint
	case "discord":
		endpoint = oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		}
	default:
		logger.Critical("Invalid OAuth provider: %s", server.Config.OAuth.Provider)
		os.Exit(1)
	}

	server.OAuth = &oauth2.Config{
		ClientID:     server.Config.OAuth.Id,
		ClientSecret: server.Config.OAuth.Secret,
		Endpoint:     endpoint,
		RedirectURL:  server.Config.OAuth.RedirectUrl,
		Scopes:       server.Config.OAuth.Scopes,
	}
}

func (server *Server) GetOAuthRedirect(ctx *gin.Context) {
	if server.Debug {
		ctx.Redirect(http.StatusFound, "/oauth/debug-login")
		return
	}

	url := server.OAuth.AuthCodeURL(ctx.Request.Header.Get("Referer"))
	ctx.Redirect(http.StatusFound, url)
}

func (server *Server) GetDebugSession(ctx *gin.Context) {
	redir := ctx.Query("state")

	dbUser := database.CreateUser("discord", "0", "DebugCow", "")
	session := sessions.Default(ctx)
	session.Set("userId", convert.UUIDToString(dbUser.Id))
	session.Set("displayName", dbUser.DisplayName)

	err := session.Save()

	if err != nil {
		logger.Error("Error saving debug session: %v", err)
	}

	ctx.Redirect(http.StatusFound, redir)
}

func (server *Server) GetOAuthCallback(ctx *gin.Context) {
	authCode := ctx.Query("code")
	redir := ctx.Query("state")
	token, err := server.OAuth.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		// todo - can any of these be handled?
		logger.Error("OAuth exchange error: %v", err)
		return
	}

	integrationName := strings.ToLower(server.Config.OAuth.Provider)
	providerUser := integrations.GetUser(integrationName, token.AccessToken)
	if providerUser != nil {
		dbUser := database.CreateUser(integrationName, providerUser.UserId, providerUser.DisplayName, providerUser.AvatarUrl)
		session := sessions.Default(ctx)
		session.Set("userId", convert.UUIDToString(dbUser.Id))
		session.Set("displayName", dbUser.DisplayName)
		err = session.Save()
		if err != nil {
			logger.Error("Error saving session: %v", err)
		}

		ctx.Redirect(http.StatusFound, redir)
	} else {
		// TODO error page
		logger.Error("Unable to lookup provider user for %s", integrationName)
	}
}

func (server *Server) SetupOAuthRoutes() {
	logger.Info("Setting up OAuth routes...")

	group := server.Gin.Group("/oauth")
	{
		group.GET("/redirect", server.GetOAuthRedirect)
		group.GET("/callback", server.GetOAuthCallback)
		group.GET("/debug-login", server.GetDebugSession)
	}
}
