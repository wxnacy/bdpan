package bdpan

import (
	"bdpan/common"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wxnacy/go-tasker"
)

func TaskDownloadDir(file *FileInfoDto, to string, isSync bool) error {
	tasker := NewDownloadTasker(to)
	tasker.FromFile = file
	err := tasker.Exec(isSync)
	if err != nil {
		return err
	}
	total := len(tasker.GetTasks())
	succ := total - len(tasker.GetErrorTasks())
	Log.Infof("下载完成: %d/%d", succ, total)
	return nil
}

type DownloadTaskInfo struct {
	From     string
	FromFile *FileInfoDto
	To       string
}

func NewDownloadTasker(to string) *DownloadTasker {
	t := DownloadTasker{Tasker: tasker.NewTasker()}
	t.To = to
	return &t
}

type DownloadTasker struct {
	*tasker.Tasker
	// 迁移的地址
	From     string
	FromFile *FileInfoDto
	Froms    []string
	To       string
	fromDir  string      // 文件夹地址
	toDir    string      // 真实的保存目录
	dler     *Downloader // 下载器
	total    int         // 总数
	succ     int         // 成功
}

func (m *DownloadTasker) Build() error {
	if m.From != "" {
		file, err := GetFileByPath(m.From)
		if err != nil {
			return err
		}
		m.FromFile = file
	}

	if m.FromFile != nil {
		if m.FromFile.IsDir() {
			// 判定下载来源是否为文件夹
			m.fromDir = m.FromFile.Path
		}
		if common.DirExists(filepath.Dir(m.To)) && !common.DirExists(m.To) {
			err := os.Mkdir(m.To, common.PermDir)
			if err != nil {
				return err
			}
			m.toDir = m.To
		} else {
			m.toDir = filepath.Join(m.To, filepath.Base(m.From))
		}
	}

	m.dler = &Downloader{}
	return nil
}

func (m *DownloadTasker) AfterRun() error {
	return nil
}

func (m *DownloadTasker) BuildTasks() error {
	// 构建下载文件夹时候的任务集合
	if m.fromDir != "" {
		files, err := GetDirAllFiles(m.fromDir)
		if err != nil {
			return err
		}
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			m.total++
			to := filepath.Join(m.toDir, f.GetFilename())
			info := DownloadTaskInfo{FromFile: f, To: to}
			m.AddTask(&tasker.Task{Info: info})
		}
	}
	return nil
}

func (m DownloadTasker) RunTask(task *tasker.Task) error {
	info := task.Info.(DownloadTaskInfo)
	return m.dler.DownloadFile(info.FromFile, info.To)
}

func (m *DownloadTasker) BeforeRun() error {
	if !common.DirExists(m.To) {
		return fmt.Errorf("%s 目录不存在", m.To)
	}
	if !common.DirExists(m.toDir) {
		Log.Debugf("创建目录: %s", m.toDir)
		err := os.Mkdir(m.toDir, common.PermDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *DownloadTasker) Exec(isSync bool) error {
	var err error
	err = m.Build()
	if err != nil {
		return err
	}
	err = m.BuildTasks()
	if err != nil {
		return err
	}
	err = m.BeforeRun()
	if err != nil {
		return err
	}
	if isSync {
		m.Tasker.Config.UseProgressBar = false
		err = m.SyncRun(m.RunTask)
	} else {
		m.dler.DisableLog = true
		err = m.Run(m.RunTask)
	}
	if err != nil {
		return err
	}
	return m.AfterRun()
}
