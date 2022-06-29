/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/jeremyphua/mypass/edit"
	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:     "rename",
	Short:   "Rename an entry in the password vault",
	Example: "mypass rename money/ocbc",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		site := args[0]
		edit.Rename(site)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
