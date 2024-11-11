package server

import (
	"codejam.io/database"
	"github.com/emicklei/pgtalk/convert"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// VerifyAdminAccess checks if the user has admin access by retrieving the session user account
// and verifying if their role is 'ADMIN'.  Appropriate HTTP responses will be set automatically.
// Returns true if the user has admin access, false otherwise.
func (server *Server) VerifyAdminAccess(ctx *gin.Context) bool {
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	if userId == nil {
		ctx.Status(http.StatusUnauthorized)
		return false
	}

	// get the session user account so we can verify they're an admin - TODO, pull from cache
	currentUser, err := database.GetUser(convert.StringToUUID(userId.(string)))
	if err != nil {
		logger.Error("AdminAccessVerified: GetUser for session user error: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return false
	}

	if currentUser.Role != database.Admin {
		logger.Error("AdminAccessVerified: unauthorized user: %v", userId)
		ctx.Status(http.StatusForbidden)
		return false
	}

	return true
}

// UserIsAdmin checks if the user with the given userId has admin access by retrieving the user from the database
// and verifying if their role is 'ADMIN'. Returns true if the user has admin access, false otherwise.
func (server *Server) UserIsAdmin(userId string) (bool, error) {
	user, err := database.GetUser(convert.StringToUUID(userId))
	if err != nil {
		return false, err
	}

	if user.Role == database.Admin {
		return true, nil
	}

	return false, nil
}

// VerifyUserNotAdmin checks if the user identified by the given userId is not an admin.  This is
// for scenarios when admin accounts should not allow operations to take place on them, such as
// moderation actions.
// Appropriate HTTP responses are set automatically.
// Returns true if the user is not an admin, false otherwise.
func (server *Server) VerifyUserNotAdmin(ctx *gin.Context, userId string) bool {
	// Admin accounts cannot be locked, even by other admins
	isAdmin, err := server.UserIsAdmin(userId)
	if err != nil {
		logger.Error("VerifyUserNotAdmin: UserIsAdmin error: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return false
	}

	if isAdmin {
		logger.Error("VerifyUserNotAdmin: UserIsAdmin TRUE")
		ctx.Status(http.StatusForbidden)
		return false
	}

	return true
}

// DeserializeRequest is a helper function to handle deserialization and automatically handling errors.
// Appropriate HTTP responses are set automatically
// Returns true if the request deserialized successfully, false otherwise
func (server *Server) DeserializeRequest(ctx *gin.Context, request any) bool {
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		logger.Error("%T ShouldBindJSON error: %v", request, err)
		ctx.Status(http.StatusBadRequest)
		return false
	}
	return true
}
