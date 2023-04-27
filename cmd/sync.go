/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/wxnacy/gotool"
)

var (
	syncCommand = &SyncCommand{}
	modelPath   = bdpan.JoinStoage("sync.json")
)

type SyncMode int

const (
	ModeBackup SyncMode = iota
	ModeSync
)

type SyncModel struct {
	ID     string
	Remote string
	Local  string
	Hash   string
	Mode   SyncMode
}

func (s SyncModel) getLogContent() string {
	mode := "同步"
	if s.IsBackup() {
		mode = "备份"
	}
	return fmt.Sprintf("%s ==> %s %s", s.Local, s.Remote, mode)
}

func (s SyncModel) BuildPrintData() []bdpan.PrettyData {
	var data = make([]bdpan.PrettyData, 0)
	data = append(data, bdpan.PrettyData{
		Name:       "ID",
		Value:      s.ID,
		IsFillLeft: true})
	data = append(data, bdpan.PrettyData{Name: "Local", Value: s.Local})
	data = append(data, bdpan.PrettyData{Name: "Remote", Value: s.Remote})
	data = append(data, bdpan.PrettyData{Name: "Mode", Value: s.GetMode()})
	return data
}

func (s SyncModel) IsBackup() bool {
	if s.Mode == ModeBackup {
		return true
	} else {
		return false
	}
}

func (s SyncModel) GetMode() string {
	switch s.Mode {
	case ModeBackup:
		return "backup"
	default:
		return "sync"
	}
}

func (s *SyncModel) BuildID() {
	s.Hash = gotool.Md5(s.Remote + s.Local)
	s.ID = shortID(s.Hash)
}

func shortID(id string) string {
	return id[0:7]
}

type SyncCommand struct {
	ID       string
	Remote   string
	Local    string
	IsBackup bool // 是否为备份
	HasHide  bool
	IsOnce   bool

	IsCmdAdd  bool
	IsCmdDel  bool
	IsCmdList bool
}

func (s SyncCommand) getModels() (m map[string]*SyncModel) {
	err := gotool.FileReadForInterface(modelPath, &m)
	if err != nil {
		m = make(map[string]*SyncModel)
	}
	return
}

func (s SyncCommand) getModel(id string) (*SyncModel, bool) {
	models := s.getModels()
	m, exits := models[id]
	return m, exits
}

func (s SyncCommand) Run() error {
	Log.Debugf("arg: %#v", s)
	if s.IsCmdAdd {
		mode := ModeSync
		if s.IsBackup {
			mode = ModeBackup
		}

		model := &SyncModel{
			Remote: s.Remote,
			Local:  s.Local,
			Mode:   mode,
		}
		model.BuildID()
		Log.Debugf("add model: %#v", model)

		key := model.ID
		models := s.getModels()
		_, exits := models[key]
		if exits {
			return fmt.Errorf("已存在该记录")
		}
		models[key] = model
		s.PrintList(models)
		return gotool.FileWriteWithInterface(modelPath, models)
	} else if s.IsCmdList {
		models := s.getModels()
		s.PrintList(models)
	} else if s.IsCmdDel {
		if s.ID == "" {
			return fmt.Errorf("--delete 缺少参数 --id")
		}
		models := s.getModels()
		delete(models, s.ID)
		s.PrintList(models)
		return gotool.FileWriteWithInterface(modelPath, models)
	} else {
		// 执行同步操作
		fmt.Println("开始进行同步操作")
		models := map[string]*SyncModel{}
		if s.ID == "" {
			models = s.getModels()
		} else {
			model, ok := s.getModel(s.ID)
			if ok {
				models[s.ID] = model
			} else {
				return fmt.Errorf("ID: %s 同步任务不存在", s.ID)
			}
		}
		for {
			for _, m := range models {
				err := s.syncModel(m)
				if err != nil {
					Log.Errorf("%s 同步报错: %v\n", m.getLogContent(), err)
				}
			}
			if s.IsOnce {
				break
			}
			time.Sleep(5 * time.Second)
		}
	}
	return nil
}

func (s SyncCommand) syncModel(m *SyncModel) error {
	fmt.Printf("开始同步 %s\n", m.getLogContent())
	if m.IsBackup() {
		uTasker := bdpan.NewUploadTasker(m.Local, m.Remote)
		uTasker.IsIncludeHide = s.HasHide
		uTasker.IsRecursion = true
		err := uTasker.Exec()
		return err
	}
	return nil
}

func (s SyncCommand) PrintList(models map[string]*SyncModel) {
	prettyList := make([]bdpan.Pretty, 0)
	for _, model := range models {
		prettyList = append(prettyList, model)
	}
	bdpan.PrettyPrintList(bdpan.PrettyList(prettyList))
}

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "同步文件夹",
	Long:  `可以对本地和远程文件夹做同步和备份两种操作`,
	Run: func(cmd *cobra.Command, args []string) {
		err := syncCommand.Run()
		handleCmdErr(err)
	},
}

func init() {
	syncCmd.Flags().StringVarP(&syncCommand.ID, "id", "", "", "任务 id")
	syncCmd.Flags().StringVarP(&syncCommand.Remote, "remote", "r", "", "远程文件夹")
	syncCmd.Flags().StringVarP(&syncCommand.Local, "local", "l", "", "本地文件夹")
	syncCmd.Flags().BoolVarP(&syncCommand.HasHide, "hide", "H", false, "是否包含隐藏文件")
	syncCmd.Flags().BoolVarP(&syncCommand.IsOnce, "once", "o", false, "是否执行单次")
	syncCmd.Flags().BoolVarP(&syncCommand.IsBackup, "backup", "", false, "是否为备份目录")
	syncCmd.Flags().BoolVarP(&syncCommand.IsCmdAdd, "add", "", false, "添加同步目录")
	syncCmd.Flags().BoolVarP(&syncCommand.IsCmdDel, "delete", "", false, "删除同步目录")
	syncCmd.Flags().BoolVarP(&syncCommand.IsCmdList, "list", "", false, "列出同步目录")
	syncCmd.MarkFlagsRequiredTogether("remote", "local")
	rootCmd.AddCommand(syncCmd)
}
