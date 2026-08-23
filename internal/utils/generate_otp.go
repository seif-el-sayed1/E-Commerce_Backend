package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func GenerateOTPCode() (rawCode string, hashedCode string, err error) {
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", "", err
	}
	rawCode = fmt.Sprintf("%06d", n.Int64())
	hash := sha256.Sum256([]byte(rawCode))
	hashedCode = hex.EncodeToString(hash[:])
	return rawCode, hashedCode, nil
}
