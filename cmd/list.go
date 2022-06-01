/*
Copyright © 2022 JEREMY PHUA <jeremyphuachengtoon@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jeremyphua/password-cli/db"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list ( all | <website> )",
	Args:  cobra.MinimumNArgs(1), // make sure that only one arg can be passed
	Short: "Display one or more credentials",
	Long:  `Display one or more credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		var infos []db.Information
		// if first argument is all
		// example: password-cli list all
		if args[0] == "all" {
			// check if there are no other arguments after "all", if not exit
			if len(args) > 1 {
				fmt.Println("error: you must specify only one resource")
				os.Exit(1)
			} else {
				var err error
				infos, err = db.AllCredentials()
				if err != nil {
					fmt.Println("Something went wrong:", err)
					os.Exit(1)
				}
			}
		} else {
			// list by website
			var info db.Information
			var err error
			arg := strings.Join(args, " ")
			info, err = db.GetCredentials(arg)
			if err != nil {
				fmt.Println("Something went wrong:", err)
				os.Exit(1)
			}
			if (db.Information{}) == info {
				fmt.Printf("No information found for %s, please input another website\n", arg)
				os.Exit(1)
			}
			infos = append(infos, info)
		}

		if len(infos) == 0 {
			fmt.Println("You have no credentials stored! Why not add one?")
			return
		}
		fmt.Println("You have the following credentials:")
		for _, info := range infos {
			fmt.Println("---------------------------------------------------------------")
			fmt.Printf("Website/Application:\t%s\n", info.Site)
			fmt.Println("---------------------------------------------------------------")
			fmt.Printf("|%-35s|%-25s|\n", "USERNAME", "PASSWORD")
			fmt.Println("---------------------------------------------------------------")
			fmt.Printf("|%-35s|%-25s|\n", info.UserInfo.Username, info.UserInfo.Password)
		}
	},
}

// func credentials() string {
// 	reader := bufio.NewReader(os.Stdin)

// 	fmt.Print("Enter Username: ")
// 	username, _ := reader.ReadString('\n')

// 	fmt.Print("Enter Password: ")
// 	bytePassword, err := terminal.ReadPassword(0)
// 	if err == nil {
// 		fmt.Println("\nPassword typed: " + string(bytePassword))
// 	}
// 	password := string(bytePassword)

// 	return strings.TrimSpace(username), strings.TrimSpace(password)
// }

func init() {
	rootCmd.AddCommand(listCmd)
}
