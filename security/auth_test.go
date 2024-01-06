package security

import (
	"log"
	"testing"
)

func TestVerifyCorrectPassword(t *testing.T) {
	plain := "He11o Wor!d"
	hashedInfo, _ := generateArgon2(plain, log.Default())
	err := VerifyPass(plain, hashedInfo, log.Default())
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

}

func TestVerifyWrongPassword(t *testing.T) {
	plain := "He11o Wor!d"
	hashedInfo, _ := generateArgon2(plain, log.Default())
	err := VerifyPass("Hello world!", hashedInfo, log.Default())
	if err == nil {
		t.Fatalf("Expected to get %s but failed to get an error", ErrWrongPass)
	}
	if err != ErrWrongPass {
		t.Fatalf("Expected to get %s, got %s", ErrWrongPass, err)
	}
}

func TestVerifyWrongFormat(t *testing.T) {
	pass := "$v=19$m=3,t=2,p=1$bzZIcEFZVloyRlJUU3pXMA$AZwYaC/9RRssd/9X2SZ9Lg"

	//Call generateArgon2 to generate a b64 salt
	info, _ := generateArgon2("Password", log.Default())

	//Change password to wrong format
	info.Pw = pass

	err := VerifyPass("Password", info, log.Default())
	if err == nil {
		t.Fatalf("Expected to get %s but failed to get an error", ErrWrongFormat)
	}
	if err != ErrWrongFormat {
		t.Fatalf("Expected to get %s, got %s", ErrWrongFormat, ErrWrongPass)
	}

}

func TestVerifySpacedPassword(t *testing.T) {
	plain := "Password"
	hashedInfo, _ := generateArgon2(plain, log.Default())
	err := VerifyPass("Password ", hashedInfo, log.Default())
	if err == nil {
		t.Fatalf("Expected to get %s but failed to get an error", ErrWrongPass)
	}
	if err != ErrWrongPass {
		t.Fatalf("Expected to get %s, got %s", ErrWrongPass, err)
	}
}
