package srv

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"mlauth/pkg/conf"
	"strconv"
	"time"
)

func genToken(uid int, exp int) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"uid": strconv.Itoa(uid),
		"exp": time.Now().Add(time.Duration(exp) * time.Second).Unix(),
	})
	token, err := t.SignedString(conf.SecretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenAccessToken(uid int) (string, error) {
	return genToken(uid, conf.AccessTokenAge)
}

func GenUpdateToken(uid int) (string, error) {
	return genToken(uid, conf.UpdateTokenAge)
}

func checkToken(token string, exp int) (int, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		} else {
			return conf.SecretKey, nil
		}
	})
	if err != nil {
		return 0, err
	}

	claimsErrMsg := "invalid claims"
	c, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return 0, fmt.Errorf(claimsErrMsg)
	}

	u, ok := c["uid"]
	if !ok {
		return 0, fmt.Errorf(claimsErrMsg)
	}

	us, ok := u.(string)
	if !ok {
		return 0, fmt.Errorf(claimsErrMsg)
	}

	uid, err := strconv.Atoi(us)
	if err != nil {
		return 0, err
	}

	e, ok := c["exp"]
	if !ok {
		return 0, fmt.Errorf(claimsErrMsg)
	}

	ef, ok := e.(float64)
	if !ok {
		return 0, fmt.Errorf(claimsErrMsg)
	}

	et := time.Unix(int64(ef), 0)
	if time.Now().Sub(et).Seconds() > float64(exp) {
		return 0, fmt.Errorf("token expired")
	}

	return uid, nil
}

func CheckAccessToken(token string) (int, error) {
	return checkToken(token, conf.AccessTokenAge)
}

func CheckUpdateToken(token string) (int, error) {
	return checkToken(token, conf.UpdateTokenAge)
}
