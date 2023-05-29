/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeConfirm
)

var (
	globalArg   = &GlobalArg{}
	rootCommand = &RootCommand{
		mode: ModeNormal,
	}

	Log = bdpan.GetLogger()
)

type GlobalArg struct {
	IsVerbose bool
	AppId     string
}

type RootCommand struct {
	Path  string
	Limit int

	from  string
	opera bdpan.FileManageOpera

	T    *Terminal
	mode Mode

	leftBox        *Box
	leftSelect     *Select
	leftDir        string
	leftSelectPath string
	leftFile       *bdpan.FileInfoDto

	midBox        *Box
	midSelect     *Select
	midDir        string
	midSelectPath string // 中间需要选中的地址
	midFile       *bdpan.FileInfoDto

	rightBox *Box

	// 按键
	prevRune rune
}

func (r *RootCommand) initViewDir(file *bdpan.FileInfoDto) error {
	path := file.Path
	if path == "/" {
		r.midDir = path
		r.midSelectPath = ""
		return nil
	}
	if file.IsDir() {
		r.midDir = file.Path
		r.midSelectPath = ""
	} else {
		r.midDir = filepath.Dir(path)
		r.midSelectPath = path
	}
	r.leftDir = filepath.Dir(r.midDir)
	r.leftSelectPath = r.midDir
	r.midFile = file
	return nil
}

func (r *RootCommand) InitScreen(file *bdpan.FileInfoDto) error {
	var err error
	err = r.initViewDir(file)
	if err != nil {
		return err
	}
	// r.T.S.Clear()
	if err = r.DrawTopLeft(file.Path); err != nil {
		return err
	}
	if err = r.DrawLayout(); err != nil {
		return err
	}
	r.T.S.Show()
	if err = r.DrawSelect(); err != nil {
		return err
	}
	return nil
}

// 画布局
func (r *RootCommand) DrawLayout() error {
	t := r.T
	w, h := t.S.Size()
	Log.Debugf("window size (%d, %d)", w, h)
	// left box
	var boxWidth = int(float64(w) * 0.2)
	var startX = 0
	var startY = 1
	var endX = startX + boxWidth
	var endY = h - 2
	r.leftBox = NewBox(startX, startY, endX, endY, t.StyleDefault)
	r.T.DrawBox(*r.leftBox)
	// left select
	r.leftSelect = &Select{
		StartX:    startX + 1,
		StartY:    startY + 1,
		MaxWidth:  boxWidth - 2,
		MaxHeight: h - 4,
		StyleSelect: tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorDarkCyan),
	}
	// mid box
	startX = endX
	boxWidth = int(float64(w) * 0.4)
	endX = startX + boxWidth
	r.midBox = NewBox(startX, startY, endX, endY, t.StyleDefault)
	r.T.DrawBox(*r.midBox)
	// mid select
	r.midSelect = &Select{
		StartX:    startX + 1,
		StartY:    startY + 1,
		MaxWidth:  boxWidth - 2,
		MaxHeight: h - 4,
		StyleSelect: tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorDarkCyan),
	}
	// left box
	startX = endX
	endX = startX + int(float64(w)*0.4)
	r.rightBox = NewBox(startX, startY, endX, endY, t.StyleDefault)
	r.T.DrawBox(*r.rightBox)
	return nil
}

func (r *RootCommand) DrawSelect() error {
	// mid select
	// draw
	r.midBox.DrawText(r.T.S, r.T.StyleDefault, "load files...")
	if r.midDir != "/" {
		r.leftBox.DrawText(r.T.S, r.T.StyleDefault, "load files...")
	}
	r.T.S.Show()
	r.fillSelectItems(r.midSelect, r.midDir, r.midSelectPath)
	r.drawSelect(r.midSelect)
	// left select
	if r.midDir != "/" {
		r.fillSelectItems(r.leftSelect, r.leftDir, r.leftSelectPath)
		r.drawSelect(r.leftSelect)
	}
	// r.T.S.Show()
	return nil
}

func (r *RootCommand) initLeftSelect(s *Select, midDir string) error {
	dir := filepath.Dir(midDir)
	r.initSelect(s, dir)
	for i, item := range s.Items {
		info := item.Info.(*bdpan.FileInfoDto)
		if info.Path == midDir {
			s.SelectIndex = i
			break
		}
	}
	return nil
}
func (r *RootCommand) fillSelectItems(s *Select, dir string, selectPath string) error {
	if len(s.Items) == 0 {
		files, err := bdpan.GetDirAllFiles(dir)
		if err != nil {
			return err
		}
		var items = make([]*SelectItem, 0)
		for i, f := range files {
			item := &SelectItem{
				Info: f,
			}
			items = append(items, item)
			if f.Path == selectPath {
				s.SelectIndex = i
			}
		}
		s.Items = items
	}
	return nil
}
func (r *RootCommand) initSelect(s *Select, dir string) error {
	if len(s.Items) == 0 {
		files, err := bdpan.GetDirAllFiles(dir)
		if err != nil {
			return err
		}
		var items = make([]*SelectItem, 0)
		for _, f := range files {
			item := &SelectItem{
				Info: f,
			}
			items = append(items, item)
		}
		// items[0].IsSelect = true
		s.SelectIndex = 0
		s.Items = items
	}
	return nil
}

