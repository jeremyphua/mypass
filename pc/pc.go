package pc

import (
	"crypto/rand"

	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/nacl/secretbox"
)

// Wrapper around argon2id.CreateHash to create generic function to take in string input
// Return a 32-bytes key
func Argon2id(password string) (key [32]byte, err error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	copy(key[:], hash)
	return
}

// Randomly generate a nonce to eliminate risk of reusing nonce
func SecretboxSeal(key *[32]byte, message []byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}
	return secretbox.Seal(nonce[:], message, &nonce, key), nil
}
