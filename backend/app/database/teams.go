package database

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

type DBTeam struct {
	Id           pgtype.UUID      `db:"id"`
	EventId      pgtype.UUID      `db:"event_id"`
	Name         string           `db:"name"`
	Visibility   string           `db:"visibility"`
	Timezone     string           `db:"timezone"`
	Technologies string           `db:"technologies"`
	Availability string           `db:"availability"`
	Description  string           `db:"description"`
	CreatedOn    pgtype.Timestamp `db:"created_on"`
	InviteCode   string           `db:"invite_code"`
}

type CreateTeamMember struct {
	UserId   pgtype.UUID `db:"user_id"`
	TeamId   pgtype.UUID `db:"team_id"`
	TeamRole string      `db:"team_role"`
}

// has all the user info & role to pass to be read client-side
type DBTeamMemberInfo struct {
	DBUser          // embed the DBUser fields into the struct
	TeamRole string `db:"team_role"`
}

// For team_member table.
type DBTeamMember struct {
	TeamId    pgtype.UUID      `db:"team_id"`
	UserId    pgtype.UUID      `db:"user_id"`
	TeamRole  string           `db:"team_role"`
	CreatedOn pgtype.Timestamp `db:"created_on"`
}
//json:"createdOn-hidden"

// type DBTeamAndMember struct {
// 	DBTeam		// table teams
// 	DBTeamMember DBTeamMember // table team_members
// 	DisplayName string	`db:"display_name"` // table_user
// }

// type TeamAndMember struct {
// 	DBTeam
// 	TeamMembers	[]TeamMember
// }
// type TeamMember struct {
// 	DBTeamMember 
// 	DisplayName string
// }

type DBUserTeams struct {
	DBTeam 
	DisplayName	 	string		`db:"display_name"`
	TeamRole		string 		`db:"team_role"`
	AvatarId		string		`db:"avatar_id"`
}

func CreateTeam(team DBTeam) (pgtype.UUID, error) {
	team, err := GetRow[DBTeam](
		`INSERT INTO teams
            (event_id, name, visibility, timezone, technologies, availability, description, invite_code)
            VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, event_id, name, visibility, timezone, technologies, availability, description, created_on, invite_code
		`,
		team.EventId, team.Name, team.Visibility, team.Timezone, team.Technologies, team.Availability, team.Description, team.InviteCode)
	if err != nil {
		fmt.Println("ERROR: failed to create team: ", err)
	}
	return team.Id, err
}

// stepp 5: used to construct the GetTeamResponse struct in server/teams.go
func GetTeam(teamId pgtype.UUID) (DBTeam, error) {
	team, err := GetRow[DBTeam](
		`SELECT 
			teams.id,
			teams.event_id,
			teams.name,
			teams.visibility,
			teams.timezone,
			teams.technologies,
			teams.availability,
			teams.description,
			teams.created_on,
			teams.invite_code
		FROM teams
		WHERE teams.id = $1`,
		teamId)
	// `SELECT * FROM teams WHERE id = $1`,
	// teamId)
	if err != nil {
		logger.Error("===DB/GetTeam error: ", err)
		return DBTeam{}, err
	}
	return team, nil
}

func GetTeamByInvite(inviteCode string) (DBTeam, error) {
	team, err := GetRow[DBTeam](
		`SELECT 
			teams.id,
			teams.event_id,
			teams.name,
			teams.visibility,
			teams.timezone,
			teams.technologies,
			teams.availability,
			teams.description,
			teams.created_on,
			teams.invite_code
		FROM teams
		WHERE teams.invite_code = $1`,
		inviteCode)
	if err != nil {
		logger.Error("===DB/GetTeamByInvite error: ", err)
		return DBTeam{}, err
	}
	return team, nil
}

type DBTeamAndMember struct {
	Id           	pgtype.UUID      `db:"id"`
	EventId      	pgtype.UUID      `db:"event_id"`
	Name         	string           `db:"name"`
	Visibility   	string           `db:"visibility"`
	Timezone     	string           `db:"timezone"`
	Technologies 	string           `db:"technologies"`
	Availability 	string           `db:"availability"`
	Description  	string           `db:"description"`
	InviteCode   	string           `db:"invite_code"`
	TeamId   		pgtype.UUID      `db:"team_id"`
	UserId    		pgtype.UUID      `db:"user_id"`
	TeamRole 	 	string           `db:"team_role"`
	DisplayName 	string 			 `db:"display_name"`
	AvatarId		string			 `db:"avatar_id"`
	ServiceUserId 	string			 `db:"service_user_id"`
}

type UITeam struct {
	Id           pgtype.UUID      `db:"id"`
	EventId      pgtype.UUID      `db:"event_id"`
	Name         string           `db:"name"`
	Visibility   string           `db:"visibility"`
	Timezone     string           `db:"timezone"`
	Technologies string           `db:"technologies"`
	Availability string           `db:"availability"`
	Description  string           `db:"description"`
	InviteCode   string           `db:"invite_code"`
}

type UITeamMember struct {
	TeamId    pgtype.UUID      	  `db:"team_id"`
	UserId    pgtype.UUID      	  `db:"user_id"`
	TeamRole  string          	  `db:"team_role"`
}

type TeamMember struct {
	UITeamMember 
	DisplayName	    string 		`db:"display_name"`
	AvatarUrl		string 		`db:"avatar_id"`
	ServiceUserId 	string 		`db:"service_user_id"`
}

