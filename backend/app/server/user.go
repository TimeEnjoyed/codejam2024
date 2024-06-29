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

func validateUser(user *database.DBUser, response *models.FormResponse) {
	// Title is required
	if strings.Trim(user.DisplayName, " ") == "" {
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

func (server *Server) PutUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	if userId != nil {
		var user database.DBUser
		err := ctx.ShouldBindJSON(&user)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		response := models.NewFormResponse()

		// Perform validation
		validateUser(&user, &response)
		if len(response.Errors) > 0 {
			logger.Error("Validation Error: %v+", user)
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		user.Id = convert.StringToUUID(userId.(string))
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
		group.PUT("/", server.PutUser)
		group.GET("/logout", server.Logout)
	}

}
