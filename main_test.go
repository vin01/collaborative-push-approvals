package main

import (
        "bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
        "collaborative-push-approvals/routes"
)

func TestProtectedRoute(t *testing.T) {
	router := routes.SetupRouter()

	w := httptest.NewRecorder()
        test_data := []byte(`{"request_id":"test"}`)
	req, _ := http.NewRequest("POST", "/request/status", bytes.NewBuffer(test_data))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginLoad(t *testing.T) {
        router := routes.SetupRouter()

        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/", nil)
        router.ServeHTTP(w, req)

        assert.Contains(t, w.Body.String(), "Log In")
}
