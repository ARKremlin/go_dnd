package hasher

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	if hash == "" {
		t.Error("HashPassword return empty string")
	}
}
func TestHashPassword_DifferentForSameInput(t *testing.T) {
	hash1, err1 := HashPassword("secret")
	hash2, err2 := HashPassword("secret")
	if err1 != nil || err2 != nil {
		t.Fatalf("HashPassword return error: %v, %v", err1, err2)
	}
	if hash1 == hash2 {
		t.Error("hash1 == hash2: salt didn't work")
	}
}

func TestCheckPassword_Correct(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatalf("HashPassword return error: %v", err)
	}
	ok := CheckPassword(hash, "secret")
	if !ok {
		t.Error("CheckPassword returned false")
	}
}

func TestCheckPassword_Wrong(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatalf("HashPassword return error: %v", err)
	}
	ok := CheckPassword(hash, "other")
	if ok {
		t.Error("ожидалось false, получено true")
	}
}

func TestCheckPassword_InvalidHash(t *testing.T) {
	ok := CheckPassword("not_a_hash", "anything")
	if ok {
		t.Error("want false")
	}
}
