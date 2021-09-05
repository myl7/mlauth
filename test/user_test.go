package test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"mlauth/pkg/api"
	"mlauth/pkg/srv"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

func TestUserGet(t *testing.T) {
	r := api.Route()
	at, _ := userLogin(t, r)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/api/users/me", nil)
	assert.NoError(t, err)

	req.Header.Set("x-access-token", at)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body := struct {
		Uid         int       `json:"uid"`
		Username    string    `json:"username"`
		Email       string    `json:"email"`
		DisplayName string    `json:"display_name"`
		IsActive    bool      `json:"is_active"`
		CreatedAt   time.Time `json:"created_at"`
	}{}
	err = json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
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
	r := api.Route()
	b, err := json.Marshal(gin.H{
		"username":     "testU",
		"password":     "testPassYou",
		"email":        "testE@outlook.com",
		"display_name": "符号看象限ラブライブ한국어",
	})
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/users", bytes.NewReader(b))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body := struct {
		Uid         int       `json:"uid"`
		Username    string    `json:"username"`
		Email       string    `json:"email"`
		DisplayName string    `json:"display_name"`
		IsActive    bool      `json:"is_active"`
		CreatedAt   time.Time `json:"created_at"`
	}{}
	err = json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "testU", body.Username)
	assert.Equal(t, "testE@outlook.com", body.Email)
	assert.Equal(t, "符号看象限ラブライブ한국어", body.DisplayName)
	assert.Equal(t, false, body.IsActive)
	assert.Less(t, time.Now().UTC().Sub(body.CreatedAt), 5*time.Second)

	var email, emailBody string
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	select {
	case email = <-srv.SendEmailMockChan:
	case <-ctx.Done():
		log.Println("Can not get email")
	}
	cancel()
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	select {
	case emailBody = <-srv.SendEmailMockChan:
	case <-ctx.Done():
		log.Println("Can not get email body")
	}
	cancel()
	assert.Equal(t, "testE@outlook.com", email)

	re := regexp.MustCompile(`/emails/active/?\?active-code=[0-9a-z-]+`)
	p := re.Find([]byte(emailBody))
	assert.NotNil(t, p)

	req, err = http.NewRequest("POST", "/api"+string(p), nil)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	req, err = http.NewRequest("GET", "/api/users/me", nil)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	b, err = json.Marshal(gin.H{
		"username": "testU",
		"password": "testPassYou",
	})
	assert.NoError(t, err)

	req, err = http.NewRequest("POST", "/api/users/login", bytes.NewReader(b))
	assert.NoError(t, err)

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body2 := struct {
		Username    string `json:"username"`
		AccessToken string `json:"access_token"`
		UpdateToken string `json:"update_token"`
	}{}
	err = json.Unmarshal(w.Body.Bytes(), &body2)
	assert.NoError(t, err)
	assert.Equal(t, body2.Username, "testU")
	assert.NotEqual(t, body2.AccessToken, "")
	assert.NotEqual(t, body2.UpdateToken, "")

	at := body2.AccessToken
	w = httptest.NewRecorder()

	req, err = http.NewRequest("GET", "/api/users/me", nil)
	assert.NoError(t, err)

	req.Header.Set("x-access-token", at)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body3 := struct {
		Uid         int       `json:"uid"`
		Username    string    `json:"username"`
		Email       string    `json:"email"`
		DisplayName string    `json:"display_name"`
		IsActive    bool      `json:"is_active"`
		CreatedAt   time.Time `json:"created_at"`
	}{}
	err = json.Unmarshal(w.Body.Bytes(), &body3)
	assert.NoError(t, err)
	assert.Equal(t, body.Uid, body3.Uid)
	assert.Equal(t, true, body3.IsActive)
}
