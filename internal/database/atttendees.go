package database

import "database/sql"

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	ID      int `json:"id"`
	UserID  int `json:"userId"`
	EventID int `json:"eventId"`
}
