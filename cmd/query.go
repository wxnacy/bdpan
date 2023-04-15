/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	queryCommand *QueryCommand
)

type QueryArg struct {
	Dir string
}

func NewQueryCommand(cmd *cobra.Command) *QueryCommand {
	c := &QueryCommand{}
	cmd.Flags().StringVarP(&c.Dir, "dir", "d", "", "查询目录")
	cmd.Flags().StringVarP(&c.Key, "key", "k", "", "查询关键字")
	cmd.Flags().StringSliceVar(&c.FSIDS, "fsid", make([]string, 0), "查询 id 列表")
	return c
}

type QueryCommand struct {
	// *Command

	Dir   string
	Key   string
	FSIDS []string
}

func (q QueryCommand) Run() error {
	dir := q.Dir
	key := q.Key
	fsids := q.FSIDS

	if len(fsids) > 0 {
		_fsids := make([]uint64, 0)
		for _, fsid := range fsids {
			id, err := strconv.Atoi(fsid)
			if err != nil {
				panic(err)
			}
			_fsids = append(_fsids, uint64(id))
		}

		files, err := bdpan.GetFilesByFSIDS(_fsids)
		if err != nil {
			panic(err)
		}
		printFileInfoList(files)
		return nil
	}

	if key != "" {
		res, err := bdpan.NewFileSearchRequest(key).Dir(dir).Recursion(1).Execute()
		if err != nil {
			return err
		}
		res.Print()
		return nil
	}
	if dir != "" {
		files, err := bdpan.GetDirAllFiles(dir)
		if err != nil {
			return err
		}
		printFileInfoList(files)
		return nil
	}
	return nil
}

func printFileInfoList(files []*bdpan.FileInfoDto) {
	for _, f := range files {
		f.PrintOneLine()
	}
	fmt.Printf("Total: %d\n", len(files))
}

func runQuery(cmd *cobra.Command, args []string) error {
	return queryCommand.Run()
}

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "查询数据",
	Long:  ``,
	RunE:  runQuery,
}

func init() {
	queryCommand = NewQueryCommand(queryCmd)
	rootCmd.AddCommand(queryCmd)
}
