/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/jeremyphua/mypass/initialize"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your pass vault",
	Long:  `Initialize your pass vault and generate your master password`,
	Run: func(cmd *cobra.Command, args []string) {
		initialize.Init()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
