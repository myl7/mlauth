package dao

import "mlauth/pkg/mdl"

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
