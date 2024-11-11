package server

import (
	"codejam.io/database"
	"github.com/emicklei/pgtalk/convert"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (server *Server) GetAllUsers(ctx *gin.Context) {
	if server.VerifyAdminAccess(ctx) {
		users, err := database.GetAllUsers()
		if err != nil {
			logger.Error("Error getting all users: %v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, users)
	}
}

type PutAccountStatusRequest struct {
	AccountStatus string
}

func (server *Server) PutAccountStatus(ctx *gin.Context) {
	userIdParam := ctx.Param("id")
	var request PutAccountStatusRequest
	if server.VerifyAdminAccess(ctx) &&
		server.VerifyUserNotAdmin(ctx, userIdParam) &&
		server.DeserializeRequest(ctx, &request) {
		user, err := database.SetAccountStatus(convert.StringToUUID(userIdParam), request.AccountStatus)
		if err != nil {
			logger.Error("PutAccountStatus: SetAccountStatus error: %v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}

type PutDisplayNameRequest struct {
	DisplayName string
}

func (server *Server) PutDisplayName(ctx *gin.Context) {
	userIdParam := ctx.Param("id")
	var request PutDisplayNameRequest
	if server.VerifyAdminAccess(ctx) &&
		server.VerifyUserNotAdmin(ctx, userIdParam) &&
		server.DeserializeRequest(ctx, &request) {
		user, err := database.SetDisplayName(convert.StringToUUID(userIdParam), request.DisplayName)
		if err != nil {
			logger.Error("SetDisplayName error: %v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}

type PutDisplayNameLockRequest struct {
	Lock bool
}

func (server *Server) PutDisplayNameLock(ctx *gin.Context) {
	userIdParam := ctx.Param("id")
	var request PutDisplayNameLockRequest
	if server.VerifyAdminAccess(ctx) &&
		server.VerifyUserNotAdmin(ctx, userIdParam) &&
		server.DeserializeRequest(ctx, &request) {
		user, err := database.SetDisplayNameLock(convert.StringToUUID(userIdParam), request.Lock)
		if err != nil {
			logger.Error("SetDisplayNameLock error: %v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}

func (server *Server) PutBan(ctx *gin.Context) {
	userIdParam := ctx.Param("id")
	if server.VerifyAdminAccess(ctx) &&
		server.VerifyUserNotAdmin(ctx, userIdParam) {
		user, err := database.SetAccountStatus(convert.StringToUUID(userIdParam), "BANNED")
		if err != nil {
			logger.Error("PutBan: SetAccountStatus error: %v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}

func (server *Server) PutUnban(ctx *gin.Context) {
	userIdParam := ctx.Param("id")
	if server.VerifyAdminAccess(ctx) &&
		server.VerifyUserNotAdmin(ctx, userIdParam) {
		user, err := database.SetAccountStatus(convert.StringToUUID(userIdParam), "ACTIVE")
		if err != nil {
			logger.Error("PutUnban: SetAccountStatus error: %v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}

func (server *Server) SetupAdminUserRoutes() {
	logger.Info("Setting up Admin User routes...")

	group := server.Gin.Group("/admin/user")
	{
		group.GET("/all", server.GetAllUsers)
		group.PUT("/:id/display_name/", server.PutDisplayName)
		group.PUT("/:id/display_name_lock/", server.PutDisplayNameLock)
		group.PUT("/:id/ban/", server.PutBan)
		group.PUT("/:id/unban/", server.PutUnban)
	}

}
