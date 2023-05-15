/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var (
	syncCommand = &SyncCommand{}
	modelPath   = bdpan.JoinStoage("sync.json")
)

type SyncMode int

const (
	ActionSync bdpan.SelectAction = iota
)

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

func (s SyncCommand) getModelSlice() []*bdpan.SyncModel {
	models := bdpan.GetModels()
	modelSlice := make([]*bdpan.SyncModel, 0)
	for _, f := range models {
		modelSlice = append(modelSlice, f)
	}
	slice := bdpan.SyncModelSlice(modelSlice)
	sort.Sort(slice)
	return slice
}

func (s SyncCommand) getModelItems() []*bdpan.SelectItem {
	models := s.getModelSlice()
	items := make([]*bdpan.SelectItem, 0)
	for _, m := range models {
		item := &bdpan.SelectItem{
			Name:   m.Remote,
			Desc:   m.Desc(),
			Info:   m,
			Action: bdpan.ActionSystem,
		}
		items = append(items, item)
	}
	return items
}

func (s SyncCommand) selectSync() error {
	models := s.getModelItems()
	handle := func(item *bdpan.SelectItem) error {
		return s.handleAction(item)
	}
	return bdpan.PromptSelect("所有同步任务", models, true, 10, func(index int, s string) error {
		item := models[index]
		return handle(item)
	})
}

func (s SyncCommand) handleAction(item *bdpan.SelectItem) error {
	var err error
	switch item.Action {
	case bdpan.ActionSystem:
		return s.selectSystem(item)
	case ActionSync:
		m := item.Info.(*bdpan.SyncModel)
		err = m.Exec()
		if err != nil {
			return err
		}
		return s.selectSync()
	case bdpan.ActionDelete:
		m := item.Info.(*bdpan.SyncModel)
		err = bdpan.DeleteSyncModel(m.ID)
		return s.selectSync()
	}
	return nil
}

func (s SyncCommand) selectSystem(item *bdpan.SelectItem) error {
	systems := []*bdpan.SelectItem{
		&bdpan.SelectItem{
			Name:   "Sync",
			Desc:   "进行一次同步操作",
			Info:   item.Info,
			Action: ActionSync,
		},
		&bdpan.SelectItem{
			Name:   "Delete",
			Desc:   "删除操作",
			Info:   item.Info,
			Action: bdpan.ActionDelete,
		},
	}
	handle := func(item *bdpan.SelectItem) error {
		return s.handleAction(item)
	}
	return bdpan.PromptSelect("操作列表", systems, true, 5, func(index int, s string) error {
		item := systems[index]
		return handle(item)
	})
}

func (s SyncCommand) Run() error {
	Log.Debugf("arg: %#v", s)
	if s.IsCmdAdd {
		mode := bdpan.ModeSync
		if s.IsBackup {
			mode = bdpan.ModeBackup
		}

		model := bdpan.NewSyncModel(s.Local, s.Remote, mode)
		Log.Debugf("add model: %#v", model)

		key := model.ID
		models := bdpan.GetModels()
		_, exits := models[key]
		if exits {
			return fmt.Errorf("已存在该记录")
		}
		models[key] = model
		err := bdpan.SaveModels(models)
		if err != nil {
			return err
		}
		bdpan.PrintSyncModelList()
	} else if s.IsCmdList {
		bdpan.PrintSyncModelList()
	} else if s.IsCmdDel {
		if s.ID == "" {
			return fmt.Errorf("--delete 缺少参数 --id")
		}
		return bdpan.DeleteSyncModel(s.ID)
	} else {
		return s.selectSync()
	}
	return nil
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
