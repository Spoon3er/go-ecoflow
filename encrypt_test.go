package ecoflow

import (
	"fmt"
	"testing"
)

func TestEncryptHmacSHA256(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		message  string
		secret   string
		expected string
	}{
		{
			name:    "Basic test",
			message: "accessKey=sy5dmGY3BiRYJnW1aLtaY3UkfBf4eONx&nonce=825095&timestamp=1762031611763000000",
			secret:  "IAVxCanSqPBVpjATS9lAW4Zot0zSCAPf",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := encryptHmacSHA256(tc.message, tc.secret)
			fmt.Printf("HMAC-SHA256('%s', '%s') = %s\n", tc.message, tc.secret, actual)
		})
	}
}
