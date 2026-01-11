package util

import (
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
