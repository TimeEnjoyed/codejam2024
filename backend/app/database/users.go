package database

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

type DBUser struct {
	Id            pgtype.UUID      `db:"id"`
	ServiceName   string           `db:"service_name"`
	ServiceUserId string           `db:"service_user_id"`
	Role          string           `db:"role"`
	DisplayName   string           `db:"display_name"`
	AvatarId      string          `db:"avatar_id"`
	CreatedOn     pgtype.Timestamp `db:"created_on" json:"-"`
}

// Roles
const (
	Admin = "ADMIN"
)

func CreateUser(serviceName string, serviceUserId string, serviceDisplayName string, avatarId string) DBUser {
	user, err := GetRow[DBUser](
		`INSERT INTO users (service_name, service_user_id, display_name, avatar_id)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (service_name, service_user_id)
		 DO UPDATE
		 SET display_name = $3, avatar_id = $4
		 RETURNING *`,
		serviceName, serviceUserId, serviceDisplayName, avatarId)
	if err != nil {
		fmt.Println(serviceName, serviceUserId, serviceDisplayName, avatarId)
		logger.Error("error getting user: %v", err)
	}
	return user
}

func GetUser(userId pgtype.UUID) DBUser {
	user, err := GetRow[DBUser](
		`SELECT *
		 FROM users 
		 WHERE id = $1`,
		userId)
	if err != nil {
		logger.Error("error getting user: %v", err)
	}
	return user
}

func UpdateUser(user DBUser) (DBUser, error) {
	user, err := GetRow[DBUser](
		`UPDATE users
	 	 SET display_name = $2
	     WHERE id = $1
	     RETURNING *`,
		user.Id, user.DisplayName)
	if err != nil {
		logger.Error("error updating user: %v", err)
	}
	return user, err
}
