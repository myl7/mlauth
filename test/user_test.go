package test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"mlauth/pkg/api"
	"net/http"
	"net/http/httptest"
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
	assert.Equal(t, 200, w.Code, "body:", w.Body.String())

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
