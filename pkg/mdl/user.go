package mdl

import "time"

type User struct {
	Uid         int       `db:"uid"`
	Username    string    `db:"username"`
	Password    string    `db:"password"`
	Email       string    `db:"email"`
	DisplayName string    `db:"display_name"`
	IsActive    bool      `db:"is_active"`
	IsSuper     bool      `db:"is_super"`
	CreatedAt   time.Time `db:"created_at"`
}
