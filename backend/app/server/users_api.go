package server

import (
	"codejam.io/database"
	"codejam.io/server/models"
	"github.com/emicklei/pgtalk/convert"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func validateDisplayName(displayName string, response *models.FormResponse) {
	if displayName == "" {
		response.AddError("DisplayName", "required")
	}
}

func (server *Server) GetUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	if userId != nil {
		dbUser, err := database.GetUser(convert.StringToUUID(userId.(string)))
		if err != nil {
			logger.Error("GetUser: %v, %v", userId, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, dbUser)
		return
	}
	ctx.Status(http.StatusUnauthorized)
}

type PutProfileRequest struct {
	DisplayName string
}

func (server *Server) PutProfile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	if userId != nil {
		user, err := database.GetUser(convert.StringToUUID(userId.(string)))
		if err != nil {
			logger.Error("PutProfile: GetUser error: %v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		// admin display name lock prevents a user from changing it
		if user.LockDisplayName {
			logger.Info("DisplayName Locked")
			ctx.Status(http.StatusForbidden)
			return
		}

		var request PutProfileRequest
		err = ctx.ShouldBindJSON(&request)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		response := models.NewFormResponse()
		request.DisplayName = strings.Trim(request.DisplayName, " ")

		// Perform validation
		validateDisplayName(request.DisplayName, &response)
		if len(response.Errors) > 0 {
			logger.Error("Validation Error: %v+", user)
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		user.DisplayName = request.DisplayName
		_, err = database.UpdateUser(user)
		if err != nil {
			logger.Error("Error calling database.UpdateEvent: %v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		} else {
			response.Data = user
			ctx.JSON(http.StatusOK, response)
		}

	} else {
		logger.Error("PutUser Unauthorized: no session")
		ctx.Status(http.StatusUnauthorized)
	}
}

// Logout is a GET route for logging out a user.
// This involved clearing the session cookie, clearing/deleting the entry from the session store
func (server *Server) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	err := session.Save()
	if err != nil {
		logger.Error("Logout: error saving session: %v", err)
	}

	// Have to manually do this, not sure why the session middleware doesn't handle this...
	ctx.SetCookie(SessionCookieName, "", -1, "/", "", false, false)
	ctx.Redirect(http.StatusFound, ctx.Request.Header.Get("Referer"))
}

func (server *Server) SetupUserRoutes() {
	logger.Info("Setting up User routes...")

	group := server.Gin.Group("/user")
	{
		group.GET("/", server.GetUser)
		group.PUT("/profile", server.PutProfile)
		group.GET("/logout", server.Logout)
	}

}
