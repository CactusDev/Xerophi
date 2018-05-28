package secure

import (
	"crypto/rand"
	"reflect"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/argon2"
)

// Functions/data related to hashing/verifying passwords

// Generates a crypto-secure salt
func genSalt() []byte {
	// Make the salt 64 bytes long
	salt := make([]byte, 64)
	// Read the random salt from /dev/urandom
	_, err := rand.Read(salt)
	// If there's an error we don't want to continue
	if err != nil {
		log.Error(err)
		return nil
	}
	// Didn't error, return the salt
	return salt
}

// HashArgon hashes the password with the configured salt
func HashArgon(password string) []byte {
	return argon2.IDKey([]byte(password), genSalt(), 1, 64*1024, 4, 32)
}

// VerifyHash verifies the provided hash against the provided password
func VerifyHash(hash []byte, password string) bool {
	return reflect.DeepEqual(hash, HashArgon(password))
}
