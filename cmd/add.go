/*
Copyright © 2022 JEREMY PHUA <jeremyphuachengtoon@gmail.com>
*/
package cmd

import (
	"github.com/jeremyphua/mypass/add"
	"github.com/spf13/cobra"
)

var username string
var password string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a file or password to your vault",
	Example: "mypass add money/bank.com",
	Long:    `Add a site to your password store. This site can optionally be a part of a group by prepending a group name and slash to the site name. Will prompt for confirmation when a site path is not unique.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteName := args[0]
		add.Password(siteName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
