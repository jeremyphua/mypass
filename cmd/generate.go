/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/jeremyphua/mypass/generate"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate a secure password",
	Example: "mypass generator",
	Long:    `Prints a randomly generated password of length 12.`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		generate.Password()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
