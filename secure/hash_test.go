package secure_test

import (
	"testing"

	"github.com/CactusDev/Xerophi/secure"
)

func TestVerifyHash(t *testing.T) {
	// Yay simple two-case test (for now)!
	goodPassword := "test"
	goodHash, goodSalt := secure.HashPassword(goodPassword)

	verify, err := secure.VerifyHash(goodHash, goodPassword, goodSalt)
	if err != nil {
		t.Error(err)
	}
	if !verify {
		t.Error("Failed to verify known good hash")
	}

	badPassword := "test"
	badHash, badSalt := secure.HashPassword(badPassword)

	verify, err = secure.VerifyHash(badHash, "notTest", badSalt)
	if err != nil {
		t.Error(err)
	}
	if verify {
		t.Error("Verified known good hash")
	}
}
