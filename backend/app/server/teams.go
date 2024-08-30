package server

import (
	"fmt"
	"net/http"

	"codejam.io/database"
	"github.com/emicklei/pgtalk/convert"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math"
	"math/big"
)

type CreateTeamRequest struct {
	// sent from UI, to be processed in the server into a DBTeam structure
	// referenced in server/teams.go/CreateTeam line 53
	EventId      string
	Name         string
	Visibility   string
	Availability string
	Description  string
	Technologies string
	Timezone     string
}

type GetTeamResponse struct {
	Team    *database.DBTeam
	Event   *database.DBEvent
	Members *[]database.DBTeamMemberInfo // array(slice) of a struct
}

type InvitePayload struct {
	TeamId 		string	`json:"teamId"`
	InviteCode	string 	`json:"inviteCode"`
}

type JoinPayload struct {
	TeamId		string 	`json:"teamId"`
}

func MD5HashCode(teamName string) (string, error) {
	randNum, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return "Md5 hash error", err
	}
	hash := md5.Sum([]byte(teamName + "." + randNum.String()))
	return hex.EncodeToString(hash[:7]), nil
}

func (server *Server) signupsAllowed(eventId string) bool {
	statusCode, err := database.GetEventStatusCode(convert.StringToUUID(eventId))
	if err != nil {
		logger.Error("Error in GetEventStatusCode: %v+", err)
		return false
	}

	if statusCode == "SIGNUP" || statusCode == "STARTED" {
		return true
	} else {
		return false
	}
}

func (server *Server) GetAllTeams(ctx *gin.Context) {
	// what if there's no session => no user id
	teams, err := database.GetTeams()
	if err == nil {
		ctx.JSON(http.StatusOK, teams)
	} else {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (server *Server) GetUserTeams(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	if userId == nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	strUserId := userId.(string)

	teams, err := database.GetUserTeams(convert.StringToUUID(strUserId))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, teams)
}

// stepp 4: GET team info
// purpose is to construct the DBTeamMemberInfo
func (server *Server) GetTeamInfo(id pgtype.UUID) (*GetTeamResponse, error) {

	var teamResponse GetTeamResponse
	var team database.DBTeam
	var event database.DBEvent
	var members *[]database.DBTeamMemberInfo //user info based on teamId
	// fmt.Println("========id: ", id) prints: {[204 69 126 62 33 10 77 93 131 216 8 153 66 109 252 147] true}
	team, err := database.GetTeam(id)
	if err != nil {
		logger.Error("failed to get team: %v", err)
		return nil, err
	}

	event, err = database.GetEvent(team.EventId)
	if err != nil {
		logger.Error("failed to get event: %v", err)
		return nil, err
	}

	members, err = database.GetMembersByTeamId(team.Id)
	if err != nil {
		logger.Error("failed to get event: %v", err)
		return nil, err
	}

	// attach all 3 structures to GetTeamResponse --> nested structs turn into nested JSON (with ctx.JSON)
	teamResponse.Team = &team
	teamResponse.Event = &event
	teamResponse.Members = members
	
	return &teamResponse, nil
}

func (server *Server) sendTeamInfo(ctx *gin.Context) {
	id := convert.StringToUUID(ctx.Param("id"))

	teamResponse, err := server.GetTeamInfo(id)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
	}

	ctx.JSON(http.StatusOK, teamResponse)
}

func (server *Server) GetTeamInfoByInviteCode(ctx *gin.Context) {
	inviteCode := ctx.Param("invitecode")
	fmt.Println("\n===server getteam by invite code: ", inviteCode)

	var teamResponse GetTeamResponse
	var team database.DBTeam
	var event database.DBEvent
	var members *[]database.DBTeamMemberInfo //user info based on teamId

	team, err := database.GetTeamByInvite(inviteCode)
	if err != nil {
		logger.Error("failed to get team: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get team: %v", err)})
		return
	} else {
		fmt.Println("TEST if GetTeamByInvite was SUCCESS: ", err, " + ", team)
	}

	event, err = database.GetEvent(team.EventId)
	if err != nil {
		logger.Error("failed to get event: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get event: %v", err)})
		return
	}

	members, err = database.GetMembersByTeamId(team.Id)
	if err != nil {
		logger.Error("failed to get event: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get members: %v", err)})
		return
	}

	// attach all 3 structures to GetTeamResponse --> nested structs turn into nested JSON (with ctx.JSON)
	teamResponse.Team = &team
	teamResponse.Event = &event
	teamResponse.Members = members

	fmt.Println(teamResponse)
	ctx.JSON(http.StatusOK, teamResponse)
}

