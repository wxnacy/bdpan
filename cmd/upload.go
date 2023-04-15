/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"
	"bdpan/common"
	"errors"
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
	c.Flags().StringVarP(&cmd.To, "to", "t", bdpan.DEFAULT_UPLOAD_DIR, "保存地址")
	c.Flags().BoolVar(&cmd.IsSync, "sync", false, "是否同步上传")
	c.MarkFlagRequired("from")
	return cmd
}

type UploadCommand struct {
	From   string
	To     string
	IsSync bool
}

func (u UploadCommand) Run() error {
	from := u.From
	to := u.To
	if common.FileExists(from) {
		if strings.HasSuffix(to, "/") {
			to = filepath.Join(to, filepath.Base(from))
		}
		bdpan.Log.Infof("Upload %s to %s", from, to)
		_, err := bdpan.UploadFile(from, to)
		if err != nil {
			return err
		}
		bdpan.Log.Infof("File: %s upload success", from)
	} else if common.DirExists(from) {
		if u.IsSync {

			res, err := bdpan.UploadDir(from, to)
			if err != nil {
				return err
			}
			bdpan.Log.Infof("Success: %d", res.SuccessCount)
			bdpan.Log.Infof("Failed: %d", res.FailedCount)
		} else {
			bdpan.TaskUploadDir(from, to)
		}
	} else {
		return errors.New(fmt.Sprintf("%s 文件路径不存在", from))
	}
	return nil
}

func runUpload(cmd *cobra.Command, args []string) error {
	return uploadCommand.Run()
}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "上传文件",
	Long:  ``,
	RunE:  runUpload,
}

func init() {
	uploadCommand = NewUploadCommand(uploadCmd)
	rootCmd.AddCommand(uploadCmd)
}
