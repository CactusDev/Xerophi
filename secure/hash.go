package secure

import (
	"crypto/rand"
	"encoding/hex"
	"reflect"

	"golang.org/x/crypto/argon2"

	log "github.com/sirupsen/logrus"
)

// Functions/data related to hashing/verifying passwords

// GenSalt Generates a crypto-secure salt
func GenSalt() []byte {
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

// HashPassword hashes the password with the configured salt
func HashPassword(password string) (string, string) {
	salt := GenSalt()
	hashed := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return hex.EncodeToString(hashed), hex.EncodeToString(salt)
}

// VerifyHash verifies the provided hash against the provided password
func VerifyHash(hash string, password string, salt string) (bool, error) {
	// Decode the hex-encoded strings (likely from the DB)
	decodedHash, err := hex.DecodeString(hash)
	if err != nil {
		log.Error(err)
		return false, err
	}
	decodedSalt, err := hex.DecodeString(salt)
	if err != nil {
		log.Error(err)
		return false, err
	}

	// Create a hash using the password provided and the decoded salt
	hashed := argon2.IDKey([]byte(password), decodedSalt, 1, 64*1024, 4, 32)

	// Compare the decoded hash and the newly hashed string, they should be the
	// same
	return reflect.DeepEqual(decodedHash, hashed), nil
}
