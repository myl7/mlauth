package test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mlauth/pkg/api"
	"regexp"
	"testing"
	"time"
)

func TestUserGet(t *testing.T) {
	at, _ := userLogin(t, "testusername", "testpassword")
	body := getUserDetail(t, at)
	assert.Equal(t, 1, body.Uid)
	assert.Equal(t, "testusername", body.Username)
	assert.Equal(t, "testemail@gmail.com", body.Email)
	assert.Equal(t, "test display name", body.DisplayName)
	assert.Equal(t, true, body.IsActive)

	p, err := time.Parse(time.RFC3339, "1999-01-08T04:05:06Z")
	assert.NoError(t, err)
	assert.Equal(t, p, body.CreatedAt)
}

func TestUserRegister(t *testing.T) {
	b := serJson(t, gin.H{
		"username":     "testU",
		"password":     "testPassYou",
		"email":        "testE@outlook.com",
		"display_name": "符号看象限ラブライブ한국어",
	})
	w := reqApi(t, "POST", "/api/users", b, nil)
	body := api.UserDetail{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "testU", body.Username)
	assert.Equal(t, "testE@outlook.com", body.Email)
	assert.Equal(t, "符号看象限ラブライブ한국어", body.DisplayName)
	assert.Equal(t, false, body.IsActive)
	assert.Less(t, time.Now().UTC().Sub(body.CreatedAt), 5*time.Second)

	email, emailBody := getEmailInfo(t)
	assert.Equal(t, "testE@outlook.com", email)

	re := regexp.MustCompile(`/emails/active/?\?active-code=[0-9a-z-]+`)
	p := re.Find([]byte(emailBody))
	assert.NotNil(t, p)

	w = reqApi(t, "POST", "/api"+string(p), nil, nil)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	at, _ := userLogin(t, "testU", "testPassYou")
	body2 := getUserDetail(t, at)
	assert.Equal(t, body.Uid, body2.Uid)
	assert.Equal(t, true, body2.IsActive)
}

func TestUserEditExceptEmail(t *testing.T) {
	at, _ := userLogin(t, "anotherusername", "anotherpassword")
	b := serJson(t, gin.H{
		"password":     "testPassYou",
		"display_name": "符号看象限ラブライブ한국어",
	})
	w := reqApi(t, "PUT", "/api/users/me", b, &at)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body := api.UserDetail{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, body.Uid, body.Uid)
	assert.Equal(t, "符号看象限ラブライブ한국어", body.DisplayName)

	ensureNoEmail(t)
	at, _ = userLogin(t, "anotherusername", "testPassYou")
	getUserDetail(t, at)
}

func TestUserEditEmail(t *testing.T) {
	at, _ := userLogin(t, "username3", "password3")
	b := serJson(t, gin.H{
		"email": "tellmeyouremail@outlook.com",
	})
	w := reqApi(t, "PUT", "/api/users/me", b, &at)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body := api.UserDetail{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, 3, body.Uid)
	assert.Equal(t, "hello@qq.com", body.Email)

	email, emailBody := getEmailInfo(t)
	assert.Equal(t, "tellmeyouremail@outlook.com", email)

	re := regexp.MustCompile(`/emails/change-email/?\?verify-code=[0-9a-z-]+`)
	p := re.Find([]byte(emailBody))
	assert.NotNil(t, p)

	w = reqApi(t, "POST", "/api"+string(p), nil, nil)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body2 := getUserDetail(t, at)
	assert.Equal(t, 3, body2.Uid)
	assert.Equal(t, "tellmeyouremail@outlook.com", body2.Email)
}
