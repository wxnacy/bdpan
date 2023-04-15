/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"

	"github.com/spf13/cobra"
)

var (
	deleteCommand *DeleteCommand
)

func NewDeleteCommand(c *cobra.Command) *DeleteCommand {
	cmd := &DeleteCommand{}
	c.Flags().StringVarP(&cmd.Path, "path", "p", "", "文件地址")
	return cmd
}

type DeleteCommand struct {
	Path string
}

func (d DeleteCommand) Run() error {
	path := d.Path
	if path != "" {
		err := bdpan.DeleteFile(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	return deleteCommand.Run()
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除文件",
	Long:  ``,
	RunE:  runDelete,
}

func init() {
	deleteCommand = NewDeleteCommand(deleteCmd)
	rootCmd.AddCommand(deleteCmd)
}
