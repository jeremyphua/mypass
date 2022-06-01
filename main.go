/*
Copyright © 2022 JEREMY PHUA <jeremyphuachengtoon@gmail.com>
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jeremyphua/password-cli/cmd"
	"github.com/jeremyphua/password-cli/db"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func main() {
	dbPath := filepath.Join(basepath, "password.db")
	err := db.Init(dbPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	cmd.Execute()
}
