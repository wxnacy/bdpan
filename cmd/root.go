/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	globalArg   = &GlobalArg{}
	rootCommand = &RootCommand{}

	Log = bdpan.GetLogger()
)

type GlobalArg struct {
	IsVerbose bool
	AppId     string
}

type RootCommand struct {
	Dir string
}

func promptDir(dir string) error {
	files, err := bdpan.GetDirAllFiles(dir)
	if err != nil {
		return err
	}
	printTemple := " {{ .GetFileTypeIcon }} {{ .GetFilename | cyan }}"
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   promptui.IconSelect + printTemple,
		Inactive: " " + printTemple,
		Selected: "\U0001f449  {{ .GetFileTypeIcon }}  {{ .Path | red | cyan }}",
		Details: `
--------- Detail ----------
{{ "FSID:" }}	{{ .FSID }}
{{ "Name:" }}	{{ .GetFilename }}
{{ "Filetype:" }}	{{ .GetFileTypeIcon }} {{ .GetFileType }}
{{ "Size:" }}	{{ .GetSize }}
{{ "Path:" }}	{{ .Path }}
{{ "MD5:" }}	{{ .MD5 }}
{{ "CTime:" }}	{{ .GetServerCTime }}
{{ "MTime:" }}	{{ .GetServerMTime }}
{{ "LCTime:" }}	{{ .GetLocalCTime }}
{{ "LMTime:" }}	{{ .GetLocalMTime }}
`,
	}

	searcher := func(input string, index int) bool {
		pepper := files[index]
		name := strings.Replace(strings.ToLower(pepper.GetFilename()), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "选择查看的文件",
		Items:     files,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
		IsVimMode: true,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	return promptDir(files[i].Path)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bdpan",
	Short: "百度网盘命令行工具",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		err = promptDir(rootCommand.Dir)
		handleCmdErr(err)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&globalArg.AppId, "app-id", "", "指定添加 App Id")
	rootCmd.PersistentFlags().BoolVarP(&globalArg.IsVerbose, "verbose", "v", false, "打印赘余信息")

	rootCmd.PersistentFlags().StringVarP(&rootCommand.Dir, "dir", "d", "/", "文件夹")
	// 运行前全局命令
	cobra.OnInitialize(func() {
		// 打印 debug 日志
		if globalArg.IsVerbose {
			bdpan.SetLogLevel(logrus.DebugLevel)
		}
	})
}
