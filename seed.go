package main

import "time"

// Loan is a struct that represents a loan in this system.
// TableName: dt_m_loan
type Loan struct {
	CreatedAt time.Time `json:"created_at,format:'2006-01-02'" db:"created_at"`
	Code      string    `json:"name" db:"code"`
	ID        int64     `json:"id,omitzero" db:"id"`
	Amount    int       `json:"age" db:"age"`
}
