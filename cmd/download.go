/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
doc: https://pan.baidu.com/union/doc/pkuo3snyp
*/
package cmd

import (
	"bdpan"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	downloadCommand *DownloadCommand
)

func NewDownloadCommand(c *cobra.Command) *DownloadCommand {
	cmd := &DownloadCommand{}

	c.Flags().StringVarP(&cmd.From, "from", "f", "", "网盘文件地址")
	c.Flags().StringVarP(&cmd.To, "to", "t", bdpan.DefaultDownloadDir, "保存地址")
	c.Flags().BoolVar(&cmd.IsSync, "sync", false, "是否同步进行")
	c.MarkFlagRequired("from")
	return cmd
}

type DownloadCommand struct {
	From   string
	To     string
	IsSync bool
}

func (d DownloadCommand) Run() error {
	from := d.From
	to := d.To
	to, err := homedir.Expand(to)
	if err != nil {
		return err
	}
	file, err := bdpan.GetFileByPath(from)
	if err != nil {
		return err
	}
	if file.IsDir() {
		err = bdpan.TaskDownloadDir(file, to, d.IsSync)
	} else {
		dler := &bdpan.Downloader{}
		err = dler.DownloadFile(file, to)
	}
	if err != nil {
		Log.Error(err)
		return nil
	}

	return nil
}

func runDownload(cmd *cobra.Command, args []string) error {
	return downloadCommand.Run()
}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "下载文件",
	Long:  ``,
	RunE:  runDownload,
}

func init() {
	downloadCommand = NewDownloadCommand(downloadCmd)
	rootCmd.AddCommand(downloadCmd)
}
