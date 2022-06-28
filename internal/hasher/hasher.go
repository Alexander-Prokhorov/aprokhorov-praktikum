package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
)

func HashHMAC(src string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(src))
	return string(h.Sum(nil))
}
