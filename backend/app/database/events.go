package database

import "github.com/jackc/pgx/v5/pgtype"

type DBEvent struct {
	Id              pgtype.UUID      `db:"id"`
	StatusId        int              `db:"status_id"`
	Title           string           `db:"title"`
	Description     string           `db:"description"`
	Rules           string           `db:"rules"`
	Timeline        string           `db:"timeline"`
	OrganizerUserId pgtype.UUID      `db:"organizer_user_id" json:"-"`
	MaxTeams        int              `db:"max_teams"`
	StartsAt        pgtype.Timestamp `db:"starts_at"`
	EndsAt          pgtype.Timestamp `db:"ends_at"`
	CreatedOn       pgtype.Timestamp `db:"created_on" json:"-"`
}

type DBEventStatus struct {
	Id          int    `db:"id"`
	Code        string `db:"code"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

func CreateEvent(organizerUserId pgtype.UUID) (DBEvent, error) {
	var event DBEvent
	event, err := GetRow[DBEvent](
		`INSERT INTO events
           (status_id, title, description, rules, organizer_user_id)
         VALUES
           ((SELECT id FROM statuses WHERE code = 'PLANNING'),
            '',
            '',
            '',
            $1)
         RETURNING *
         `,
		organizerUserId)
	return event, err
}

func GetEvent(eventId pgtype.UUID) (DBEvent, error) {
	event, err := GetRow[DBEvent](
		`SELECT * FROM events WHERE id = $1`,
		eventId)
	return event, err
}

func GetEvents() ([]DBEvent, error) {
	result, err := GetRows[DBEvent](`SELECT * FROM events`)
	return result, err
}

// GetActiveEvent will return what is assumed to be a single active event.
func GetActiveEvent() (DBEvent, error) {
	result, err := GetRow[DBEvent](
		`SELECT * FROM events
         WHERE status_id IN (
           SELECT id FROM statuses WHERE code NOT IN ('PLANNING', 'ENDED') 
         )
         LIMIT 1`)
	return result, err
}

func UpdateEvent(event DBEvent) (DBEvent, error) {
	event, err := GetRow[DBEvent](
		`UPDATE events
         SET status_id=$2,
             title=$3,
             description=$4,
             rules=$5,
             max_teams=$6,
             starts_at=$7,
             ends_at=$8
         WHERE id=$1
         RETURNING *`,
		event.Id, event.StatusId, event.Title, event.Description, event.Rules, event.MaxTeams, event.StartsAt, event.EndsAt)
	return event, err
}

func GetStatuses() ([]DBEventStatus, error) {
	statuses, err := GetRows[DBEventStatus](`SELECT * FROM statuses`)
	return statuses, err
}