func (server *Server) CreateTeam(ctx *gin.Context) {
	// ctx of *gin.Context has HTTP request info.
	// Step 4: Post Team Data API (TWO PARTS 1) create team 2) add team members)
	session := sessions.Default(ctx)
	userId := session.Get("userId")
    if userId == nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	strUserId := userId.(string)

	var team database.DBTeam
	var teamReq CreateTeamRequest

	err := ctx.ShouldBindJSON(&teamReq)
	if err != nil {
		logger.Error("CreateTeam Request ShouldBindJSON error: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	if !server.signupsAllowed(teamReq.EventId) {
		ctx.Status(http.StatusForbidden)
		return
	}

	// CONVERT teamReq to team
	team.EventId = convert.StringToUUID(teamReq.EventId)
	team.Name = teamReq.Name
	team.Availability = teamReq.Availability
	team.Description = teamReq.Description
	team.Visibility = teamReq.Visibility
	team.Technologies = teamReq.Technologies
	team.Timezone = teamReq.Timezone

	md5code, err := MD5HashCode(team.Name)
	if err != nil {
		logger.Error("Error - Md5HashCode failed", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	team.InviteCode = md5code

	fmt.Printf("%+v", team)
	
	// INSERTS TEAM into DB
	// PART 1/2 DONE
	teamUUID, err := database.CreateTeam(team)
	if err != nil {
		logger.Error("Error trying to CreateTeam(team)")
		ctx.Status(http.StatusBadRequest)
		return
	}
	fmt.Println(convert.StringToUUID(strUserId), teamUUID,)

	// PART 2/2 DONE
	// construct TeamMember
	_, err = database.AddTeamMember(convert.StringToUUID(strUserId), teamUUID, "owner")

	if err == nil {
		fmt.Println("Successfully added team member")
		ctx.JSON(http.StatusCreated, map[string]pgtype.UUID{
			"id": teamUUID,
		})
	} else {
		fmt.Println(err)
		logger.Error("AddTeamMember error: %v for user %s", err, strUserId)
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (server *Server) UpdateTeam(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	fmt.Println(ctx.Request.Body)
	if userId != nil {
		var team database.DBTeam
		err := ctx.ShouldBindJSON(&team)  //"message incoming data to this struct"
		if err != nil {
			logger.Error("UpdateEvent Request ShouldBindJSON error: %v", err)
			ctx.Status(http.StatusBadRequest)
			return
		}
		team, err = database.UpdateTeam(team)
		if err != nil {
			logger.Error("Error calling database.UpdateEvent: %v", err)
			ctx.Status(http.StatusInternalServerError)
		} else {
			ctx.JSON(http.StatusOK, team)
		}
	} else {
		ctx.Status(http.StatusUnauthorized)
	}
}

func (server *Server) UpdateTeamMembers(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	if userId != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	//database.AddTeamMember(convert.StringToUUID(userId.(string), ))
}

func (server *Server) MemberJoin(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	fmt.Println("userID = ", userId)

	if userId == nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	strUserId := userId.(string)
	uuidUserId := convert.StringToUUID(strUserId)

	var teamId JoinPayload

	if err := ctx.ShouldBindJSON(&teamId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// teamId prints: {cc457e3e-210a-4d5d-83d8-0899426dfc93}
	uuidTeamId := convert.StringToUUID(teamId.TeamId)
	teamInfo, err := server.GetTeamInfo(uuidTeamId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}
	// TODO: check if the user already exists in the team before adding the member. 
	// use teamInfo: 
	// type GetTeamResponse struct {
	// 	Team    *database.DBTeam
	// 	Event   *database.DBEvent
	// 	Members *[]database.DBTeamMemberInfo // array(slice) of a struct
	// }

	fmt.Println("print (*teamInfo.Members)[0].DBUser.Id = ", convert.UUIDToString((*teamInfo.Members)[0].DBUser.Id)) 
	// ^ dereferences pointer to get to actual slice of database.DBTeamMemberInfo, 
	// then accesses the first element in dereferenced slice


	// verifies team is public
	var isPublic string = teamInfo.Team.Visibility
	if isPublic == "private" {
		ctx.Status(http.StatusForbidden)
		return 
	}
	_, err = database.AddTeamMember(uuidUserId, uuidTeamId, "member")
	if err != nil {
		ctx.JSON(http.StatusConflict, err)
		return
	}
	strTeamId := convert.UUIDToString(teamInfo.Team.Id)
	ctx.JSON(http.StatusOK, strTeamId)
}

func (server *Server) MemberInvite(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId:= session.Get("userId")

	var payload InvitePayload
	// bind the message context to the structure, and do an error check
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inviteCode := payload.InviteCode
	teamId := payload.TeamId

	fmt.Println("===316", inviteCode, teamId, userId)

	//var team database.DBTeamMemberInfo
	// to add a member, i need to add a user to team_member table
}


func (server *Server) SetupTeamRoutes() {
	group := server.Gin.Group("/team")
	{
		group.POST("/", server.CreateTeam)
		group.POST("/edit/:teamid", server.UpdateTeamMembers) // for admin to remove people
		group.POST("/join", server.MemberJoin)
		group.POST("/:invitecode", server.MemberInvite)
		//group.GET("/", server.GetAllTeams)
		group.GET("/:id", server.sendTeamInfo)
		group.GET("/invite/:invitecode", server.GetTeamInfoByInviteCode)
		// group.PUT("/:id", server.UpdateTeam)
		// Step 3: Post Team Data API
	}

	server.Gin.GET("/teams", server.GetUserTeams) // I think this works rofl
	server.Gin.GET("/teams/browse", server.GetAllTeams)
}
