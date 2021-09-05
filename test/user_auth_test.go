package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mlauth/pkg/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func userLogin(t *testing.T, r *gin.Engine) (string, string) {
	w := httptest.NewRecorder()
	b, err := json.Marshal(gin.H{
		"username": "testusername",
		"password": "testpassword",
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
	assert.Equal(t, body.Username, "testusername")
	assert.NotEqual(t, body.AccessToken, "")
	assert.NotEqual(t, body.UpdateToken, "")

	return body.AccessToken, body.UpdateToken
}

func TestUserLogin(t *testing.T) {
	r := api.Route()
	_, _ = userLogin(t, r)
}

func TestUserRenew(t *testing.T) {
	r := api.Route()
	_, ut := userLogin(t, r)
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
