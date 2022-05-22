/*
Copyright © 2022 JEREMY PHUA <jeremyphuachengtoon@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "password-cli",
	Short: "A tool to manage your password",
	Long:  `password-cli is a CLI tool to enable users to manage their passwords through the terminal.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
