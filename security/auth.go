package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	time           = 3         //Iterations
	memory         = 64 * 1024 //Memory
	threads        = 2         //Max threads (or parallel)
	keyLen         = 32
	expectedLength = 4 //Length of argon2 after splitting
)

var (
	ErrWrongPass   = errors.New("password is wrong")
	ErrWrongFormat = errors.New("argon2 is formatted wrong")
)

type Info struct {
	Pw   string
	Salt string
}

func NewInfoStruct() *Info {
	return &Info{}
}

func generateRandomBytes(l *log.Logger) ([]byte, error) {
	l.Printf("[INFO] generateRandomBytes called")
	bytes := make([]byte, 26)
	_, err := rand.Read(bytes)

	if err != nil {
		l.Printf("[ERROR] Failed to insert random bytes in function generateRandomBytes. Error %s", err)
		return nil, err
	}
	l.Printf("[INFO] Successfully generated random bytes")
	return bytes, nil
}

func generateArgon2(password string, l *log.Logger) (string, string, error) {
	l.Printf("[INFO] generateArgon2 called")

	salt, err := generateRandomBytes(l)
	if err != nil {
		return "", "", err
	}
	l.Printf("[INFO] generating and encoding argon2")
	hashedPass := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	b64Pass := base64.RawStdEncoding.EncodeToString(hashedPass)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)

	a2S := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%s", argon2.Version, time, threads, b64Pass)
	l.Printf("[INFO] argon2 generated")

	return a2S, b64Salt, nil

}
func VerifyPass(original string, hashed *Info, l *log.Logger) error {
	l.Printf("[INFO] VerifyPass called")

	l.Printf("[INFO] Extracting password")
	pass, err := extractPass(hashed.Pw, l)

	if err != nil {
		l.Printf("[ERROR] Failed to extract password. Error %s", err)
		return err
	}

	l.Printf("[INFO] Decoding base64")
	salt, err := base64.RawStdEncoding.DecodeString(hashed.Salt)
	if err != nil {
		l.Printf("[ERROR] Failed to decode base64. Error %s", err)
		return err
	}
	l.Printf("Hashing and encoding original")
	originalHashed := argon2.IDKey([]byte(original), salt, time, memory, threads, keyLen)
	originalEncoded := base64.RawStdEncoding.EncodeToString(originalHashed)
	l.Printf("[INFO] Checking for equality")
	equals := subtle.ConstantTimeCompare([]byte(originalEncoded), []byte(pass)) == 1

	if !equals {
		l.Printf("[INFO] Detected wrong password. Returning")
		return ErrWrongPass
	}

	return nil

}

func extractPass(arg2String string, l *log.Logger) (string, error) {
	l.Printf("[INFO] function extractPass ecalled")
	arg2Arr := strings.Split(arg2String, "$")

	if len(arg2Arr) != expectedLength {
		l.Printf("[ERROR] argon2 is not formatted correctly. Returning ErrWrongFormat")
		return "", ErrWrongFormat
	}
	arg2Arr = strings.Split(arg2Arr[3], ",")
	var pw string
	_, err := fmt.Sscanf(arg2Arr[2], "p=%s", &pw)
	if err != nil {
		l.Printf("[ERROR] Failed to scan string. Error %s", err)
		return "", err
	}
	return pw, nil
}
