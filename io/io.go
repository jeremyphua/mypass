package io

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	VaultFileName  = "vault.json"
	ConfigFileName = "config"
)

// return error if password vault directory does not exist
func PassDirExists() (bool, error) {
	d, err := GetPassDir()
	if err != nil {
		return false, err
	}

	dirInfo, err := os.Stat(d)
	if err == nil && dirInfo.IsDir() {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Returns dir of password vault
func GetPassDir() (d string, err error) {
	home, err := getHomeDir()
	if err == nil {
		d = filepath.Join(home, ".mypass")
		return
	}
	return
}

// Returns user home dir, Example: C:\Users\<name of user>
func getHomeDir() (d string, err error) {
	d, err = os.UserHomeDir()
	return
}

// Returns password vault
func GetVaultFile() (d string, err error) {
	p, err := GetPassDir()
	if err == nil {
		d = filepath.Join(p, VaultFileName)
	}
	return
}

// PassConfigExists is used to determine if the passgo config
// file exists in the user's passgo directory.
func PassConfigExists() (bool, error) {
	c, err := GetConfigPath()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(c)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetConfigPath is used to get the user's passgo directory.
func GetConfigPath() (p string, err error) {
	d, err := GetPassDir()
	if err == nil {
		p = filepath.Join(d, ConfigFileName)
	}
	return
}

func SaveFile(s string) (err error) {
	if exists, err := PassConfigExists(); err != nil {
		log.Fatalf("Could not find config file: %s", err.Error())
	} else if !exists {
		log.Fatalf("pass config could not be foound %s", err.Error())
	}
	sBytes, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		log.Fatalf("Could not marshal config file: %s", err.Error())
	}
	path, err := GetConfigPath()
	if err != nil {
		log.Fatalf("Could not get config file path: %s", err.Error())
	}
	err = ioutil.WriteFile(path, sBytes, 0666)
	return
}

func PromptPass(prompt string) (pass string, err error) {
	fd := int(os.Stdin.Fd())
	fmt.Printf("%s: ", prompt)
	passBytes, err := terminal.ReadPassword(fd)
	fmt.Println("")
	return string(passBytes), err
}
