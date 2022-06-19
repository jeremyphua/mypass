package show

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/disiqueira/gotree"
	"github.com/jeremyphua/mypass/io"
	"github.com/jeremyphua/mypass/pc"
)

var (
	lastPrefix      = "└──"
	regPrefix       = "├──"
	innerPrefix     = "|  "
	innerLastPrefix = "   "
)

// list all sites
func ListAll() {
	allSites := GetSiteInfoByGroup()

	showResults(allSites)

}

func GetSiteInfoByGroup() (allSites map[string]io.SiteFile) {
	allSites = map[string]io.SiteFile{}
	sf := getSiteFileContent()
	for _, s := range sf {
		slashIndex := strings.LastIndex(s.Name, "/")
		group := ""
		if slashIndex > 0 {
			group = s.Name[:slashIndex]
		}
		name := s.Name[slashIndex+1:]
		si := io.SiteInfo{
			Name: name,
		}
		if allSites[group] == nil {
			allSites[group] = []io.SiteInfo{}
		}
		allSites[group] = append(allSites[group], si)
	}
	return
}

func showResults(allSites map[string]io.SiteFile) {
	vault := gotree.New("Vault")
	for group, siteList := range allSites {
		subvault := vault.Add(group)
		for _, site := range siteList {
			subvault.Add(site.Name)
		}
	}
	fmt.Println(vault.Print())
}

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
	showUsernameAndPassword(siteInfo, masterPrivKey)
}

// GetSiteInfo returns the site information for that particular entry
// What we need from SiteInfo is the public key for the site
func GetSiteInfo(searchFor string) (si io.SiteInfo) {
	sf := getSiteFileContent()
	for _, site := range sf {
		if site.Name == searchFor {
			return site
		}
	}
	return
}

func showUsernameAndPassword(siteInfo io.SiteInfo, masterPrivKey [32]byte) {
	vault, err := io.GetVaultFolder()
	if err != nil {
		log.Fatalf("Could not get vault: %s", err.Error())
	}
	encFilePath := filepath.Join(vault, siteInfo.Name)
	encryptedPassword, err := ioutil.ReadFile(encFilePath)
	password, ok := pc.BoxOpen(encryptedPassword, &siteInfo.PubKey, &masterPrivKey)
	if !ok {
		log.Fatalf("Error decryption password")
	}
	fmt.Printf("Username: %-20s\n", siteInfo.Username)
	fmt.Printf("Password: %-20s\n", password)
}

func getSiteFileContent() (sf io.SiteFile) {
	siteFile, err := io.GetSiteFile()
	if err != nil {
		log.Fatalf("Could not get site file: %s", err.Error())
	}
	fileBytes, err := ioutil.ReadFile(siteFile)
	if err != nil {
		log.Fatalf("Could not read site file: %s", err.Error())
	}
	json.Unmarshal(fileBytes, &sf)
	return
}
