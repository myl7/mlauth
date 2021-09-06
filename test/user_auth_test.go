package test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestUserLogin(t *testing.T) {
	userLogin(t, "testusername", "testpassword")
}

func TestUserRenew(t *testing.T) {
	_, ut := userLogin(t, "testusername", "testpassword")
	b := serJson(t, gin.H{
		"update_token": ut,
	})
	w := reqApi(t, "POST", "/api/users/renew", b, nil)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body := struct {
		UpdateToken string `json:"update_token"`
		AccessToken string `json:"access_token"`
	}{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, ut, body.UpdateToken)
}

func TestUserRecover(t *testing.T) {
	b := serJson(t, gin.H{
		"username": "username4",
	})
	w := reqApi(t, "POST", "/api/users/recover", b, nil)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	email, emailBody := getEmailInfo(t)
	assert.Equal(t, "user4@128.com", email)

	re := regexp.MustCompile(`/emails/recover/?\?recover-code=[0-9a-z-]+`)
	p := re.Find([]byte(emailBody))
	assert.NotNil(t, p)

	b = serJson(t, gin.H{
		"password": "passwordRecover",
	})
	w = reqApi(t, "POST", "/api"+string(p), b, nil)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	at, _ := userLogin(t, "username4", "passwordRecover")
	body := getUserDetail(t, at)
	assert.Equal(t, 4, body.Uid)
}
