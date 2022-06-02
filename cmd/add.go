/*
Copyright © 2022 JEREMY PHUA <jeremyphuachengtoon@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jeremyphua/mypass/db"
	"github.com/spf13/cobra"
)

var username string
var password string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new credentials",
	Long:  `Add username and password for the associated website or application.`,
	Run: func(cmd *cobra.Command, args []string) {
		site := strings.Join(args, " ")
		if site == "" {
			fmt.Println("Please enter a valid website or application")
			os.Exit(1)
		}
		db.AddCredentials(site, username, password)
		fmt.Printf("You have added new credentials to %s\n Username: %s \n Password: %s", site, username, password)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&username, "username", "u", "", "Your username")
	addCmd.Flags().StringVarP(&password, "password", "p", "", "Your password")
	addCmd.MarkFlagRequired("username")
	addCmd.MarkFlagRequired("password")
}
