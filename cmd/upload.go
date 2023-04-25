/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"
	"bdpan/common"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	uploadCommand *UploadCommand
)

func NewUploadCommand(c *cobra.Command) *UploadCommand {
	cmd := &UploadCommand{}

	c.Flags().StringVarP(&cmd.From, "from", "f", "", "文件来源")
	c.Flags().StringVarP(&cmd.To, "to", "t", bdpan.DefaultUploadDir, "保存地址")
	c.Flags().BoolVarP(&cmd.IsRecursion, "recursion", "r", false, "是否递归子文件夹文件")
	c.Flags().BoolVarP(&cmd.IsIncludeHide, "hide", "H", false, "是否上传隐藏文件")
	c.Flags().BoolVarP(&cmd.IsSync, "sync", "s", false, "是否同步上传")
	c.MarkFlagRequired("from")
	return cmd
}

type UploadCommand struct {
	From          string
	To            string
	IsSync        bool
	IsRecursion   bool // 是否递归子文件夹文件
	IsIncludeHide bool // 是否上传隐藏文件
}

func (u UploadCommand) Run() error {
	from := u.From
	to := u.To
	if common.FileExists(from) {
		// 获取准确上传地址
		toFile, err := bdpan.GetFileByPath(to)
		if err != nil && !strings.Contains(err.Error(), "找不到") {
			return err
		}
		if toFile == nil {
			if strings.HasSuffix(to, "/") {
				to = filepath.Join(to, filepath.Base(from))
			}
		} else {
			if toFile.IsDir() {
				to = filepath.Join(to, filepath.Base(from))
			}
		}
		Log.Infof("Upload %s to %s", from, to)
		_, err = bdpan.UploadFile(from, to)
		if err != nil {
			return err
		}
		Log.Infof("File: %s upload success", from)
	} else if common.DirExists(from) {
		return bdpan.TaskUploadDir(from, to, u.IsSync, u.IsRecursion, u.IsIncludeHide)
	} else {
		return fmt.Errorf("%s 不存在", from)
	}
	return nil
}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "上传文件",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := uploadCommand.Run()
		handleCmdErr(err)
	},
}

func init() {
	uploadCommand = NewUploadCommand(uploadCmd)
	rootCmd.AddCommand(uploadCmd)
}
