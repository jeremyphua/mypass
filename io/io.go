package io

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	SiteFileName    = "sites.json"
	ConfigFileName  = "masterpass"
	VaultFolderName = "vault"
)

type ConfigFile struct {
	MasterPassKey       []byte
	MasterPrivKeySealed []byte
	MasterPubKey        [32]byte
}

// SiteInfo represents a single saved password entry.
type SiteInfo struct {
	PubKey   [32]byte
	Name     string
	Username string
}

// contents of sites.json
type SiteFile []SiteInfo

// return error if mypass directory does not exist
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

// Returns dir of application,
// Example: C:\Users\<name of user>\.mypass
func GetPassDir() (d string, err error) {
	home, err := getHomeDir()
	if err == nil {
		d = filepath.Join(home, ".mypass")
		return
	}
	return
}

// Returns user home dir
// Example: C:\Users\<name of user>
func getHomeDir() (d string, err error) {
	d, err = os.UserHomeDir()
	return
}

// Returns dir of sites.json file
// Example: C:\Users\<name of user>\.mypass\sites.json
func GetSiteFile() (d string, err error) {
	p, err := GetPassDir()
	if err == nil {
		d = filepath.Join(p, SiteFileName)
	}
	return
}

// Determine if config file already exist
func ConfigFileExists() (bool, error) {
	c, err := GetConfigFile()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(c)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetConfigFile is used to get the user's config file
// Example: C:\Users\<name of user>\.mypass\masterpass
func GetConfigFile() (p string, err error) {
	d, err := GetPassDir()
	if err == nil {
		p = filepath.Join(d, ConfigFileName)
	}
	return
}

// Check if sites.json exist
func SiteFileExists() (bool, error) {
	p, err := GetPassDir()
	if err != nil {
		return false, err
	}
	SiteFilePath := filepath.Join(p, SiteFileName)
	_, err = os.Stat(SiteFilePath)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Check if vault folder exist
// Vault folder is used to store all the encrypted password in their respective folder
func VaultExists() (bool, error) {
	v, err := GetVaultFolder()
	if err != nil {
		return false, err
	}
	dirInfo, err := os.Stat(v)
	if err == nil && dirInfo.IsDir() {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Get vault folder dir
// Example: C:\Users\<name of user>\.mypass\vault
func GetVaultFolder() (v string, err error) {
	d, err := GetPassDir()
	if err == nil {
		v = filepath.Join(d, VaultFolderName)
	}
	return
}

// Add SiteInfo to sites.json
func (s *SiteInfo) AddFile(fileBytes []byte, filename string) error {
	vault, err := GetVaultFolder()
	if err != nil {
		return err
	}
	// Make sure that the file directory exists.
	fileDirExists, err := ConfigFileExists()
	if err != nil {
		return err
	}
	if !fileDirExists {
		err = os.Mkdir(vault, 0700)
		if err != nil {
			log.Fatalf("Could not create passgo encrypted file dir: %s", err.Error())
		}
	}
	encFilePath := filepath.Join(vault, filename)
	dir, _ := filepath.Split(encFilePath)
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		log.Fatalf("Could not create subdirectory: %s", err.Error())
	}
	err = ioutil.WriteFile(encFilePath, fileBytes, 0666)
	if err != nil {
		return err
	}

	// We still need to add this site info to the bytes.
	return s.AddSite()
}

// AddSite is used by individual password entries to update the vault.
func (s *SiteInfo) AddSite() (err error) {
	siteFile := GetSites()
	for _, si := range siteFile {
		if s.Name == si.Name {
			return errors.New("Could not add site with duplicate name")
		}
	}
	siteFile = append(siteFile, *s)
	return UpdateSiteFile(siteFile)
}

// Returns SiteFile which is a slice of SiteInfo
func GetSites() (s SiteFile) {
	si, err := GetSiteFile()
	if err != nil {
		log.Fatalf("Could not get site file: %s", err.Error())
	}
	siteFileContents, err := ioutil.ReadFile(si)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Could not open site file. Run mypass init.: %s", err.Error())
		}
		log.Fatalf("Could not read site file: %s", err.Error())
	}
	err = json.Unmarshal(siteFileContents, &s)
	if err != nil {
		log.Fatalf("Could not unmarshal site info: %s", err.Error())
	}
	return
}

// UpdateSiteFile is used to replace the current sites.json.
func UpdateSiteFile(s SiteFile) (err error) {
	si, err := GetSiteFile()
	if err != nil {
		log.Fatalf("Could not get site file: %s", err.Error())
	}
	siteFileContents, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		log.Fatalf("Could not marshal site info: %s", err.Error())
	}

	// Write the site with the newly appended site to the file.
	err = ioutil.WriteFile(si, siteFileContents, 0666)
	return
}

func (c *ConfigFile) SaveFile() (err error) {
	if exists, err := ConfigFileExists(); err != nil {
		log.Fatalf("Could not find config file: %s", err.Error())
	} else if !exists {
		log.Fatalf("pass config could not be found %s", err.Error())
	}
	cBytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		log.Fatalf("Could not marshal config file: %s", err.Error())
	}
	cfg, err := GetConfigFile()
	if err != nil {
		log.Fatalf("Could not get config file: %s", err.Error())
	}
	err = ioutil.WriteFile(cfg, cBytes, 0666)
	return
}

func PromptPass(prompt string) (pass string, err error) {
	fd := int(os.Stdin.Fd())
	fmt.Printf("%s: ", prompt)
	passBytes, err := terminal.ReadPassword(fd)
	fmt.Println("")
	return string(passBytes), err
}

func PromptUsername(path string) (input string) {
	fmt.Printf("Enter your username for %s: ", path)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("An error occured while reading input. Please try again", err)
	}
	// remove the delimeter from the string
	input = strings.TrimSuffix(input, "\r\n")
	return
}
