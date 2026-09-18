package token

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const testSecret = "test_secret_for_unit_tests"

func TestGenerateToken_Success(t *testing.T) {
	userID := uuid.New()
	secret := []byte(testSecret)

	tokenString, err := GenerateToken(userID, "dm", secret, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	if tokenString == "" {
		t.Error("token string is empty")
	}
}

func TestValidateToken_Success(t *testing.T) {
	userID := uuid.New()
	secret := []byte(testSecret)
	tokenString, err := GenerateToken(userID, "dm", secret, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	claims, err := ValidateToken(tokenString, secret)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID: got %v want %v", claims.UserID, userID)
	}
	if claims.Role != "dm" {
		t.Errorf("Role: got %q want %q", claims.Role, "dm")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := []byte(testSecret)
	otherSecret := []byte("different_secret")

	tokenString, err := GenerateToken(userID, "dm", secret, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = ValidateToken(tokenString, otherSecret)
	if !errors.Is(err, jwt.ErrSignatureInvalid) {
		t.Errorf("Signature invalid: got %v want %v", err, jwt.ErrSignatureInvalid)
	}
}

func TestValidateToken_Expired(t *testing.T) {
	userID := uuid.New()
	secret := []byte(testSecret)

	tokenString, err := GenerateToken(userID, "dm", secret, -1*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = ValidateToken(tokenString, secret)
	if err == nil {
		t.Fatal("ValidateToken should have failed with expired token")
	}
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("expected ErrTokenExpired, got: %v", err)
	}
}

func TestValidateToken_Garbage(t *testing.T) {
	secret := []byte(testSecret)

	_, err := ValidateToken("not.a.jwt", secret)
	if err == nil {
		t.Fatalf("expected error for garbage input, got  nil")
	}
	if !errors.Is(err, jwt.ErrTokenMalformed) {
		t.Errorf("expected ErrTokenMalformed, got  %v", err)
	}
}
