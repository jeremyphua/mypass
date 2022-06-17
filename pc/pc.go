package pc

import (
	"crypto/rand"
	"encoding/json"
	"log"

	"github.com/alexedwards/argon2id"
	"github.com/jeremyphua/mypass/io"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
)

var customArgon2idParams = &argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  1,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   8,
}

// Wrapper around argon2id.CreateHash to create generic function to take in string input
// Return a 32-bytes key
func Argon2id(password string) (key []byte, err error) {
	hash, err := argon2id.CreateHash(password, customArgon2idParams)
	key = []byte(hash)
	return
}

// Wrapper around secretbox.Seal
// Convert key byte slice to 32-bytes
// Randomly generate a nonce to eliminate risk of reusing nonce
func SecretboxSeal(key []byte, message []byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}
	var keyArr [32]byte
	copy(keyArr[:], key)
	return secretbox.Seal(nonce[:], message, &nonce, &keyArr), nil
}

// Wrapper around secretbox.Open
// Convert key byte slice to 32-bytes
func SecretboxOpen(key []byte, encrypted []byte) ([]byte, bool) {
	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])
	var keyArr [32]byte
	copy(keyArr[:], key)
	return secretbox.Open(nil, encrypted[24:], &decryptNonce, &keyArr)

}

// Wrapper around box.Seal to create a randomly generated nonce
func BoxSeal(message []byte, pub *[32]byte, priv *[32]byte) (out []byte, err error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}
	return box.Seal(nonce[:], message, &nonce, pub, priv), nil
}

// Retrieve master private key
func GetMasterPrivKey() []byte {
	pass, err := io.PromptPass("Please enter master password")
	if err != nil {
		log.Fatalf("Could not read password: %s", err.Error())
	}

	c, err := io.GetConfigFile()
	if err != nil {
		log.Fatalf("Could not get config file: %s", err.Error())
	}

	var configFile io.ConfigFile
	configFileBytes, err := io.util.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Could not read config file: %s", err.Error())
	}

	err = json.Unmarshal(configFileBytes, &configFile)
	if err != nil {
		log.Fatalf("Could not read unmarshal config file: %s", err.Error())
	}

	passKey, err := pc.Argon2id(pass)
	if err != nil {
		log.Fatalf("Error hashing password using Argon2id: %s", err.Error())
	}

	masterPrivKey, ok := SecretboxOpen(passKey, configFile.MasterPrivKeySealed)
	if !ok {
		log.Fatalf("Wrong master password")
	}
	return masterPrivKey
}
