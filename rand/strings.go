package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RemeberTokenBytes = 32

// Bytes help us generate n random bytes.
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String will generate a byite slice of size nBytes
// and then return a string that is base64 URL encoded version
// of that byte slise.
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken generate aa remember token of predeterminded byte size.
func RememberToken() (string, error) {
	return String(RemeberTokenBytes)
}