func (r *RootCommand) drawSelect(s *Select) error {
	if s.Items == nil {
		return nil
	}
	for i, item := range s.GetDrawItems() {
		info := item.Info.(*bdpan.FileInfoDto)
		text := fmt.Sprintf(" %s %s", info.GetFileTypeIcon(), info.GetFilename())
		style := r.T.StyleDefault
		if i == s.SelectIndex {
			style = s.StyleSelect
		}
		// Log.Infof("drawSelect item %s style %#v", info.GetFilename(), style)
		r.T.DrawLineText(s.StartX, s.StartY+i, s.MaxWidth, style, text)
	}
	r.DrawSelectItem(s)
	return nil
}

func (r *RootCommand) DrawEventKey(ev *tcell.EventKey) error {
	// 写入 rune
	runeStr := strings.ReplaceAll(strconv.QuoteRune(ev.Rune()), "'", "")
	if runeStr == " " {
		runeStr = "space"
	}
	err := r.DrawBottomRight(runeStr)
	if err != nil {
		return err
	}
	// 写入 key
	keyStr, ok := tcell.KeyNames[ev.Key()]
	if ok {
		err = r.DrawBottomRight(keyStr)
		if err != nil {
			return err
		}
	}
	return nil
}

// 左上角输入内容
func (r *RootCommand) DrawTopLeft(text string) error {
	w, _ := r.T.S.Size()
	return r.T.DrawText(0, 0, w-1, 0, r.T.StyleDefault, text)
}

// 左下角输入内容
func (r *RootCommand) DrawBottomLeft(text string) error {
	w, h := r.T.S.Size()
	return r.T.DrawText(0, h-1, w-10, h-1, r.T.StyleDefault, text)
}

// 右下角输入内容
func (r *RootCommand) DrawBottomRight(text string) error {
	w, h := r.T.S.Size()
	drawW := 10
	return r.T.DrawLineText(w-drawW-1, h-1, drawW, r.T.StyleDefault, text)
}

// 绘制选中的 select item
func (r *RootCommand) DrawSelectItem(s *Select) {
	selectItem := s.GetSeleteItem()
	info := selectItem.Info.(*bdpan.FileInfoDto)
	r.rightBox.DrawText(r.T.S, r.T.StyleDefault, info.GetPretty())
}

// 获取被选中的文件对象
func (r *RootCommand) GetSelectInfo() *bdpan.FileInfoDto {
	return r.getSelectInfo(r.midSelect)
}

func (r *RootCommand) getSelectInfo(s *Select) *bdpan.FileInfoDto {
	item := s.GetSeleteItem()
	info := item.Info.(*bdpan.FileInfoDto)
	Log.Infof("GetSelectInfo %s", info.Path)
	return info
}

func (r *RootCommand) ListenEventKeyInModeConfirm(ev *tcell.EventKey) error {
	// 处理退出的快捷键
	if ev.Rune() == 'q' || ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
		r.mode = ModeNormal
		return nil
	}
	return nil
}
func (r *RootCommand) ListenEventKeyInModeNormal(ev *tcell.EventKey) error {
	// 处理退出的快捷键
	if ev.Rune() == 'q' || ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
		return nil
	}
	err := r.DrawEventKey(ev)
	if err != nil {
		return err
	}
	switch ev.Rune() {
	case 'j':
		if r.midSelect.MoveDownSelect(1) {
			err := r.drawSelect(r.midSelect)
			if err != nil {
				return err
			}
		}
	case 'k':
		if r.midSelect.MoveUpSelect(1) {
			err := r.drawSelect(r.midSelect)
			if err != nil {
				return err
			}
		}
	case 'l':
		r.InitScreen(r.GetSelectInfo())
	case 'y':
		switch r.prevRune {
		case 0:
			r.prevRune = 'y'
			// case 'y':
			// info := r.GetSelectInfo()
		}
	case 'h':
		leftSelectFile := r.getSelectInfo(r.leftSelect)
		file := &bdpan.FileInfoDto{
			Path:     filepath.Dir(leftSelectFile.Path),
			FileType: 1,
		}
		r.InitScreen(file)
	default:
		switch ev.Key() {
		case tcell.KeyCtrlL:
			r.T.S.Sync()
		case tcell.KeyEnter:
			r.DrawBottomLeft("确定要下载?(y/N)")
			r.mode = ModeConfirm
		}
	}
	return nil
}

func (r *RootCommand) Exec(args []string) error {
	var err error
	var file *bdpan.FileInfoDto
	file = &bdpan.FileInfoDto{
		Path:     r.Path,
		FileType: 1,
	}
	if r.Path != "/" {
		file, err = bdpan.GetFileByPath(r.Path)
		if err != nil {
			return err
		}
	}
	bdpan.SetOutputFile()
	bdpan.SetLogLevel(logrus.DebugLevel)
	t, err := NewTerminal()
	if err != nil {
		return err
	}
	defer t.Quit()
	r.T = t
	r.InitScreen(file)
	for {
		// Update screen
		t.S.Show()
		// Poll event
		ev := t.S.PollEvent()
		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			t.S.Clear()
			t.S.Sync()
			r.InitScreen(r.midFile)
		case *tcell.EventKey:
			switch r.mode {
			case ModeNormal:
				err = r.ListenEventKeyInModeNormal(ev)
			case ModeConfirm:
				err = r.ListenEventKeyInModeConfirm(ev)
			}
			if err != nil {
				return err
			}
		}
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "bdpan",
	Short:   "百度网盘命令行工具",
	Long:    ``,
	Version: "0.3.0",
	Run: func(cmd *cobra.Command, args []string) {
		handleCmdErr(rootCommand.Exec(args))
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
