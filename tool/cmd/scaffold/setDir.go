/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package scaffold

import (
	"fmt"

	"github.com/spf13/cobra"
)

var directory string

// setDirCmd represents the setDir command
var setDirCmd = &cobra.Command{
	Use:   "setDir",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setDir:", directory)
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setDirCmd.PersistentFlags().String("foo", "", "A help for foo")
	setDirCmd.Flags().StringVarP(&directory, "dir", "d", "", "Directory to set")
	setDirCmd.MarkFlagRequired("dir")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setDirCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	ScaffoldCmd.AddCommand(setDirCmd)
}
