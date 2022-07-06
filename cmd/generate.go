/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/jeremyphua/mypass/generate"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate a secure password",
	Example: "mypass generator",
	Long:    `Prints a randomly generated password. The default length is 20.`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		password := generate.Password()
		fmt.Println(password)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
