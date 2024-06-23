package server

import (
	"codejam.io/database"
	"codejam.io/server/models"
	"github.com/emicklei/pgtalk/convert"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/mrz1836/go-sanitize"
	"net/http"
	"strings"
)

func sanitizeEvent(event *database.DBEvent) {
	event.Title = sanitize.Scripts(event.Title)
	event.Description = sanitize.Scripts(event.Description)
	event.Rules = sanitize.Scripts(event.Rules)
}

func validateEvent(event database.DBEvent, response *models.FormResponse) {
	// Title is required
	if strings.Trim(event.Title, " ") == "" {
		response.AddError("Title", "required")
	}
}

func (server *Server) GetAllEvents(ctx *gin.Context) {
	events, err := database.GetEvents()
	if err == nil {
		ctx.JSON(http.StatusOK, events)
	} else {
		ctx.Status(http.StatusInternalServerError)
	}
}

func (server *Server) GetEvent(ctx *gin.Context) {
	id := ctx.Param("id")
	event, err := database.GetEvent(convert.StringToUUID(id))
	if err == nil {
		ctx.JSON(http.StatusOK, event)
	} else {
		logger.Error("GetEvent error: %v", err)
		ctx.Status(http.StatusInternalServerError)
	}
}

func (server *Server) GetActiveEvent(ctx *gin.Context) {
	event, err := database.GetActiveEvent()
	if err == nil {
		ctx.JSON(http.StatusOK, event)
	} else if err == pgx.ErrNoRows {
		ctx.Status(http.StatusNoContent)
	} else {
		logger.Error("GetActiveEvent error: %v", err)
		ctx.Status(http.StatusInternalServerError)
	}
}

func (server *Server) PostEvent(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	if userId != nil {
		event, err := database.CreateEvent(convert.StringToUUID(userId.(string)))
		if err == nil {
			ctx.JSON(http.StatusOK, event)
		} else {
			logger.Error("CreateEvent error: %v for user %s", err, userId)
			ctx.Status(http.StatusInternalServerError)
		}
	} else {
		ctx.Status(http.StatusUnauthorized)
	}
}

func (server *Server) PutEvent(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	if userId != nil {
		var event database.DBEvent
		err := ctx.ShouldBindJSON(&event)
		if err != nil {
			logger.Error("UpdateEvent Request ShouldBindJSON error: %v", err)
			ctx.Status(http.StatusBadRequest)
			return
		}

		// Check Authorization
		user := database.GetUser(convert.StringToUUID(userId.(string)))
		if user.Role != database.Admin {
			logger.Error("PutEvent unauthorized user: %v", userId)
			ctx.Status(http.StatusUnauthorized)
			return
		}

		response := models.NewFormResponse()

		sanitizeEvent(&event)
		validateEvent(event, &response)

		// Perform validation
		if len(response.Errors) > 0 {
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		// Update the event in the DB
		event, err = database.UpdateEvent(event)
		if err != nil {
			logger.Error("Error calling database.UpdateEvent: %v", err)
			ctx.Status(http.StatusInternalServerError)
		} else {
			logger.Info("User %v updated Event %v", userId, event.Id)
			response.Data = event
			ctx.JSON(http.StatusOK, response)
		}
	} else {
		logger.Error("PutEvent Unauthorized: no session")
		ctx.Status(http.StatusUnauthorized)
	}
}

func (server *Server) GetStatuses(ctx *gin.Context) {
	statuses, err := database.GetStatuses()
	if err != nil {
		logger.Error("Error in GetStatuses: %v", err)
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.JSON(http.StatusOK, statuses)
	}
}

func (server *Server) SetupEventRoutes() {
	group := server.Gin.Group("/event")
	{
		group.GET("/", server.GetAllEvents)
		group.GET("/active", server.GetActiveEvent)
		group.GET("/:id", server.GetEvent)
		group.PUT("/:id", server.PutEvent)
		group.POST("/", server.PostEvent)
		group.GET("/statuses", server.GetStatuses)
	}
}
