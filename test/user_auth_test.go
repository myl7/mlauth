package test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mlauth/pkg/api"
	"mlauth/pkg/srv"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

func userLogin(t *testing.T, r *gin.Engine, username string, password string) (string, string) {
	w := httptest.NewRecorder()
	b, err := json.Marshal(gin.H{
		"username": username,
		"password": password,
	})
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/users/login", bytes.NewReader(b))
	assert.NoError(t, err)

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body := struct {
		Username    string `json:"username"`
		AccessToken string `json:"access_token"`
		UpdateToken string `json:"update_token"`
	}{}
	err = json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, body.Username, username)
	assert.NotEqual(t, body.AccessToken, "")
	assert.NotEqual(t, body.UpdateToken, "")

	return body.AccessToken, body.UpdateToken
}

func TestUserLogin(t *testing.T) {
	r := api.Route()
	_, _ = userLogin(t, r, "testusername", "testpassword")
}

func TestUserRenew(t *testing.T) {
	r := api.Route()
	_, ut := userLogin(t, r, "testusername", "testpassword")
	w := httptest.NewRecorder()
	b, err := json.Marshal(gin.H{
		"update_token": ut,
	})
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/users/renew", bytes.NewReader(b))
	assert.NoError(t, err)

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body := struct {
		UpdateToken string `json:"update_token"`
		AccessToken string `json:"access_token"`
	}{}
	err = json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, ut, body.UpdateToken)
}

func TestUserRecover(t *testing.T) {
	r := api.Route()
	b, err := json.Marshal(gin.H{
		"username": "username4",
	})
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/users/recover", bytes.NewReader(b))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	var email, emailBody string
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	select {
	case email = <-srv.SendEmailMockChan:
	case <-ctx.Done():
		assert.NotEmpty(t, email, "Can not get email")
	}
	cancel()
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	select {
	case emailBody = <-srv.SendEmailMockChan:
	case <-ctx.Done():
		assert.NotEmpty(t, emailBody, "Can not get email body")
	}
	cancel()
	assert.Equal(t, "user4@128.com", email)

	re := regexp.MustCompile(`/emails/recover/?\?recover-code=[0-9a-z-]+`)
	p := re.Find([]byte(emailBody))
	assert.NotNil(t, p)

	b, err = json.Marshal(gin.H{
		"password": "passwordRecover",
	})
	assert.NoError(t, err)

	req, err = http.NewRequest("POST", "/api"+string(p), bytes.NewReader(b))
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	at, _ := userLogin(t, r, "username4", "passwordRecover")
	body := getUserDetail(t, r, at)
	assert.Equal(t, 4, body.Uid)
}
