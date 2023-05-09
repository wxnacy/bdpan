/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
doc: https://pan.baidu.com/union/doc/pkuo3snyp
*/
package cmd

import (
	"bdpan"
	"os"

	"github.com/spf13/cobra"
)

var (
	downloadCommand *DownloadCommand
)

func NewDownloadCommand(c *cobra.Command) *DownloadCommand {
	cmd := &DownloadCommand{}

	pwd, _ := os.Getwd()
	c.Flags().StringVarP(&cmd.outputDir, "output-dir", "d", pwd, "保存目录。默认为当前目录")
	c.Flags().StringVarP(&cmd.outputPath, "output-path", "o", "", "保存地址。覆盖已存在文件，优先级比 --output-dir 高")

	c.Flags().BoolVar(&cmd.IsSync, "sync", false, "是否同步进行")
	c.Flags().BoolVarP(&cmd.isRecursion, "recursion", "r", false, "是否递归下载")

	return cmd
}

type DownloadCommand struct {
	From        string
	outputDir   string
	outputPath  string
	IsSync      bool
	isRecursion bool
}

func (d DownloadCommand) Run() error {
	from := d.From
	file, err := bdpan.GetFileByPath(from)
	if err != nil {
		return err
	}
	Log.Debugf("是否同步: %v", d.IsSync)
	Log.Info("开始下载")
	if file.IsDir() {
		dlTasker := bdpan.NewDownloadTasker(file)
		dlTasker.Path = d.outputPath
		dlTasker.Dir = d.outputDir
		dlTasker.IsRecursion = d.isRecursion
		err = dlTasker.Exec()
		if err == nil {
			total := len(dlTasker.GetTasks())
			succ := total - len(dlTasker.GetErrorTasks())
			Log.Infof("下载完成: %d/%d", succ, total)
		}
	} else {
		dler := bdpan.NewDownloader()
		dler.UseProgressBar = true
		dler.Path = d.outputPath
		dler.Dir = d.outputDir
		if globalArg.IsVerbose {
			dler.EnableVerbose()
		}
		err = dler.DownloadFile(file)
	}
	if err != nil {
		Log.Error(err)
		return nil
	}

	return nil
}

func runDownload(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		downloadCommand.From = args[0]
	}
	return downloadCommand.Run()
}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "下载文件",
	Example: `  bdpan download /apps/video.mp4				下载文件
  bdpan download /apps/video.mp4 -d ~/Downloads			指定下载目录
  bdpan download /apps/video.mp4 -o ~/Downloads/1.mp4		指定下载地址
	`,
	DisableFlagsInUseLine: true,
	Long:                  ``,
	RunE:                  runDownload,
}

func init() {
	downloadCommand = NewDownloadCommand(downloadCmd)
	rootCmd.AddCommand(downloadCmd)
}
