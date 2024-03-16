package tools

import (
	"crypto"
	"encoding/hex"
)

func HashString(toHash string) string{
	hasher := crypto.SHA256.New()
	hasher.Write([]byte(toHash))
	hashedBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}