package webhook

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

// VerifySignature checks if the provided signature matches the expected HMAC-SHA512 hash
// of the payload and timestamp using the secret.
func VerifySignature(payload []byte, signature string, timestamp string, secret string) bool {
	if secret == "" {
		return false
	}

	// 1. Concatenate payload and timestamp
	valueToDigest := string(payload) + timestamp

	// 2. Compute HMAC-SHA512
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(valueToDigest))
	expectedMAC := mac.Sum(nil)
	// expectedSignature := hex.EncodeToString(expectedMAC)

	// 3. Compare signatures
	// Using hmac.Equal would be better for constant time comparison if we had bytes,
	// but here we are comparing hex strings.
	// Let's decode the hex string to bytes to use hmac.Equal for security.
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	return hmac.Equal(sigBytes, expectedMAC)
}
