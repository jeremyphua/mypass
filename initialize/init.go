package initialize

import (
	"fmt"
	"log"
	"os"

	"github.com/alexedwards/argon2id"
	"github.com/jeremyphua/mypass/io"
)

var (
	needsDir      bool
	hasSiteFile   bool
	hasMasterPass bool
	hasVault      bool
)

// Initialize a new password vault in the home directory and their respective folders
// What will be initialize:
// 1. application dir -> C:\Users\<name of user>\.mypass
// 2. masterpassword file -> C:\Users\<name of user>\.mypass\masterpass
// 3. sites file -> C:\Users\<name of user>\.mypass\sites.json
// 4. vault folder -> C:\Users\<name of user>\.mypass\vault
func Init() {

	checkDirAndFoldersExists()

	// check if application dir valid
	passDir, err := io.GetPassDir()
	if err != nil {
		log.Fatalf("Could not get pass dir: %s", err.Error())
	}

	// check if password vault dir valid
	siteFile, err := io.GetSiteFile()
	if err != nil {
		log.Fatalf("Could not get site file: %s", err.Error())
	}

	// check if master password dir valid
	masterPass, err := io.GetMasterPassword()
	if err != nil {
		log.Fatalf("Could not get master password: %s", err.Error())
	}

	// check if vault dir valid
	vault, err := io.GetVaultFolder()
	if err != nil {
		log.Fatalf("Could not get vault: %s", err.Error())
	}

	/*
		prompt for master password to allow user to run init the second time
		if they quits before password vault is fully initialized
		If this step is placed after creation of folders, there will be error when initializing folders
		if user quits during password prompt
	*/
	pass, err := io.PromptPass("Please enter your password")
	if err != nil {
		log.Fatalf("Could not read password: %s", err.Error())
	}

	// if password vault does not exist, create folder C:\Users\<name of user>\.mypass
	if needsDir {
		CreateAppDir(passDir)
	}

	// Don't accidentally delete master password file or any file with similar name
	if hasMasterPass {
		log.Fatalf("Master password file already found")
	}

	CreateMasterpassFile(masterPass)

	// Create and initialize the sites.json
	if !hasSiteFile {
		CreateSiteFile(siteFile)
	}

	// Create vault folder if not exist
	if !hasVault {
		CreateVaultFolder(vault)
	}

	// Create hash of master password
	passKey, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		log.Fatalf("Could not hash master password: %s", err.Error())
	}

	// Save master password to file
	if err = io.SaveFile(passKey); err != nil {
		log.Fatalf("Could not write to master password file: %s", err.Error())
	} else {
		fmt.Printf("Successfully written hashed password to masterpass file\n")
	}

	fmt.Println("Password Vault successfully initialized")
}

// Check if dir and respective folders exist and update respective variables
func checkDirAndFoldersExists() {
	if dirExists, err := io.PassDirExists(); err == nil { // Check application dir exists
		if !dirExists {
			needsDir = true
		} else {
			if _, err := io.MasterPasswordExists(); err == nil { // Check masterpassword file exists
				hasMasterPass = true
			}
			if _, err := io.SiteFileExists(); err == nil { // Check site file exists
				hasSiteFile = true
			}
			if _, err := io.VaultExists(); err == nil { // Check vault folder exists
				hasVault = true
			}
		}
	} else {
		fmt.Println(err.Error())
	}

}

func CreateAppDir(passDir string) {
	err := os.Mkdir(passDir, 0700)
	if err != nil {
		log.Fatalf("Could not create mypass vault: %s", err.Error())
	} else {
		fmt.Printf("Successfully created directory to store passwords at: %s\n", passDir)
	}
}

// Create master password file with secure permission.
// os.Create() leaves file world-readable.
func CreateMasterpassFile(masterPass string) {
	config, err := os.OpenFile(masterPass, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Could not create passgo config: %s", err.Error())
	}
	config.Close()
	fmt.Printf("Successfully created masterpass file to store encrypted master password at: %s\n", masterPass)
}

// Create file, with secure permissions.
func CreateSiteFile(siteFile string) {
	sf, err := os.OpenFile(siteFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Could not create site file: %s", err.Error())
	}
	// Initialize an empty SiteFile
	siteFileContents := []byte("[]")
	_, err = sf.Write(siteFileContents)
	if err != nil {
		log.Fatalf("Could not save site file: %s", err.Error())
	}
	sf.Close()
	fmt.Printf("Successfully created site file to store information at: %s\n", siteFile)
}

func CreateVaultFolder(vault string) {
	err := os.Mkdir(vault, 0700)
	if err != nil {
		log.Fatalf("Could not create vault folder: %s", err.Error())
	} else {
		fmt.Printf("Successfully created directory to store encrypted passwords at: %s\n", vault)
	}
}
