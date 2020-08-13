package credentials

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type Credential struct {
	ID   int64
	Hash string
	Salt string
}

func GenerateHash(password string) (string, string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", "", err
	}
	passwordHash, saltHash := generateSaltedHash(password, salt)
	return passwordHash, saltHash, err
}

func generateSaltedHash(password string, salt []byte) (string, string) {
	h := sha256.New()
	h.Write(append(salt, []byte(password)...))
	return fmt.Sprintf("%x", h.Sum(nil)), base64.StdEncoding.EncodeToString(salt)
}
