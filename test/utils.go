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

func getUserDetail(t *testing.T, r *gin.Engine, at string) api.UserDetail {
	req, err := http.NewRequest("GET", "/api/users/me", nil)
	assert.NoError(t, err)

	req.Header.Set("x-access-token", at)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body := api.UserDetail{}
	err = json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)

	return body
}
