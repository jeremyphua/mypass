/*
Copyright © 2022 JEREMY PHUA <jeremyphuachengtoon@gmail.com>
*/
package cmd

import (
	"github.com/jeremyphua/mypass/edit"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove a specific site from the vault by specifying the site-path",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteName := args[0]
		edit.DeleteSite(siteName)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
