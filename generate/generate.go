package generate

import (
	"log"

	"github.com/jeremyphua/mypass/pc"
)

func Password() string {
	pass, err := pc.GeneratePassword()
	if err != nil {
		log.Fatalf("Could not generate password: %s", err.Error())
	}
	return pass
}
