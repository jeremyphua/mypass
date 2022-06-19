/*
Copyright © 2022 JEREMY PHUA <jeremyphuachengtoon@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete credentials",
	Long:  `Remove credentials for the specified website`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
