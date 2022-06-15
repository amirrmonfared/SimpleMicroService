// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"database/sql"
)

type User struct {
	ID         int64         `json:"id"`
	Email      string        `json:"email"`
	FirstName  string        `json:"first_name"`
	LastName   string        `json:"last_name"`
	Password   string        `json:"password"`
	UserActive sql.NullInt64 `json:"user_active"`
	CreatedAt  sql.NullTime  `json:"created_at"`
	UpdatedAt  sql.NullTime  `json:"updated_at"`
}
