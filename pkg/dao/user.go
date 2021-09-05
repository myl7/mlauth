package dao

import (
	"fmt"
	"mlauth/pkg/mdl"
)

func SelectUser(uid int) (mdl.User, error) {
	db, err := getDb()
	if err != nil {
		return mdl.User{}, err
	}

	u := mdl.User{}
	sql := "SELECT * FROM users WHERE uid = $1"
	err = db.Get(&u, sql, uid)
	if err != nil {
		return mdl.User{}, err
	}

	return u, nil
}

func SelectUserByUsername(username string) (mdl.User, error) {
	db, err := getDb()
	if err != nil {
		return mdl.User{}, err
	}

	u := mdl.User{}
	sql := "SELECT * FROM users WHERE username = $1"
	err = db.Get(&u, sql, username)
	if err != nil {
		return mdl.User{}, err
	}

	return u, nil
}

func UpdateUser(uid int, uEdit mdl.User) (mdl.User, error) {
	uEdit.Uid = uid
	db, err := getDb()
	if err != nil {
		return mdl.User{}, err
	}

	u := mdl.User{}
	sql := `UPDATE users SET display_name = :display_name, password = :password, email = :email, is_active = :is_active
        WHERE uid = :uid RETURNING *`
	rows, err := db.NamedQuery(sql, uEdit)
	if err != nil {
		return mdl.User{}, err
	}

	if !rows.Next() {
		return mdl.User{}, fmt.Errorf("no return values")
	}

	err = rows.StructScan(&u)
	if err != nil {
		return mdl.User{}, err
	}

	return u, nil
}

func InsertUser(uCreate mdl.User) (mdl.User, error) {
	db, err := getDb()
	if err != nil {
		return mdl.User{}, err
	}

	u := mdl.User{}
	sql := `INSERT INTO users (username, password, email, display_name, is_active, is_super, created_at)
        VALUES (:username, :password, :email, :display_name, :is_active, :is_super, :created_at) RETURNING *`
	rows, err := db.NamedQuery(sql, uCreate)
	if err != nil {
		return mdl.User{}, err
	}

	if !rows.Next() {
		return mdl.User{}, fmt.Errorf("no return values")
	}

	err = rows.StructScan(&u)
	if err != nil {
		return mdl.User{}, err
	}

	return u, nil
}
