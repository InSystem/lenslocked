package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// HMAC is a wrapper around package "crypto/hmac" that
// allows us to use it a little bit easer in our code.
type HMAC struct {
	hmac hash.Hash
}

// NewHMAC creates and returns a new hmac object.
func NewHMAC(key string) HMAC {
	h := hmac.New(sha256.New, []byte(key))
	return HMAC{
		hmac: h,
	}
}

// Hash will hash the provided input string using HMAC with the secret key
// provided when the HMAC was created.
func (h HMAC) Hash(input string) string {
	h.hmac.Reset()
	// nolint
	h.hmac.Write([]byte(input))
	b := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}
