/*
Copyright © 2022 JEREMY PHUA <jeremyphuachengtoon@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jpct96/password-cli/db"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all credentials",
	Long:  `A list of all your websites and their credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		infos, err := db.AllCredentials()
		if err != nil {
			fmt.Println("Something went wrong:", err)
			os.Exit(1)
		}
		if len(infos) == 0 {
			fmt.Println("You have no credentials stored! Why not add one?")
			return
		}
		fmt.Println("You have the following credentials:")
		for _, info := range infos {
			fmt.Println("---------------------------------------------------------")
			fmt.Printf("Website/Application:\t%s\n", info.Site)
			fmt.Println("---------------------------------------------------------")
			fmt.Printf("|%-18s|%-18s|\n", "USERNAME", "PASSWORD")
			fmt.Println("---------------------------------------")
			fmt.Printf("|%-18s|%-18s|\n", info.UserInfo.Username, info.UserInfo.Password)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
