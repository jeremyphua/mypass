package show

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/jeremyphua/mypass/io"
	"github.com/jeremyphua/mypass/pc"
)

// Site will print out the password of the site that matches path
func Site(path string) {
	// get site information from sites.json
	// site, err := GetSiteInfo(path)
	siteInfo := GetSiteInfo(path)
	if (io.SiteInfo{}) == siteInfo {
		log.Fatalf("Site with path %s not found", path)
	}

	// get master private key
	masterPrivKey := pc.GetMasterPrivKey()

	// show password
	showPassword(siteInfo, masterPrivKey)
}

// GetSiteInfo returns the site information for that particular entry
// What we need from SiteInfo is the public key for the site
func GetSiteInfo(searchFor string) (si io.SiteInfo) {
	siteFile, err := io.GetSiteFile()
	if err != nil {
		log.Fatalf("Could not get site file: %s", err.Error())
	}
	fileBytes, err := ioutil.ReadFile(siteFile)
	if err != nil {
		log.Fatalf("Could not read site file: %s", err.Error())
	}
	var s []io.SiteInfo
	json.Unmarshal(fileBytes, &s)
	for _, site := range s {
		if site.Name == searchFor {
			return site
		}
	}
	return
}

func showPassword(siteInfo io.SiteInfo, masterPrivKey [32]byte) {
	vault, err := io.GetVaultFolder()
	if err != nil {
		log.Fatalf("Could not get vault: %s", err.Error())
	}
	encFilePath := filepath.Join(vault, siteInfo.Name)
	encryptedPassword, err := ioutil.ReadFile(encFilePath)
	decrypted, ok := pc.BoxOpen(encryptedPassword, &siteInfo.PubKey, &masterPrivKey)
	if !ok {
		log.Fatalf("Error decryption password")
	}
	fmt.Println(string(decrypted))
}
