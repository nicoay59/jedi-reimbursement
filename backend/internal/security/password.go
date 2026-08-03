package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

const (
	passwordAlgorithm  = "pbkdf2_sha256"
	passwordIterations = 210000
	passwordSaltLength = 16
	passwordKeyLength  = 32
)

func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", fmt.Errorf("password minimal 8 karakter")
	}

	salt := make([]byte, passwordSaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("membuat salt password: %w", err)
	}

	key := pbkdf2SHA256(
		[]byte(password),
		salt,
		passwordIterations,
		passwordKeyLength,
	)

	return strings.Join([]string{
		passwordAlgorithm,
		strconv.Itoa(passwordIterations),
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	}, "$"), nil
}

func VerifyPassword(encodedHash, password string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 4 || parts[0] != passwordAlgorithm {
		return false
	}

	iterations, err := strconv.Atoi(parts[1])
	if err != nil || iterations < 1 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}

	expectedKey, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil || len(expectedKey) == 0 {
		return false
	}

	actualKey := pbkdf2SHA256(
		[]byte(password),
		salt,
		iterations,
		len(expectedKey),
	)

	return subtle.ConstantTimeCompare(actualKey, expectedKey) == 1
}

func pbkdf2SHA256(
	password []byte,
	salt []byte,
	iterations int,
	keyLength int,
) []byte {
	hashLength := sha256.Size
	blockCount := (keyLength + hashLength - 1) / hashLength
	derivedKey := make([]byte, 0, blockCount*hashLength)

	for block := 1; block <= blockCount; block++ {
		mac := hmac.New(sha256.New, password)
		_, _ = mac.Write(salt)

		counter := make([]byte, 4)
		binary.BigEndian.PutUint32(counter, uint32(block))
		_, _ = mac.Write(counter)

		u := mac.Sum(nil)
		t := append([]byte(nil), u...)

		for iteration := 1; iteration < iterations; iteration++ {
			mac = hmac.New(sha256.New, password)
			_, _ = mac.Write(u)
			u = mac.Sum(nil)

			for index := range t {
				t[index] ^= u[index]
			}
		}

		derivedKey = append(derivedKey, t...)
	}

	return derivedKey[:keyLength]
}
