package database

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

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
	// "/invite/:invitecode"
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

func MapToTeamAndMember(data []DBTeamAndTeamMember) []TeamAndMembers {
	// instantiates array to store output, mapped by team id (uuid) for key
	teamMap := make(map[pgtype.UUID]*TeamAndMembers)
	for _, item := range data {
		// Check if team already exists in the map. This map loopkup returns:
		// 1) value associated with the key if it exsits
		// 2) boolean indicating whether key was found in the map
		team, ok := teamMap[item.TeamId]
		if !ok {
			// Create a new team
			team = &TeamAndMembers{
				UITeam: UITeam{
					Id:           item.Id,
					EventId:      item.EventId,
					Name:         item.Name,
					Visibility:   item.Visibility,
					Timezone:     item.Timezone,
					Technologies: item.Technologies,
					Availability: item.Availability,
					Description:  item.Description,
					CreatedOn:    item.CreatedOn,
					InviteCode:   item.InviteCode,
				},
				TeamMembers: []TeamMember{},
			}
			teamMap[item.TeamId] = team
		}
		// Add team member to TeamMembers slice
		member := TeamMember{
			UITeamMember: UITeamMember{
				TeamId:              item.TeamId,
				UserId:              item.UserId,
				MembershipCreatedOn: item.MembershipCreatedOn,
				TeamRole:            item.TeamRole,
			},
			DisplayName:   item.DisplayName,
			AvatarId:      item.AvatarId,
			ServiceUserId: item.ServiceUserId,
		}
		team.TeamMembers = append(team.TeamMembers, member)
	}
	// Convert map back to slice
	var result []TeamAndMembers
	for _, team := range teamMap {
		result = append(result, *team)
	}
	fmt.Println(result)
	return result
}

func GetTeams() (*[]TeamAndMembers, error) {
	teamAndMember, err := GetRows[DBTeamAndTeamMember]( // returns { team 1: { userA: {display_name: "momo"}}, team 1...}
		`SELECT 
			t.id,
			t.event_id, 
			t.name, 
			t.visibility,
			t.timezone,
			t.technologies,
			t.availability,
			t.description,
            t.created_on AS team_created_on,
			t.invite_code,
			u.display_name,
			u.avatar_id,
			u.service_user_id,
			tm.team_id,
			tm.user_id,
			tm.created_on AS membership_created_on,
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

	return &UITeamAndMember, err
}

func MapToUserTeamAndMember(userTeams *[]UserTeamAndMembers) []UserTeam {
	// initialize a json dictionary
	teamMap := make(map[pgtype.UUID]*UserTeam)

	// loop through all the rows from the query.
	// UserTeams aka []UserTeam is a single row from query.
	for _, ut := range *userTeams {
		// If team doesn't exist in map, create it
		if _, exists := teamMap[ut.Id]; !exists {
			teamMap[ut.Id] = &UserTeam{
				Id:            ut.Id,
				EventId:       ut.EventId,
				Name:          ut.Name,
				Description:   ut.Description,
				Visibility:    ut.Visibility,
				Technologies:  ut.Technologies,
				Availability:  ut.Availability,
				TeamCreatedOn: ut.TeamCreatedOn,
				InviteCode:    ut.InviteCode,
				TeamMembers:   []UserTeamMember{},
			}
		}

		// create a member
		TeamMember := UserTeamMember{
			UserId:        ut.UserId,
			DisplayName:   ut.DisplayName,
			TeamRole:      ut.TeamRole,
			AvatarId:      ut.AvatarId,
			ServiceUserId: ut.ServiceUserId,
			UserCreatedOn: ut.MembershipCreatedOn,
		}

		// add member to teamMap dictionary
		teamMap[ut.Id].TeamMembers = append(teamMap[ut.Id].TeamMembers, TeamMember)
	}

	// Convert map to slice
	teams := make([]UserTeam, 0, len(teamMap))
	for _, team := range teamMap {
		teams = append(teams, *team)
	}
	return teams
}

func GetUserTeams(userId pgtype.UUID) (*[]UserTeam, error) {
	result, err := GetRows[UserTeamAndMembers](
		`SELECT 
            t.id,
            t.event_id,
            t.name,
            t.visibility,
            t.timezone,
            t.technologies,
            t.availability,
            t.description,
            t.created_on AS team_created_on,
			t.invite_code,
			tm.team_id,
			tm.user_id, 
            tm.team_role,
            tm.created_on AS membership_created_on,
			tmu.display_name,
            tmu.avatar_id,
            tmu.service_user_id,
			mtm.team_role AS current_user_role
		FROM team_members mtm
		INNER JOIN teams t ON t.id = mtm.team_id
		INNER JOIN team_members tm ON t.id = tm.team_id
		INNER JOIN users tmu ON tmu.id = tm.user_id
		WHERE mtm.user_id = $1
		ORDER BY team_created_on`,
		userId)
	if err != nil {
		logger.Error("failed to retrieve team info for user: %v", userId, err)
	}

	UIUserTeamAndMember := MapToUserTeamAndMember(&result)
	return &UIUserTeamAndMember, err
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

type DBTeamInviteCode struct {
	InviteCode string `db:"invite_code"`
}

func GetTeamInviteCode(teamId pgtype.UUID) (inviteCode DBTeamInviteCode, err error) {
	fmt.Println(teamId)
	teamInviteCode, err := GetRow[DBTeamInviteCode](
		`SELECT teams.invite_code
		FROM teams 
		WHERE teams.id = $1`, teamId)
	if err != nil {
		logger.Error("failed to retrieve invite code for teamId %v: %w", teamId, err)
	}
	return teamInviteCode, err
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

func RemoveTeamMember(teamId pgtype.UUID, userId pgtype.UUID) {

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
