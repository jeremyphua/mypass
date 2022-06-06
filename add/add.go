package add

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alexedwards/argon2id"
	"github.com/jeremyphua/mypass/io"
)

func Password(name string) {

	HandleVaultExist()

	pass, err := io.PromptPass(fmt.Sprintf("Please enter your password for %s", name))
	if err != nil {
		log.Fatalf("Could not read password: %s", err.Error())
	}

	passKey, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		log.Fatalf("Could not hash master password: %s", err.Error())
	}

	vaultDir, err := io.GetVaultFolder()
	if err != nil {
		log.Fatalf("Could not get vault folder dir: %s", err.Error())
	}

	filePath := filepath.Join(vaultDir, name) // filepath is the full path of the file
	dir, _ := filepath.Split(filePath)        // dir is the path up to the final separator
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		log.Fatalf("Could not subdirectory: %s", err.Error())
	}

	ioutil.WriteFile(filePath, []byte(passKey), 0666)

	siteInfo := io.SiteInfo{
		Name:           name,
		HashedPassword: passKey,
	}

	err = siteInfo.AddFile()
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
