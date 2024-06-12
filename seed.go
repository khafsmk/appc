package main

import "time"

// User is a struct that represents a user
// TableName: users
type User struct {
	CreatedAt time.Time `json:"created_at,format:'2006-01-02'" db:"created_at"`
	Name      string    `json:"name" db:"name"`
	ID        int64     `json:"id,omitzero" db:"id"`
	Age       int       `json:"age" db:"age"`
}
