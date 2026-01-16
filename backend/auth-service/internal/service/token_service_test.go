package service

import (
	"auth-service/internal/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func newTestManager() (TokenManager, config.JWTConfig) {
	cfg := config.JWTConfig{
		JWTAccessTokenSecret:  "secret_access",
		JWTRefreshTokenSecret: "secret_refresh",
		JWTAccessTokenExp:     15 * time.Minute,
		JWTRefreshTokenExp:    24 * time.Hour,
	}
	return NewTokenManager(cfg), cfg
}

func TestManager_GenerateTokens(t *testing.T) {
	mgr, cfg := newTestManager()
	userID := "user-123"

	access, refresh, err := mgr.GenerateTokens(userID)

	if err != nil {
		t.Fatalf("GenerateTokens() unexpected error: %v", err)
	}

	if access == "" {
		t.Error("GenerateTokens() returned empty access token")
	}
	if refresh == "" {
		t.Error("GenerateTokens() returned empty refresh token")
	}

	if _, err := mgr.ValidateToken(access, []byte(cfg.JWTAccessTokenSecret)); err != nil {
		t.Errorf("GenerateTokens() access token invalid signature: %v", err)
	}

	if _, err := mgr.ValidateToken(refresh, []byte(cfg.JWTRefreshTokenSecret)); err != nil {
		t.Errorf("GenerateTokens() refresh token invalid signature: %v", err)
	}
}

func TestManager_GenerateAccessToken(t *testing.T) {
	mgr, cfg := newTestManager()
	userID := "user-access-only"

	token, err := mgr.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error: %v", err)
	}

	claims, err := mgr.ValidateToken(token, []byte(cfg.JWTAccessTokenSecret))
	if err != nil {
		t.Fatalf("Generated access token is invalid: %v", err)
	}

	if claims.Subject != userID {
		t.Errorf("GenerateAccessToken() subject = %v, want %v", claims.Subject, userID)
	}
}

func TestManager_GenerateRefreshToken(t *testing.T) {
	mgr, cfg := newTestManager()
	userID := "user-refresh-only"

	token, err := mgr.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error: %v", err)
	}

	claims, err := mgr.ValidateToken(token, []byte(cfg.JWTRefreshTokenSecret))
	if err != nil {
		t.Fatalf("Generated refresh token is invalid: %v", err)
	}

	if claims.Subject != userID {
		t.Errorf("GenerateRefreshToken() subject = %v, want %v", claims.Subject, userID)
	}
}

func TestManager_ValidateRefreshToken(t *testing.T) {
	mgr, cfg := newTestManager()
	userID := "user-validate-refresh"

	validRefresh, _ := mgr.GenerateToken(userID, time.Hour, []byte(cfg.JWTRefreshTokenSecret))
	accessToken, _ := mgr.GenerateToken(userID, time.Hour, []byte(cfg.JWTAccessTokenSecret))

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Success: Valid Refresh Token",
			token:   validRefresh,
			wantErr: false,
		},
		{
			name:    "Fail: Access Token passed as Refresh (Wrong Signature)",
			token:   accessToken,
			wantErr: true,
		},
		{
			name:    "Fail: Garbage string",
			token:   "invalid.token.structure",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := mgr.ValidateRefreshToken(tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if claims == nil || claims.Subject != userID {
					t.Errorf("ValidateRefreshToken() wrong claims returned")
				}
			}
		})
	}
}

func TestManager_GenerateToken(t *testing.T) {
	mgr, _ := newTestManager()
	customSecret := []byte("custom-secret")
	userID := "user-custom"

	ttl := time.Second * 5
	tokenStr, err := mgr.GenerateToken(userID, ttl, customSecret)
	if err != nil {
		t.Fatalf("GenerateToken() error: %v", err)
	}

	token, _ := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return customSecret, nil
	})

	if !token.Valid {
		t.Error("GenerateToken() produced invalid token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.Subject != userID {
		t.Errorf("GenerateToken() wrong subject: got %v, want %v", claims.Subject, userID)
	}
}

func TestManager_ValidateToken(t *testing.T) {
	mgr, _ := newTestManager()
	secret := []byte("test-secret")
	userID := "user-validate"

	validToken, _ := mgr.GenerateToken(userID, time.Hour, secret)

	expiredToken, _ := mgr.GenerateToken(userID, -time.Minute, secret)

	wrongSecretToken, _ := mgr.GenerateToken(userID, time.Hour, []byte("wrong-secret"))

	tests := []struct {
		name       string
		tokenInput string
		secret     []byte
		wantErr    bool
	}{
		{
			name:       "Success: Valid Token",
			tokenInput: validToken,
			secret:     secret,
			wantErr:    false,
		},
		{
			name:       "Fail: Expired Token",
			tokenInput: expiredToken,
			secret:     secret,
			wantErr:    true,
		},
		{
			name:       "Fail: Wrong Secret",
			tokenInput: wrongSecretToken,
			secret:     secret,
			wantErr:    true,
		},
		{
			name:       "Fail: Malformed String",
			tokenInput: "not.a.jwt",
			secret:     secret,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := mgr.ValidateToken(tt.tokenInput, tt.secret)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if claims.Subject != userID {
					t.Errorf("ValidateToken() subject = %v, want %v", claims.Subject, userID)
				}
			}
		})
	}
}
