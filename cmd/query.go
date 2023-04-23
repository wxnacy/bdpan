/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"
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
	cmd.Flags().StringVarP(&c.Path, "path", "p", "", "文件地址")
	cmd.Flags().StringSliceVar(&c.FSIDS, "fsid", make([]string, 0), "查询 id 列表")
	return c
}

type QueryCommand struct {
	Dir   string // 目录
	Key   string // 搜索关键字
	Path  string // 文件地址
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
		bdpan.PrintFileInfoList(files)
		return nil
	}

	if key != "" {
		files, err := bdpan.SearchFiles(dir, key)
		if err != nil {
			return err
		}
		bdpan.PrintFileInfoList(files)
		return nil
	}
	if dir != "" {
		files, err := bdpan.GetDirAllFiles(dir)
		if err != nil {
			return err
		}
		bdpan.PrintFileInfoList(files)
		return nil
	}

	path := queryCommand.Path
	if path != "" {
		file, err := bdpan.GetFileByPath(path)
		if err != nil {
			return err
		}
		file.PrettyPrint()
		return nil
	}
	return nil
}

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "查询数据",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := queryCommand.Run()
		handleCmdErr(err)
	},
}

func init() {
	queryCommand = NewQueryCommand(queryCmd)
	rootCmd.AddCommand(queryCmd)
}
