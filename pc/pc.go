package pc

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func BoxOpen(encrypted []byte, pub *[32]byte, priv *[32]byte) ([]byte, bool) {
	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])
	return box.Open(nil, encrypted[24:], &decryptNonce, pub, priv)
}

// Retrieve master private key
func GetMasterPrivKey() (masterPrivKey [32]byte) {
	pass, err := io.PromptPass("Please enter master password")
	if err != nil {
		log.Fatalf("Could not read password: %s", err.Error())
	}

	cfile, err := io.GetConfigFile()
	if err != nil {
		log.Fatalf("Could not get config file: %s", err.Error())
	}

	var c io.ConfigFile

	configFileBytes, err := ioutil.ReadFile(cfile)
	if err != nil {
		log.Fatalf("Could not read config file: %s", err.Error())
	}

	err = json.Unmarshal(configFileBytes, &c)
	if err != nil {
		log.Fatalf("Could not read unmarshal config file: %s", err.Error())
	}

	validateMasterPassword(pass, string(c.MasterPassKey))

	masterPrivKeySlice, ok := SecretboxOpen(c.MasterPassKey, c.MasterPrivKeySealed)

	if !ok {
		log.Fatalf("Failed to get master private key")
	}

	copy(masterPrivKey[:], masterPrivKeySlice)
	fmt.Println("Authentication success!")
	return
}

func validateMasterPassword(input string, encryptedMasterPassword string) {
	match, err := argon2id.ComparePasswordAndHash(input, encryptedMasterPassword)
	if err != nil {
		log.Fatalf("Error comparing password: %s", err.Error())
	}
	if !match {
		log.Fatalf("Wrong master password")
	}
}
