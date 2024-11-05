/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/wxnacy/bdpan"

	"github.com/spf13/cobra"
)

var deleteCommand = &ManageCommand{opera: bdpan.OperaDelete}

func runDelete(cmd *cobra.Command, args []string) {
	handleCmdErr(deleteCommand.Exec(args))
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除文件",
	Long:  ``,
	Run:   runDelete,
}

func init() {
	deleteCmd.Flags().StringVarP(&deleteCommand.path, "path", "p", "", "操作地址")
	rootCmd.AddCommand(deleteCmd)
}
