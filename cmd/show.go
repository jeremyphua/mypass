/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/jeremyphua/mypass/show"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:     "show",
	Example: "mypass show money/ocbc",
	Short:   "Print the password of a mypass site.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		show.Site(path)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
