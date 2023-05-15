/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"

	"github.com/spf13/cobra"
)

var copyCommand = &ManageCommand{opera: bdpan.OperaCopy}

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "复制文件",
	Run: func(cmd *cobra.Command, args []string) {
		handleCmdErr(copyCommand.Exec(args))
	},
}

func init() {
	copyCmd.Flags().StringVarP(&copyCommand.path, "path", "p", "", "操作地址")
	copyCmd.Flags().StringVarP(&copyCommand.to, "to", "o", "", "复制目标")
	rootCmd.AddCommand(copyCmd)
}
