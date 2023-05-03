/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/wxnacy/go-pretty"
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

type SyncModelSlice []*SyncModel

func (s SyncModelSlice) GetWriter() io.Writer {
	return os.Stdout
}

func (s SyncModelSlice) List() []pretty.Pretty {
	slice := make([]pretty.Pretty, 0)
	for _, v := range s {
		slice = append(slice, v)
	}
	return slice
}

func (s SyncModelSlice) Len() int {
	return len(s)
}

func (s SyncModelSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s SyncModelSlice) Less(i, j int) bool { return s[i].CreateTime.Before(s[j].CreateTime) }

type SyncModel struct {
	ID         string
	Remote     string
	Local      string
	Hash       string
	Mode       SyncMode
	CreateTime time.Time
}

func (s SyncModel) getLogContent() string {
	mode := "同步"
	if s.IsBackup() {
		mode = "备份"
	}
	return fmt.Sprintf("%s ==> %s %s", s.Local, s.Remote, mode)
}

func (s SyncModel) BuildPretty() []pretty.Field {
	var data = make([]pretty.Field, 0)
	data = append(data, pretty.Field{
		Name:       "ID",
		Value:      s.ID,
		IsFillLeft: true})
	data = append(data, pretty.Field{Name: "Local", Value: s.Local})
	data = append(data, pretty.Field{Name: "Remote", Value: s.Remote})
	data = append(data, pretty.Field{Name: "Mode", Value: s.GetMode()})
	data = append(data, pretty.Field{Name: "CreateTime", Value: s.CreateTime.Format("2006-01-02 15:04:05")})
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
			Remote:     s.Remote,
			Local:      s.Local,
			Mode:       mode,
			CreateTime: time.Now(),
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
	modelSlice := make([]*SyncModel, 0)
	for _, f := range models {
		modelSlice = append(modelSlice, f)
	}
	slice := SyncModelSlice(modelSlice)
	sort.Sort(slice)
	pretty.PrintList(slice)
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
	syncCmd.Flags().StringVarP(&syncCommand.Local, "local", "L", "", "本地文件夹")
	syncCmd.Flags().BoolVarP(&syncCommand.HasHide, "hide", "H", false, "是否包含隐藏文件")
	syncCmd.Flags().BoolVarP(&syncCommand.IsOnce, "once", "o", false, "是否执行单次")
	syncCmd.Flags().BoolVarP(&syncCommand.IsBackup, "backup", "", false, "是否为备份目录")
	syncCmd.Flags().BoolVarP(&syncCommand.IsCmdAdd, "add", "", false, "添加同步目录")
	syncCmd.Flags().BoolVarP(&syncCommand.IsCmdDel, "delete", "", false, "删除同步目录")
	syncCmd.Flags().BoolVarP(&syncCommand.IsCmdList, "list", "", false, "列出同步目录")
	syncCmd.MarkFlagsRequiredTogether("remote", "local")
	rootCmd.AddCommand(syncCmd)
}
