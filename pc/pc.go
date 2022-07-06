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

const (
	// password length
	pwLength = 20

	// ASCII bound values
	// https://design215.com/toolbox/ascii-utf8.php
	// Uppercase ASCII bound values
	upperCaseLowerbound = 65
	upperCaseUpperbound = 90

	// Lowercase ASCII bound values
	lowerCaseLowerbound = 97
	lowerCaseUpperbound = 122

	// Digit ASCII bound values
	digitCaseLowerbound = 48
	digitCaseUpperbound = 57

	// Symbol ASCII bound values
	symbolGroupOneLowerbound = 33
	symbolGroupOneUpperbound = 47

	symbolGroupTwoLowerbound = 58
	symbolGroupTwoUpperbound = 64

	symbolGroupThreeLowerbound = 91
	symbolGroupThreeUpperbound = 96

	symbolGroupFourLowerbound = 123
	symbolGroupFourUpperbound = 126
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

// wrapper around box.Open
func BoxOpen(encrypted []byte, pub *[32]byte, priv *[32]byte) ([]byte, bool) {
	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])
	return box.Open(nil, encrypted[24:], &decryptNonce, pub, priv)
}

// Reencrypt new password using BoxSeal
func ReEncrypt(s io.SiteInfo, password string) (io.SiteInfo, []byte) {
	var c io.ConfigFile
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Could not generate site key: %s", err.Error())
	}
	config, err := io.GetConfigFile()
	if err != nil {
		log.Fatalf("Could not get config file: %s", err.Error())
	}
	configContents, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatalf("Could not read contents of config: %s", err.Error())
	}
	err = json.Unmarshal(configContents, &c)

	masterPub := c.MasterPubKey

	passSealed, err := BoxSeal([]byte(password), &masterPub, priv)
	if err != nil {
		log.Fatalf("Could not seal new site password: %s", err.Error())
	}

	return io.SiteInfo{
		Name:     s.Name,
		PubKey:   *pub,
		Username: s.Username,
	}, passSealed

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

func GeneratePassword() (password string, err error) {
	// make a slice of random bytes
	letters := make([]byte, 10000)

	// read random bytes
	_, err = rand.Read(letters)
	if err != nil {
		return
	}

	password = ""
	for _, letter := range letters {
		// Check letter is in the range of printable characters
		if letter > 32 && letter < 127 {
			password += string(letter)
		}
		// If length of password reach 12, check if it is valid
		if len(password) == pwLength {
			if validPassword(password) {
				return
			}
			// trim left character of password
			password = password[1:]
		}
	}
	return
}

func validPassword(password string) bool {
	isUpper := false
	isLower := false
	isSymbol := false
	isDigit := false

	for i := 0; i < len(password); i++ {
		if isASCIIUpper(password[i]) {
			isUpper = true
		}
		if isASCIILower(password[i]) {
			isLower = true
		}
		if isASCIISymbol(password[i]) {
			isSymbol = true
		}
		if isASCIIDigit(password[i]) {
			isDigit = true
		}
		if isUpper && isLower && isSymbol && isDigit {
			return true
		}
	}
	return false
}

func isASCIIUpper(letter byte) bool {
	return checkBound(letter, upperCaseLowerbound, upperCaseUpperbound)
}

func isASCIILower(letter byte) bool {
	return checkBound(letter, lowerCaseLowerbound, lowerCaseUpperbound)
}

func isASCIISymbol(letter byte) bool {
	if checkBound(letter, symbolGroupOneLowerbound, symbolGroupOneUpperbound) {
		return true
	}
	if checkBound(letter, symbolGroupTwoLowerbound, symbolGroupTwoUpperbound) {
		return true
	}
	if checkBound(letter, symbolGroupThreeLowerbound, symbolGroupThreeUpperbound) {
		return true
	}
	if checkBound(letter, symbolGroupFourLowerbound, symbolGroupFourUpperbound) {
		return true
	}
	return false
}

func isASCIIDigit(letter byte) bool {
	return checkBound(letter, digitCaseLowerbound, digitCaseUpperbound)
}

func checkBound(letter byte, lowerBound, upperBound int) bool {
	if int(letter) >= lowerBound && int(letter) <= upperBound {
		return true
	}
	return false
}
