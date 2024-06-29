package database

import (
	"github.com/emicklei/pgtalk/convert"
	"github.com/jackc/pgx/v5/pgtype"
)

type DBUser struct {
	Id              pgtype.UUID      `db:"id" `
	ServiceName     string           `db:"service_name"`
	ServiceUserId   string           `db:"service_user_id"`
	ServiceUserName string           `db:"service_user_name"`
	Role            string           `db:"role"`
	DisplayName     string           `db:"display_name"`
	AvatarUrl       *string          `db:"avatar_url"`
	AccountStatus   string           `db:"account_status"`
	LockDisplayName bool             `db:"lock_display_name"`
	CreatedOn       pgtype.Timestamp `db:"created_on" json:"-"`
}

// Roles
const (
	Admin = "ADMIN"
)

func CreateUser(serviceName string, serviceUserId string, serviceDisplayName string, avatarUrl string) DBUser {
	user, err := GetRow[DBUser](
		`INSERT INTO users (service_name, service_user_id, service_user_name, display_name, avatar_url)
		 VALUES ($1, $2, $3, $3, $4)
		 ON CONFLICT (service_name, service_user_id)
		 DO UPDATE
		 SET service_user_name = $3, avatar_url = $4
		 RETURNING *`,
		serviceName, serviceUserId, serviceDisplayName, avatarUrl)
	if err != nil {
		logger.Error("error getting user: %v", err)
	}
	return user
}

func GetUser(userId pgtype.UUID) (DBUser, error) {
	user, err := GetRow[DBUser](
		`SELECT 
           *
		 FROM users 
		 WHERE id = $1`,
		userId)
	if err != nil {
		logger.Error("error getting user: id: %s, error: %v", convert.UUIDToString(userId), err)
	}
	return user, err
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

func GetAllUsers() ([]DBUser, error) {
	users, err := GetRows[DBUser](
		`SELECT *
         FROM users
         ORDER BY (
           CASE -- put empty roles after actual roles
             WHEN role is not null and length(role) > 0 then role
             ELSE 'zzz'
           END 
         ), service_user_name`)
	return users, err
}

func SetAccountStatus(userId pgtype.UUID, status string) (DBUser, error) {
	user, err := GetRow[DBUser](
		`UPDATE users
         SET account_status = $2
         WHERE id= $1 
         RETURNING *`,
		userId, status)
	if err != nil {
		logger.Error("SetAccountStatus Error: %v", err)
	}
	return user, err
}

func SetDisplayName(userId pgtype.UUID, displayName string) (DBUser, error) {
	user, err := GetRow[DBUser](
		`UPDATE users
         SET display_name = $2
         WHERE id = $1
         RETURNING *`,
		userId, displayName)
	if err != nil {
		logger.Error("SetDisplayName Error: %v", err)
	}
	return user, err
}

func SetDisplayNameLock(userId pgtype.UUID, locked bool) (DBUser, error) {
	user, err := GetRow[DBUser](
		`UPDATE users
         SET lock_display_name = $2
         WHERE id = $1
         RETURNING *`,
		userId, locked)
	if err != nil {
		logger.Error("SetDisplayNameLock Error: %v", err)
	}
	return user, err
}
