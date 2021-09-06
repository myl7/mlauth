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
	"testing"
	"time"
)

var Router = api.Route()

func reqApi(t *testing.T, method string, path string, body []byte, at *string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, bytes.NewReader(body))
	assert.NoError(t, err)

	if at != nil {
		req.Header.Set("x-access-token", *at)
	}

	w := httptest.NewRecorder()
	Router.ServeHTTP(w, req)
	return w
}

func serJson(t *testing.T, in interface{}) []byte {
	b, err := json.Marshal(in)
	assert.NoError(t, err)

	return b
}

func userLogin(t *testing.T, username string, password string) (string, string) {
	b := serJson(t, gin.H{
		"username": username,
		"password": password,
	})
	w := reqApi(t, "POST", "/api/users/login", b, nil)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body := struct {
		Username    string `json:"username"`
		AccessToken string `json:"access_token"`
		UpdateToken string `json:"update_token"`
	}{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, body.Username, username)
	assert.NotEqual(t, body.AccessToken, "")
	assert.NotEqual(t, body.UpdateToken, "")

	return body.AccessToken, body.UpdateToken
}

func getUserDetail(t *testing.T, at string) api.UserDetail {
	w := reqApi(t, "GET", "/api/users/me", nil, &at)
	assert.Equal(t, 200, w.Code, "body: %s", w.Body.String())

	body := api.UserDetail{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)

	return body
}

func getEmailInfo(t *testing.T) (string, string) {
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
	return email, emailBody
}

func ensureNoEmail(t *testing.T) {
	select {
	case <-srv.SendEmailMockChan:
		<-srv.SendEmailMockChan
		assert.NotNil(t, nil, "No email edit but triggers email sending")
	default:
	}
}
