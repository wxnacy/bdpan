/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wxnacy/bdpan"

	"github.com/atotto/clipboard"
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

func NewBackFileAction(dir string) FileAction {
	backupFile := &bdpan.FileInfoDto{
		Filename: "../",
		Path:     dir,
		FileType: 1,
	}
	return FileAction{
		Name:   backupFile.GetFilename(),
		Icon:   backupFile.GetFileTypeIcon(),
		File:   backupFile,
		Action: ActionBack,
		Desc:   "返回上层目录",
	}
}
func NewSystemDirAction(dir string) FileAction {
	backupFile := &bdpan.FileInfoDto{
		Filename: "操作",
		Path:     dir,
		FileType: 1,
	}
	return FileAction{
		Name:   backupFile.GetFilename(),
		Icon:   "漣",
		File:   backupFile,
		Action: ActionViewFile,
		Desc:   "操作当前目录",
	}
}

func NewViewFileActions(files []*bdpan.FileInfoDto) []FileAction {
	actions := make([]FileAction, 0)
	for _, f := range files {
		action := ActionViewFile
		if f.IsDir() {
			action = ActionViewDir
		}
		actions = append(actions, FileAction{
			Name:   f.GetFilename(),
			Icon:   f.GetFileTypeIcon(),
			File:   f,
			Action: action,
			Desc:   f.GetPretty(),
		})
	}
	return actions
}
func NewFileActions(file *bdpan.FileInfoDto) []FileAction {
	actions := make([]FileAction, 0)
	actions = append(actions, NewBackFileAction(file.Path))
	actions = append(actions, FileAction{
		Name:   "CopyPath",
		File:   file,
		Icon:   "",
		Action: ActionCopyPath,
		Desc:   "复制文件地址到剪切板",
	})
	actions = append(actions, FileAction{
		Name:   "Download",
		Icon:   "",
		File:   file,
		Action: ActionDownload,
		Desc:   "下载文件到 " + bdpan.DefaultDownloadDir,
	})
	actions = append(actions, FileAction{
		Name:   "Copy",
		Icon:   "",
		File:   file,
		Action: ActionCopy,
		Desc:   "复制文件",
	})
	actions = append(actions, FileAction{
		Name:   "Move",
		Icon:   "",
		File:   file,
		Action: ActionMove,
		Desc:   "剪切文件",
	})
	actions = append(actions, FileAction{
		Name:   "Paste",
		Icon:   "",
		File:   file,
		Action: ActionPaste,
		Desc:   "粘贴文件",
	})
	actions = append(actions, FileAction{
		Name:   "Delete",
		Icon:   "",
		File:   file,
		Action: ActionDelete,
		Desc:   "删除文件",
	})
	return actions
}

type Action int

const (
	ActionBack Action = iota
	ActionDownload
	ActionDelete
	ActionRename
	ActionCopyPath
	ActionViewDir
	ActionViewFile
	ActionCopy
	ActionMove
	ActionPaste
)

type FileAction struct {
	File   *bdpan.FileInfoDto
	Name   string
	Icon   string
	Action Action
	Desc   string
}

type RootCommand struct {
	Path  string
	Limit int

	from  string
	opera bdpan.FileManageOpera
}

func (r *RootCommand) Run() error {
	return r.viewPath(r.Path)
}

func (r *RootCommand) viewPath(path string) error {
	if path == "/" {
		return r.viewDir(path)
	}
	file, err := bdpan.GetFileByPath(r.Path)
	if err != nil {
		return err
	}
	if file.IsDir() {
		return r.viewCurrDir(file)
	} else {
		return r.viewFile(file)
	}
}

func (r *RootCommand) viewFile(file *bdpan.FileInfoDto) error {
	actions := NewFileActions(file)

	i, err := r.promptSelect("选择您想要的操作", actions, false)
	if err != nil {
		return err
	}
	return r.handleAction(actions[i])
}

// 返回上传目录
func (r *RootCommand) viewBackDir(path string) error {
	parentDir := filepath.Dir(strings.TrimRight(path, "/"))
	return r.viewDir(parentDir)

}

func (r *RootCommand) viewCurrDir(file *bdpan.FileInfoDto) error {
	if file.IsDir() {
		return r.viewDir(file.Path)
	} else {
		return r.viewDir(filepath.Dir(file.Path))
	}
}