type TeamAndMember struct {
	UITeam
	TeamMembers	[]TeamMember
}

func MapToTeamAndMember(data []DBTeamAndMember) []TeamAndMember{
	// instantiates array to store output, mapped by team id (uuid) for key
	teamMap := make(map[pgtype.UUID]*TeamAndMember)
	for _, item := range data {
		// Check if team already exists in the map. This map loopkup returns:
		// 1) value associated with the key if it exsits
		// 2) boolean indicating whether key was found in the map
		team, ok := teamMap[item.TeamId]
		if !ok {
			// Create a new team
			team = &TeamAndMember{
				UITeam: UITeam {
					Id: 		item.TeamId,
					EventId:      item.EventId,
					Name:         item.Name,
					Visibility:   item.Visibility,
					Timezone:     item.Timezone,
					Technologies: item.Technologies,
					Availability: item.Availability,
					Description:  item.Description,
					InviteCode:   item.InviteCode,
				},
				TeamMembers: []TeamMember{},
			}
			teamMap[item.TeamId] = team
		}
		// Add team member to TeamMembers slice
		member := TeamMember{
			UITeamMember: UITeamMember{
				TeamId:		item.TeamId,
				UserId: 	item.UserId,
				TeamRole: 	item.TeamRole,
			},
			DisplayName: item.DisplayName,
			AvatarUrl: item.AvatarId,
			ServiceUserId: item.ServiceUserId,
		}
		team.TeamMembers = append(team.TeamMembers, member)
	}
	// Convert map back to slice
	var result []TeamAndMember
	for _, team := range teamMap {
		result = append(result, *team)
	}
	fmt.Println(result)
	return result
}

func GetTeams() (*[]TeamAndMember, error){
	teamAndMember, err := GetRows[DBTeamAndMember]( // returns { team 1: { userA: {display_name: "momo"}}, team 1...}
		`SELECT 
			t.id,
			t.event_id, 
			t.name, 
			t.visibility,
			t.timezone,
			t.technologies,
			t.availability,
			t.description,
			t.invite_code,
			u.display_name,
			u.avatar_id,
			u.service_user_id,
			tm.team_id,
			tm.user_id,
			tm.team_role
			FROM teams t
			INNER JOIN team_members tm ON (tm.team_id = t.id)
			INNER JOIN users u ON (u.id = tm.user_id)
            ORDER BY t.id
		`,
		)
		if err != nil {
			return nil, err
		}
	for _, t := range teamAndMember {
		fmt.Printf("%v\n", t)
	}
	UITeamAndMember := MapToTeamAndMember(teamAndMember)
	
	//fmt.Printf("%T\n", t)
	// var TeamMember TeamMember
	// var DisplayName string
	// var TeamAndMember 
	// var TeamAndMembers []TeamAndMember

	// // loop through teamAndMembers, if 

	return &UITeamAndMember, err
}


func GetUserTeams(userId pgtype.UUID) ([]DBUserTeams, error) {
	result, err := GetRows[DBUserTeams](
		`SELECT
			u.display_name,
			t.*, 
			tm.team_role
		FROM users u
		INNER JOIN team_members tm ON u.id = tm.user_id
		INNER JOIN teams t ON tm.team_id = t.id
		WHERE u.id = $1`,
		userId)
	if err != nil {
		fmt.Println("that didn't work: database.GetUserTeams")
		return nil, err
	}

	return result, err  // Try look at the table
}

func UpdateTeam(team DBTeam) (DBTeam, error) {
	event, err := GetRow[DBTeam](
		`UPDATE teams
            SET name=$2,
                visibility=$3,
				timezone=$4,
				technologies=$5,
				availability=$6,
				description=$7,
		WHERE id=$1
		RETURNING *`,
		team.Id, team.Name, team.Visibility, team.Timezone, team.Technologies, team.Availability, team.Description)
	return event, err
}

// fields: userid, teamid, role
// called at server/teams.go createTeam & when someone clicks "join team" 
// DONT MESS WITH BELOW. IT WORKS.
func AddTeamMember(userId pgtype.UUID, teamUUID pgtype.UUID, role string) (userID pgtype.UUID, err error) {
	// userId prints something like: {[22 162 173 240 222 76 79 42 174 62 196 207 243 22 25 78] true}
	
	teamMember, err := GetRow[CreateTeamMember](
		`INSERT INTO team_members
			(user_id, team_id, team_role)
			VALUES ($1, $2, $3)
		RETURNING user_id, team_id, team_role`, userId, teamUUID, role)
	if err != nil {
		fmt.Println(err)
		return userId, err
	}
	return teamMember.UserId, err
}

func GetMembersByTeamId(teamId pgtype.UUID) (*[]DBTeamMemberInfo, error) {
	// In Go, you never return slice-data.
	// Having * in sig means I'm returning the slice-header, which means I need & in my return
	// Not having * means I'm returning a small copy of the slice-header, no need for & in my return
	members, err := GetRows[DBTeamMemberInfo](
		// select all the info of a user (a user row) and their tm.role ()
		`SELECT u.*, tm.team_role
			FROM team_members tm
			INNER JOIN users u on (u.id = tm.user_id)
			WHERE tm.team_id = $1`,
		teamId)
	return &members, err
}
