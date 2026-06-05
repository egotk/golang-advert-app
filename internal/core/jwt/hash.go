package corejwt

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	hash := base64.RawURLEncoding.EncodeToString(sum[:])

	return hash
}
