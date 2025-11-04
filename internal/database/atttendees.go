package database

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	ID      int `json:"id"`
	UserID  int `json:"userId"`
	EventID int `json:"eventId"`
}

func (a *AttendeeModel) Insert(attendee *Attendee) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO attendees (event_id, user_id) VALUES ($1, $2) RETURNING id`

	err := a.DB.QueryRowContext(
		ctx,
		query,
		attendee.EventID,
		attendee.UserID,
	).Scan(&attendee.ID)

	if err != nil {
		return nil, err
	}

	return attendee, nil
}

func (a *AttendeeModel) GetByEventAndAttendee(eventID, userID int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, event_id, user_id FROM attendees WHERE event_id = $1 AND user_id = $2`

	var attendee Attendee
	err := a.DB.QueryRowContext(
		ctx,
		query,
		eventID,
		userID,
	).Scan(
		&attendee.ID,
		&attendee.EventID,
		&attendee.UserID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &attendee, nil
}

func (a *AttendeeModel) GetAttendeesByEvent(eventID int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT u.id, u.name, u.email
		FROM users u
		JOIN attendees a ON u.id = a.user_id
		WHERE a.event_id = $1
	`

	rows, err := a.DB.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var attendees []*User

	for rows.Next() {
		var attendee User
		err := rows.Scan(&attendee.ID, &attendee.Name, &attendee.Email)
		if err != nil {
			return nil, err
		}

		attendees = append(attendees, &attendee)
	}

	return attendees, nil
}

func (a *AttendeeModel) Delete(eventID, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM attendees WHERE event_id = $1 AND user_id = $2`

	res, err := a.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		log.Println("Attendee removal failed. No rows exist in db!!")
	}

	return nil
}
