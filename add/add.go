package add

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jeremyphua/mypass/io"
	"github.com/jeremyphua/mypass/pc"
	"golang.org/x/crypto/nacl/box"
)

func AddPassword(name string) {

	HandleVaultExist()

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
		log.Fatalf("Could not read config file contents: %s", err.Error())
	}

	err = json.Unmarshal(configContents, &c)
	if err != nil {
		log.Fatalf("Could not unmarshal config file contents: %s", err.Error())
	}

	masterPub := c.MasterPubKey

	// prompt for username
	username := io.Prompt(fmt.Sprintf("Enter your username for %s: ", name))

	// prompt for password
	pass, err := io.PromptPass(fmt.Sprintf("Please enter your password for %s", name))
	if err != nil {
		log.Fatalf("Could not read password: %s", err.Error())
	}

	passSealed, err := pc.BoxSeal([]byte(pass), &masterPub, priv)
	if err != nil {
		log.Fatalf("Could not seal new site password: %s", err.Error())
	}

	si := io.SiteInfo{
		PubKey:   *pub,
		Name:     name,
		Username: username,
	}

	err = si.AddFile(passSealed, name)

	if err != nil {
		log.Fatalf("Could not save site info to file: %s", err.Error())
	} else {
		fmt.Printf("Successfully added password to %s", name)
	}
}

// Handling whether vault exist
func HandleVaultExist() {
	if vf, err := io.VaultExists(); err != nil {
		log.Fatalf("Could not get vault: %s", err.Error())
	} else if vf == false {
		log.Fatalf("Vault does not exist. Run mypass init")
	}
}
