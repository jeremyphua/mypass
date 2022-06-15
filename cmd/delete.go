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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete credentials",
	Long:  `Remove credentials for the specified website`,
	Run: func(cmd *cobra.Command, args []string) {
		site := strings.Join(args, " ")
		err := db.DeleteCredentials(site)
		if err != nil {
			fmt.Println("Something went wrong:", err)
			os.Exit(1)
		}
		fmt.Printf("Credentials for website %s deleted\n", site)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
