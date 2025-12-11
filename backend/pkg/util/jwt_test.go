package util

import (
	"testing"
)

func TestGenerateAndParseJWT(t *testing.T) {
	userID := int64(123456)
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
	t.Logf("userID: %d", parsedUserID)
}
