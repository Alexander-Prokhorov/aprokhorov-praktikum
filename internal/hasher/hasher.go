package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func HashHMAC(src string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(src))
	return fmt.Sprintf("%x", h.Sum(nil))
}
