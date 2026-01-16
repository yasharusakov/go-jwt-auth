package util

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSetRefreshTokenCookie(t *testing.T) {
	recorder := httptest.NewRecorder()
	token := "test_refresh_token"

	SetRefreshTokenCookie(recorder, token, 24*time.Hour, true)

	res := recorder.Result()
	cookies := res.Cookies()

	if len(cookies) == 0 {
		t.Fatal("No cookies set")
	}

	cookie := cookies[0]

	if cookie.Name != "refresh_token" {
		t.Errorf("Wrong cookie name: got %s, want refresh_token", cookie.Name)
	}
	if cookie.Value != token {
		t.Errorf("Wrong cookie value: got %s, want %s", cookie.Value, token)
	}
	if !cookie.HttpOnly {
		t.Error("Cookie should be HttpOnly")
	}
	if cookie.Path != "/" {
		t.Error("Cookie path should be /")
	}
	if !cookie.Secure {
		t.Error("Cookie should be Secure in production")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Error("Cookie SameSite should be Lax")
	}
}

func TestRemoveRefreshTokenCookie(t *testing.T) {
	recorder := httptest.NewRecorder()

	RemoveRefreshTokenCookie(recorder, true)

	res := recorder.Result()
	cookies := res.Cookies()

	if len(cookies) == 0 {
		t.Fatal("No cookies set for removal")
	}

	cookie := cookies[0]

	if cookie.Value != "" {
		t.Error("Cookie value should be empty")
	}
	if cookie.MaxAge >= 0 {
		t.Error("Cookie MaxAge should be negative")
	}
	if !cookie.Expires.Before(time.Now()) {
		t.Error("Cookie Expires should be in the past")
	}
}
