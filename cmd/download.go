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
	"strings"

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

func (d DownloadCommand) getToFilePath(from, to string) (string, error) {
	var path string
	if common.DirExists(to) {
		path = filepath.Join(to, filepath.Base(from))
	} else {
		toDir := filepath.Dir(to)
		if !common.DirExists(toDir) {
			return "", errors.New(fmt.Sprintf("%s 目录不存在", toDir))
		} else {
			path = to
		}
	}
	return path, nil
}

func (d DownloadCommand) downloadFile(file *bdpan.FileInfoDto, from, to string) error {
	path, err := d.getToFilePath(from, to)
	if err != nil {
		return err
	}
	if path == "" {
		return errors.New("保存地址获取失败")
	}
	// Log.Infof("保存地址: %s", path)
	if common.FileExists(path) {
		bdpan.Log.Warnf("文件已存在: %s", path)
		return nil
	}
	Log.Infof("获取文件内容: %s", from)
	bytes, err := bdpan.GetFileBytes(file.FSID)
	if err != nil {
		return err
	}

	Log.Infof("开始写入文件: %s", path)
	err = os.WriteFile(path, bytes, common.PermFile)
	if err != nil {
		return err
	}
	Log.Info("下载成功")
	return nil
}

func (d DownloadCommand) downloadDir(dir *bdpan.FileInfoDto, from, to string) error {
	if !common.DirExists(to) {
		return errors.New(fmt.Sprintf("%s 目录不存在", to))
	}
	Log.Info("搜索目录内的文件")
	files, err := bdpan.GetDirAllFiles(from)
	if err != nil {
		return err
	}
	total := 0
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		total++
	}
	Log.Infof("找到 %d 个可下载文件", total)
	Log.Info("开始下载")
	to = filepath.Join(to, filepath.Base(from))
	if !common.DirExists(to) {
		Log.Infof("创建目录: %s", to)
		err = os.Mkdir(to, common.PermDir)
		if err != nil {
			return err
		}
	}
	succ := 0
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fromPath := filepath.Join(from, f.GetFilename())
		err = d.downloadFile(
			f, fromPath,
			filepath.Join(to, f.GetFilename()),
		)
		if err != nil {
			Log.Errorf("%s 下载失败: %v", fromPath, err)
		} else {
			succ++
		}
	}
	Log.Infof("下载完成: %d/%d", succ, total)
	return nil
}

func (d DownloadCommand) Run() error {
	from := d.From
	to := d.To
	from = strings.TrimRight(from, "/")
	to, err := homedir.Expand(to)
	if err != nil {
		return err
	}
	Log.Infof("开始搜索文件: %s", from)
	file, err := bdpan.GetFileByPath(from)
	if err != nil {
		return err
	}
	Log.Infof("查询到文件类型为: %s", file.GetFileType())
	if file.IsDir() {
		err = d.downloadDir(file, from, to)
	} else {
		err = d.downloadFile(file, from, to)
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
