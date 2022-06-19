/*
Copyright © 2022 JEREMY PHUA <jeremyphuachengtoon@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/jeremyphua/mypass/io"
	"github.com/jeremyphua/mypass/show"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mypass",
	Short: "A tool to manage your password",
	Long:  `Prints the content of your vault. If you have not initialized your vault, please run the init subcommand to get started.`,
	Run: func(cmd *cobra.Command, args []string) {
		if exists, _ := io.VaultExists(); exists {
			show.ListAll()
		} else {
			cmd.Help()
		}
	},
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
