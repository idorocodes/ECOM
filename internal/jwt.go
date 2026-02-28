package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var secretKey = []byte("9^.2&f=2==))&&7we")

// Internal function to encode a string as bytes without padding
func base64Encode(src []byte) string {
	return strings.TrimRight(base64.RawURLEncoding.EncodeToString(src), "=")
}

// Create a new JWT token
func CreateToken(userId string, role string) (string, error) {
	// A valid jwt token consists of three core parts seperated with (.)
	// 1. JSON Header : This contains the meta data of the json web token, like the the algorithm for the token
	// 2. Payload : This contains the user identity and other additional information like roles
	// 3. Signature : This is the signature that ensures and  verifies that the server actually signed the json web token

	// Header
	header, _ := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})

	headerEncoded := base64Encode(header)

	// Payload
	payload, _ := json.Marshal(map[string]interface{}{
		"sub":  userId,
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Minute * 30).Unix(),
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
func VerifyToken(token string) (bool, string,string, error) {
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return false, "","", errors.New("Wrong Token Length!")
	}

	payloadAndHeader := parts[0] + "." + parts[1] // Combine payload and header
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(payloadAndHeader))
	signature := base64Encode(h.Sum(nil)) // construct signature

	if !hmac.Equal([]byte(parts[2]), []byte(signature)) {
		fmt.Printf("DEBUG: Expected %s, got %s\n", signature, parts[2])
		return false, "","", errors.New("invalid signature")
	}

	payloadBytes, _ := base64.RawURLEncoding.DecodeString(parts[1]) // Header to check time
	var claims map[string]interface{}
	json.Unmarshal(payloadBytes, &claims)

	headerBytes, _ := base64.RawURLEncoding.DecodeString(parts[0])
	var header map[string]interface{}
	json.Unmarshal(headerBytes, &header)

	//Check the algorithm
	if header["alg"] != "HS256" {
		return false, "","", errors.New("unsupported signing algorithm")
	}

	userId, ok := claims["sub"].(string)
	if len(userId) == 0 || !ok {
		return false, "","", errors.New("missing subject claim")
	}

	role, ok := claims["role"].(string)
	if len(role) == 0 || !ok {
		return false, "","", errors.New("missing role")
	}

	if iat, ok := claims["iat"].(float64); ok {
		if time.Now().Unix() < int64(iat) {
			return false, "","", errors.New("token used before issued")
		}
	}

	// Check if token has expired
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return false, "","", errors.New("Token Expired!")
		}
	}

	// Return the payload and true if all token is verified and correct
	return true, userId,role, nil
}

func HashPassowrd(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))

}
