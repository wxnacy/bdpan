/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/wxnacy/bdpan"

	"github.com/spf13/cobra"
	"github.com/wxnacy/go-tools"
)

var (
	listCommand *ListCommand
	ORDERS      = []string{ORDER_NAME, ORDER_TIME, ORDER_SIZE}
)

const (
	ORDER_TIME = "time"
	ORDER_NAME = "name"
	ORDER_SIZE = "size"
)

func NewListCommand(c *cobra.Command) *ListCommand {
	cmd := &ListCommand{}
	c.Flags().StringVarP(&cmd.Dir, "dir", "d", "", "查询目录")
	c.Flags().BoolVarP(&cmd.IsRecursion, "recursion", "r", false, "是否遍历子目录，默认否")
	c.Flags().BoolVarP(&cmd.IsDesc, "desc", "", false, "是否为倒序，默认否")
	c.Flags().Int32VarP(&cmd.Start, "start", "s", 0, "查询起点，默认为0")
	c.Flags().Int32VarP(&cmd.Limit, "limit", "l", 1000, "查询数目，默认为1000")
	c.Flags().StringVarP(&cmd.Order, "order", "o", "name", "排序字段:time(修改时间)，name(文件名)，size(大小，目录无大小)")
	c.MarkFlagRequired("dir")
	return cmd
}

type ListCommand struct {
	Dir         string
	IsRecursion bool
	IsDesc      bool
	Start       int32
	Limit       int32
	Order       string
}

func (l ListCommand) Run() error {
	dir := l.Dir
	var recursion int32
	if l.IsRecursion {
		recursion = 1
	}
	var desc int32
	if l.IsDesc {
		desc = 1
	}
	res, err := bdpan.NewFileListAllRequest(dir).Recursion(recursion).Limit(
		l.Limit).Order(l.Order).Desc(desc).Start(l.Start).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "查询失败 %s", err.Error())
		return err
	}
	bdpan.PrintFileInfoList(res.List)
	if res.HasMore == 1 {
		fmt.Printf("下一页查询命令 %d\n", res.Cursor)

	}
	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	return listCommand.Run()
}

func validArgs(cmd *cobra.Command, args []string) error {
	if !tools.ArrayContainsString(ORDERS, listCommand.Order) {
		return errors.New(fmt.Sprintf("order: %s not in %v", listCommand.Order, ORDERS))
	}
	return nil
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "展示文件列表",
	Long:  ``,
	Args:  validArgs,
	RunE:  runList,
}

func init() {
	listCommand = NewListCommand(listCmd)
	rootCmd.AddCommand(listCmd)
}
