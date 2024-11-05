/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/wxnacy/bdpan"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清理缓存数据",
	Run: func(cmd *cobra.Command, args []string) {
		err := bdpan.CleanCache()
		handleCmdErr(err)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
