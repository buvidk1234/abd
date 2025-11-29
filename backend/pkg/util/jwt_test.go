package util

import (
	"testing"
)

func TestGenerateAndParseJWT(t *testing.T) {
	userID := "test_user_123"
	t.Log("begin")

	token, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	t.Logf("Generated Token: %s", token)

	parsedUserID, err := ParseToken(token)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}
	t.Logf("userID: %s", parsedUserID)
}