func (r *RootCommand) viewDir(dir string) error {
	actions := make([]FileAction, 0)
	// 添加上层目录
	if dir != "/" {
		actions = append(actions, NewBackFileAction(dir))
	}
	actions = append(actions, NewSystemDirAction(dir))
	files, err := bdpan.GetDirAllFiles(dir)
	actions = append(actions, NewViewFileActions(files)...)
	if err != nil {
		return err
	}
	i, err := r.promptSelect("选择查看的文件", actions, true)
	if err != nil {
		return err
	}
	return r.handleAction(actions[i])
}

func (r *RootCommand) promptSelect(
	label string, actions []FileAction, hasSearch bool,
) (index int, err error) {
	activeTempleFmt := "%s {{ .Icon }} {{ .Name | %s}}"
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   fmt.Sprintf(activeTempleFmt, promptui.IconSelect, "blue"),
		Inactive: fmt.Sprintf(activeTempleFmt, " ", "cyan"),
		// Selected: "\U0001f449  {{ .GetFileTypeIcon }}  {{ .Path | red | cyan }}",
		Details: `
--------- {{.File.Path}} ----------
{{ .Desc }}
`,
	}
	searcher := func(input string, index int) bool {
		pepper := actions[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:        label,
		Items:        actions,
		Templates:    templates,
		Size:         r.Limit,
		IsVimMode:    true,
		HideSelected: true, // 隐藏选择后输出的内容
	}
	if hasSearch {
		prompt.Searcher = searcher
	}
	index, _, err = prompt.Run()
	if err != nil {
		return
	}
	return
}

// 返回上传目录
func (r *RootCommand) handleAction(action FileAction) error {
	var err error
	file := action.File
	switch action.Action {
	case ActionViewDir:
		return r.viewDir(file.Path)
	case ActionViewFile:
		return r.viewFile(file)
	case ActionBack:
		return r.viewBackDir(file.Path)
	case ActionDelete:
		isConfirm := bdpan.PromptConfirm(fmt.Sprintf("确认删除 %s", file.Path))
		if isConfirm {
			err := bdpan.DeleteFile(file.Path)
			if err != nil {
				return err
			}
		}
		if file.IsDir() {
			return r.viewBackDir(file.Path)
		}
		return r.viewCurrDir(file)
	case ActionDownload:
		cmd := &DownloadCommand{
			isRecursion: true,
		}
		err = cmd.Download(file)
		if err != nil {
			return err
		}
		return r.viewCurrDir(file)
	case ActionCopyPath:
		err := clipboard.WriteAll(file.Path)
		if err != nil {
			return err
		}
		Log.Infof("%s 已经复制到剪切板", file.Path)
		return r.viewCurrDir(file)
	case ActionCopy:
		r.from = file.Path
		r.opera = bdpan.OperaCopy
		Log.Infof("%s 文件已复制", file.Path)
		return r.viewCurrDir(file)
	case ActionMove:
		r.from = file.Path
		r.opera = bdpan.OperaMove
		Log.Infof("%s 文件已剪切", file.Path)
		return r.viewCurrDir(file)
	case ActionPaste:
		if r.from == "" || r.opera < 0 {
			return r.viewCurrDir(file)
		}
		dir := filepath.Dir(file.Path)
		if file.IsDir() {
			dir = file.Path
		}
		to := filepath.Join(dir, filepath.Base(r.from))
		c := &ManageCommand{
			opera: r.opera,
			path:  r.from,
			to:    to,
		}
		err = c.Exec([]string{})
		if err != nil {
			return err
		}
		return r.viewCurrDir(file)
	}
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "bdpan",
	Short:   "百度网盘命令行工具",
	Long:    ``,
	Version: "0.2.1",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		err = rootCommand.Run()
		handleCmdErr(err)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Execute()
}

func init() {
	// 全局参数
	rootCmd.PersistentFlags().StringVar(&globalArg.AppId, "app-id", "", "指定添加 App Id")
	rootCmd.PersistentFlags().BoolVarP(&globalArg.IsVerbose, "verbose", "v", false, "打印赘余信息")

	// root 参数
	rootCmd.PersistentFlags().StringVarP(&rootCommand.Path, "path", "p", "/", "直接查看文件")
	rootCmd.PersistentFlags().IntVarP(&rootCommand.Limit, "limit", "l", 10, "查询数目")
	// 运行前全局命令
	cobra.OnInitialize(func() {
		// 打印 debug 日志
		if globalArg.IsVerbose {
			bdpan.SetLogLevel(logrus.DebugLevel)
		}
	})
}
