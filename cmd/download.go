/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
doc: https://pan.baidu.com/union/doc/pkuo3snyp
*/
package cmd

import (
	"bdpan"
	"bdpan/common"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
	path := filepath.Join(to, filepath.Base(from))
	if common.FileExists(path) {
		return errors.New(fmt.Sprintf("下载失败 %s 已存在\n", path))
	}
	file, err := bdpan.GetFileByPath(from)
	if err != nil {
		fmt.Fprintf(os.Stderr, "下载失败 %s\n", err.Error())
		return err
	}
	// TODO: 判定 to 的类型
	bytes, err := bdpan.GetFileBytes(file.FSID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "下载失败 %s\n", err.Error())
		return err
	}

	err = os.WriteFile(path, bytes, common.PermFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "下载失败 %s\n", err.Error())
		return err
	}
	fmt.Printf("%s 下载成功\n", path)

	return nil
}

func runDownload(cmd *cobra.Command, args []string) error {
	return downloadCommand.Run()
}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: runDownload,
}

func init() {
	downloadCommand = NewDownloadCommand(downloadCmd)
	rootCmd.AddCommand(downloadCmd)
}
