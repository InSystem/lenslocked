package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RemeberTokenBytes = 32

// Bytes help ug generate n random bytes.
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//NBytes return a number of bytes used in the base64
//URL encoding string
func NBytes(base64String string) (int,error){
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil{
		return -1, err
	}
	return len(b), nil
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
