package initialize

import (
	"fmt"
	"log"
	"os"

	"github.com/alexedwards/argon2id"
	"github.com/jeremyphua/mypass/io"
)

// Initialize a new password vault in the home directory
func Init() {
	var needsDir bool
	var hasVault bool

	// check if password vault directory exists
	if dirExists, err := io.PassDirExists(); err == nil {
		if !dirExists {
			needsDir = true
		} else {
			if _, err := io.PassConfigExists(); err == nil {
				hasVault = true
			}
		}
	}

	// check if password dir is a valid path
	passDir, err := io.GetPassDir()
	if err != nil {
		log.Fatalf("Could not get pass dir: %s", err.Error())
	}

	// check if password vault exist
	vaultFile, err := io.GetVaultFile()
	if err != nil {
		log.Fatalf("Could not get vault file: %s", err.Error())
	}

	// check if config
	configFile, err := io.GetConfigPath()
	if err != nil {
		log.Fatalf("Could not get pass config: %s", err.Error())
	}

	/*
		prompt for master password to allow user to run init the second time
		if they quits before password vault is fully initialized
	*/
	pass, err := io.PromptPass("Please enter your password")
	if err != nil {
		log.Fatalf("Could not read password: %s", err.Error())
	}
	fmt.Println(pass)

	// if password vault does not exist, create
	if needsDir {
		err = os.Mkdir(passDir, 0700)
		if err != nil {
			log.Fatalf("Could not create mypass vault: %s", err.Error())
		} else {
			fmt.Printf("Successfully created directory to store passwords at: %s\n", passDir)
		}
	}

	// Don't just go around deleting things for users or prompting them
	// to delete things. Make them do this manaully. Maybe this saves 1
	// person an afternoon.
	if hasVault {
		log.Fatalf("Config file already found")
	}

	// Create file with secure permission.  os.Create() leaves file world-readable.
	config, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Could not create passgo config: %s", err.Error())
	}
	config.Close()

	// Create and initialize the site vault
	if !hasVault {
		// Create file, with secure permissions.
		vf, err := os.OpenFile(vaultFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatalf("Could not create pass vault: %s", err.Error())
		}
		// Initialize an empty SiteFile
		vaultFileContents := []byte("[]")
		_, err = vf.Write(vaultFileContents)
		if err != nil {
			log.Fatalf("Could not save site file: %s", err.Error())
		}
		vf.Close()
	}

	// KDF master password
	passKey, err := argon2id.CreateHash(pass, argon2id.DefaultParams)

	if err = io.SaveFile(passKey); err != nil {
		log.Fatalf("Could not write to config file: %s", err.Error())
	}

	fmt.Println("Password Vault successfully initialized")
}
