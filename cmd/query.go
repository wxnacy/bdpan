/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"
	"fmt"
	"strconv"

	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"github.com/wxnacy/gotool"
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
		files, err := bdpan.SearchFiles(dir, key)
		if err != nil {
			return err
		}
		printFileInfoList(files)
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
	idMaxLen := len("fsid")
	filenameMaxLen := len("name")
	sizeLen := len("Size")
	for _, f := range files {
		var length int
		length = runewidth.StringWidth(f.GetFilename())
		if length > filenameMaxLen {
			// filenameMaxLen = len(f.GetFilename())
			filenameMaxLen = length
		}
		length = len(strconv.Itoa(int(f.FSID)))
		if length > idMaxLen {
			idMaxLen = length
		}
		length = len(gotool.FormatSize(int64(f.Size)))
		if length > sizeLen {
			sizeLen = length
		}
	}
	idFmt := fmt.Sprintf("%%%ds", idMaxLen+1)
	sizeFmt := fmt.Sprintf(" %%-%ds", sizeLen+1)
	format := fmt.Sprintf("%s %%s %%-8s %-s %%-19s %%-19s\n", idFmt, sizeFmt)
	fmt.Printf(
		format,
		"FSID",
		runewidth.FillRight("name", filenameMaxLen+1),
		"Filetype",
		"  Size",
		"  ctime",
		"  mtime",
	)
	for _, f := range files {
		fmt.Printf(
			format,
			strconv.Itoa(int(f.FSID)),
			runewidth.FillRight(f.GetFilename(), filenameMaxLen+1),
			f.GetFileType(),
			gotool.FormatSize(int64(f.Size)),
			f.GetServerCTime(),
			f.GetServerMTime(),
		)
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
