package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var secretKey = []byte("9^.2&f=2==))&&7we")

// Internal function to encode a string as bytes without padding
func base64Encode(src []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(src), "=")
}

// Create a new JWT token
func CreateToken(userId string) (string, error) {
	// A valid jwt token consists of three core parts seperated with (.)
	// 1. JSON Header : This contains the meta data of the json web token, like the the algorithm for the token
	// 2. Payload : This contains the user identity and other additional information like roles
	// 3. Signature : This is the signature that ensures and  verifies that the server actually signed the json web token

	// Header
	header, _ := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})

	headerEncoded := base64Encode(header)

	// Payload
	payload, _ := json.Marshal(map[string]interface{}{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	payloadEncoded := base64Encode(payload)

	// Signature
	payloadAndHeader := headerEncoded + "." + payloadEncoded
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(payloadAndHeader))
	signature := base64Encode(h.Sum(nil))

	// Return the jwt
	return payloadAndHeader + "." + signature, nil
}

// Verify token
func VerifyToken(token string) (bool, string, error) {
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return false, "", errors.New("Wrong Token Length!")
	}

	payloadAndHeader := parts[0] + "." + parts[1] // Combine payload and header
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(payloadAndHeader))
	signature := base64Encode(h.Sum(nil)) // construct signature

	if parts[2] != signature { // check if signature is same as normal
		return false, "", errors.New("Wrong Token Signature!")
	}

	payloadBytes, _ := base64.RawStdEncoding.DecodeString(parts[1]) // Header to check time
	var claims map[string]interface{}
	json.Unmarshal(payloadBytes, &claims)

	// Check if token has expired
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return false, "", errors.New("Token Expired!")
		}
	}

	// Return the payload and true if all token is verified and correct
	return true, parts[1], nil
}
