/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/jeremyphua/mypass/edit"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:     "edit",
	Example: "mypass edit money/ocbc",
	Short:   "Change the username or password of a site in the vault.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteName := args[0]
		edit.EditInformation(siteName)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
