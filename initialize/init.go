package initialize

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"

	"github.com/jeremyphua/mypass/io"
	"github.com/jeremyphua/mypass/pc"
	"golang.org/x/crypto/nacl/box"
)

var (
	needsDir      bool
	hasSiteFile   bool
	hasConfigFile bool
	hasVault      bool
)

// Initialize a new password vault in the home directory and their respective folders
// What will be initialize:
// 1. application dir -> C:\Users\<name of user>\.mypass
// 2. config file -> C:\Users\<name of user>\.mypass\masterpass
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

	// check if config file dir valid
	configFile, err := io.GetConfigFile()
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
	if hasConfigFile {
		log.Fatalf("Master password file already found")
	}

	CreateConfigFile(configFile)

	// Create and initialize the sites.json
	if !hasSiteFile {
		CreateSiteFile(siteFile)
	}

	// Create vault folder if not exist
	if !hasVault {
		CreateVaultFolder(vault)
	}

	// kdf the master password
	passKey, err := pc.Argon2id(pass)
	if err != nil {
		log.Fatalf("Error hashing password using Argon2id: %s", err.Error())
	}

	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Could not generate master key pair: %s", err.Error())
	}

	// Encrypt master private key with master password key
	masterPrivKeySealed, err := pc.SecretboxSeal(passKey, priv[:])
	if err != nil {
		log.Fatalf("Could not encrypt master key: %s", err.Error())
	}

	passConfig := io.ConfigFile{
		MasterPrivKeySealed: masterPrivKeySealed,
		MasterPubKey:        *pub,
	}

	// Save configs to file
	if err = passConfig.SaveFile(); err != nil {
		log.Fatalf("Could not write to config file: %s", err.Error())
	} else {
		fmt.Printf("Successfully written config to masterpass file\n")
	}

	fmt.Println("Password Vault successfully initialized")
}

// Check if dir and respective folders exist and update respective variables
func checkDirAndFoldersExists() {
	if dirExists, err := io.PassDirExists(); err == nil { // Check application dir exists
		if !dirExists {
			needsDir = true
		} else {
			if _, err := io.ConfigFileExists(); err == nil { // Check config file exists
				hasConfigFile = true
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

// Create mconfig file with secure permission.
// os.Create() leaves file world-readable.
func CreateConfigFile(configFile string) {
	config, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Could not create passgo config: %s", err.Error())
	}
	config.Close()
	fmt.Printf("Successfully created config file to store configs at: %s\n", configFile)
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
