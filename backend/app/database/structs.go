package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type DBTeam struct {
	Id           pgtype.UUID        `db:"id"`
	EventId      pgtype.UUID        `db:"event_id"`
	Name         string             `db:"name"`
	Visibility   string             `db:"visibility"`
	Timezone     string             `db:"timezone"`
	Technologies string             `db:"technologies"`
	Availability string             `db:"availability"`
	Description  string             `db:"description"`
	CreatedOn    pgtype.Timestamptz `db:"created_on"`
	InviteCode   string             `db:"invite_code"`
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

type DBTeamMember struct {
	TeamId    pgtype.UUID      `db:"team_id"`
	UserId    pgtype.UUID      `db:"user_id"`
	TeamRole  string           `db:"team_role"`
	CreatedOn pgtype.Timestamp `db:"user_created_on"`
}

type DBUserTeams struct {
	DBTeam
	DisplayName string `db:"display_name"`
	TeamRole    string `db:"team_role"`
	AvatarId    string `db:"avatar_id"`
}

type DBTeamAndTeamMember struct {
	Id                  pgtype.UUID        `db:"id"`
	EventId             pgtype.UUID        `db:"event_id"`
	Name                string             `db:"name"`
	Visibility          string             `db:"visibility"`
	Timezone            string             `db:"timezone"`
	Technologies        string             `db:"technologies"`
	Availability        string             `db:"availability"`
	Description         string             `db:"description"`
	CreatedOn           pgtype.Timestamptz `db:"team_created_on"`
	InviteCode          string             `db:"invite_code"`
	MembershipCreatedOn pgtype.Timestamptz `db:"membership_created_on"`
	TeamId              pgtype.UUID        `db:"team_id"`
	UserId              pgtype.UUID        `db:"user_id"`
	TeamRole            string             `db:"team_role"`
	DisplayName         string             `db:"display_name"`
	AvatarId            string             `db:"avatar_id"`
	ServiceUserId       string             `db:"service_user_id"`
}

type UITeam struct {
	Id           pgtype.UUID        `db:"id"`
	EventId      pgtype.UUID        `db:"event_id"`
	Name         string             `db:"name"`
	Visibility   string             `db:"visibility"`
	Timezone     string             `db:"timezone"`
	Technologies string             `db:"technologies"`
	Availability string             `db:"availability"`
	Description  string             `db:"description"`
	CreatedOn    pgtype.Timestamptz `db:"team_created_on"`
	InviteCode   string             `db:"invite_code"`
}

type UITeamMember struct {
	TeamId              pgtype.UUID        `db:"team_id"`
	UserId              pgtype.UUID        `db:"user_id"`
	MembershipCreatedOn pgtype.Timestamptz `db:"membership_created_on"`
	TeamRole            string             `db:"team_role"`
}

type TeamMember struct {
	UITeamMember
	DisplayName   string `db:"display_name"`
	AvatarId      string `db:"avatar_id"`
	ServiceUserId string `db:"service_user_id"`
}

type TeamAndMembers struct {
	UITeam
	TeamMembers []TeamMember
}

type UserTeamMember struct {
	UserId        pgtype.UUID        `json:"UserId"`
	DisplayName   string             `json:"DisplayName"`
	TeamRole      string             `json:"TeamRole"`
	AvatarId      string             `json:"AvatarId"`
	ServiceUserId string             `json:"ServiceUserId"`
	UserCreatedOn pgtype.Timestamptz `json:"UserCreatedOn"`
}

type UserTeam struct {
	Id            pgtype.UUID        `json:"Id"`
	EventId       pgtype.UUID        `json:"EventId"`
	Name          string             `json:"Name"`
	Description   string             `json:"Description"`
	Visibility    string             `json:"Visibility"`
	Technologies  string             `json:"Technologies"`
	Availability  string             `json:"Availability"`
	TeamCreatedOn pgtype.Timestamptz `json:"TeamCreatedOn"`
	InviteCode    string             `json:"InviteCode"`
	TeamMembers   []UserTeamMember   `json:"TeamMembers"`
}

type UserTeamAndMembers struct {
	Id                  pgtype.UUID        `db:"id"`
	EventId             pgtype.UUID        `db:"event_id"`
	Name                string             `db:"name"`
	Visibility          string             `db:"visibility"`
	Timezone            string             `db:"timezone"`
	Technologies        string             `db:"technologies"`
	Availability        string             `db:"availability"`
	Description         string             `db:"description"`
	TeamCreatedOn       pgtype.Timestamptz `db:"team_created_on"`
	InviteCode          string             `db:"invite_code"`
	TeamId              pgtype.UUID        `db:"team_id"`
	UserId              pgtype.UUID        `db:"user_id"`
	TeamRole            string             `db:"team_role"`
	MembershipCreatedOn pgtype.Timestamptz `db:"membership_created_on"`
	DisplayName         string             `db:"display_name"`
	AvatarId            string             `db:"avatar_id"`
	ServiceUserId       string             `db:"service_user_id"`
	CurrentUserRole     string             `db:"current_user_role"`
}
