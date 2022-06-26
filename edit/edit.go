package edit

import (
	"fmt"
	"log"

	"github.com/jeremyphua/mypass/io"
	"github.com/jeremyphua/mypass/pc"
)

func EditInformation(name string) {

	// prompt user whether they want to change username or password

	// if password
	// get sites.json
	sites := io.GetSites()
	for index, siteInfo := range sites {
		if siteInfo.Name == name {
			// validate master password
			// assign to empty variable because we do not need the master private key
			_ = pc.GetMasterPrivKey()
			newPass, err := io.PromptPass(fmt.Sprintf("Enter new password for %s: ", name))
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
