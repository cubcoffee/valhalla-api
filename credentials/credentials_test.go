package credentials

import "testing"

func TestGenerateHash(t *testing.T) {
	password := "admin123"
	firtHash, _, err := GenerateHash(password)
	nextHash, _, err := GenerateHash(password)
	if err != nil {
		t.Errorf("Error generating password hash")
	}
	if firtHash == nextHash {
		t.Errorf("Password hash shouldn't match")

	}
}
