package edit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jeremyphua/mypass/io"
	"github.com/jeremyphua/mypass/pc"
)

func EditInformation(name string) {
	validInput := false

	for validInput == false {
		// prompt user whether they want to change username or password
		usernameOrPassword := io.Prompt(fmt.Sprintf("Do you want to change your username or password for %s?\n", name))
		if usernameOrPassword == "password" {
			editPassword(name)
			validInput = true
		} else if usernameOrPassword == "username" {
			editUserName(name)
			validInput = true
		} else {
			fmt.Println("Invalid input. Please choose either username or password.")
		}
	}
}

func editPassword(name string) {
	// get sites.json
	sites := io.GetSites()
	for index, siteInfo := range sites {
		if siteInfo.Name == name {
			// validate master password
			// assign to empty variable because we do not need the master private key
			_ = pc.GetMasterPrivKey()
			newPass, err := io.PromptPass(fmt.Sprintf("Enter new password for %s", name))
			if err != nil {
				log.Fatalf("Could not read entered password: %s", err)
			}
			newSiteInfo, passSealed := pc.ReEncrypt(siteInfo, newPass)
			sites[index] = newSiteInfo
			// update sites.json
			err = io.UpdateSiteFile(sites)
			if err != nil {
				log.Fatalf("Could not edit %s in sites.json: %s", name, err.Error())
			}
			// update sealed password in vault folder
			err = io.UpdateVaultFile(name, passSealed)
			if err != nil {
				log.Fatalf("Could not edit password in %s: %s", name, err.Error())
			}
		}
	}
}

func editUserName(name string) {
	// get sites.json
	sites := io.GetSites()
	for index, siteInfo := range sites {
		if siteInfo.Name == name {
			// validate master password
			// assign to empty variable because we do not need the master private key
			_ = pc.GetMasterPrivKey()
			newUsername := io.Prompt(fmt.Sprintf("Enter new username for %s: ", name))
			siteInfo.Username = newUsername
			sites[index] = siteInfo
			// update sites.json
			err := io.UpdateSiteFile(sites)
			if err != nil {
				log.Fatalf("Could not edit %s in sites.json: %s", name, err.Error())
			}
		}
	}
}

func DeleteSite(site string) {
	sites := io.GetSites()
	pathIndex := -1
	for index, siteInfo := range sites {
		if siteInfo.Name == site {
			// validate master password
			// assign to empty variable because we do not need the master private key
			_ = pc.GetMasterPrivKey()
			pathIndex = index
			vault, err := io.GetVaultFolder()
			if err != nil {
				go log.Fatalf("Could not get vault: %s", err.Error())
			}
			filePath := filepath.Join(vault, site)
			err = os.Remove(filePath)
			if err != nil {
				log.Fatalf("Attempted to remove file but was unable to: %s", err.Error())
			}
			break
		}
	}
	if pathIndex == -1 {
		log.Fatalf("Could not find %s in vault", site)
	}
	sites = append(sites[:pathIndex], sites[pathIndex+1:]...)
	err := io.UpdateSiteFile(sites)
	if err != nil {
		log.Fatalf("Could not update password vault: %s", err.Error())
	}
	fmt.Printf("Successfully deleted credentials for %s", site)
}

func Rename(site string) {
	sites := io.GetSites()
	for index, siteInfo := range sites {
		if siteInfo.Name == site {
			// validate master password
			// assign to empty variable because we do not need the master private key
			_ = pc.GetMasterPrivKey()
			newSiteName := io.Prompt(fmt.Sprintf("Enter new sitename for %s: ", site))
			sites[index] = io.SiteInfo{
				PubKey:   siteInfo.PubKey,
				Name:     newSiteName,
				Username: siteInfo.Username,
			}
			err := io.UpdateSiteFile(sites)
			if err != nil {
				log.Fatalf("Could not edit %s in sites.json: %s", site, err.Error())
			}
			io.UpdateFileName(site, newSiteName)
		}
	}
}
