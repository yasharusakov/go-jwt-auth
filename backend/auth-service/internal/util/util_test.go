package util

import (
	"auth-service/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	secret := []byte("secret")
	userID := "user123"

	validToken, err := GenerateToken(userID, 1*time.Hour, secret)
	if err != nil {
		t.Fatalf("Failed to generate valid token: %v", err)
	}

	expiredToken, _ := GenerateToken(userID, -1*time.Minute, secret)

	wrongSecretToken, _ := GenerateToken(userID, 1*time.Hour, []byte("wrongsecret"))

	tests := []struct {
		name       string
		tokenInput string
		secret     []byte
		wantErr    bool
	}{
		{"Success: Valid Token", validToken, secret, false},
		{"Fail: Expired Token", expiredToken, secret, true},
		{"Fail: Wrong Secret Signature", wrongSecretToken, secret, true},
		{"Fail: Garbage String", "not.a.real.token", secret, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.tokenInput, tt.secret)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if claims.Subject != userID {
					t.Errorf("ValidateToken() got UserID = %v, want %v", claims.Subject, userID)
				}
			}
		})
	}
}

func TestGenerateTokens(t *testing.T) {
	mockJWTCfg := config.JWTConfig{
		JWTAccessTokenSecret:  "access_secret",
		JWTRefreshTokenSecret: "refresh_secret",
		JWTAccessTokenExp:     15 * time.Minute,
		JWTRefreshTokenExp:    24 * time.Hour,
	}

	userID := "user123"

	accessToken, refreshToken, err := GenerateTokens(userID, mockJWTCfg)
	if err != nil {
		t.Fatalf("GenerateTokens() error = %v", err)
	}

	if accessToken == "" || refreshToken == "" {
		t.Errorf("GenerateTokens() accessToken = %v, refreshToken = %v; want non-empty tokens", accessToken, refreshToken)
	}

	claims, err := ValidateToken(accessToken, []byte(mockJWTCfg.JWTAccessTokenSecret))
	if err != nil {
		t.Errorf("AccessToken validation failed: %v", err)
	}
	if claims.Subject != userID {
		t.Errorf("AccessToken wrong subject: got %s, want %s", claims.Subject, userID)
	}

	claimsRef, err := ValidateToken(refreshToken, []byte(mockJWTCfg.JWTRefreshTokenSecret))
	if err != nil {
		t.Errorf("RefreshToken validation failed: %v", err)
	}
	if claimsRef.Subject != userID {
		t.Errorf("RefreshToken wrong subject: got %s, want %s", claimsRef.Subject, userID)
	}
}

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
