package bdpan

import (
	"bdpan/common"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wxnacy/dler/godler"
)

func TaskDownloadDir(file *FileInfoDto, to string, isSync bool) error {
	tasker := NewDownloadTasker(to)
	tasker.FromFile = file
	if isSync {
		return tasker.SyncExec()
	} else {
		return tasker.Exec()
	}
}

type DownloadTaskInfo struct {
	From     string
	FromFile *FileInfoDto
	To       string
}

func NewDownloadTasker(to string) *DownloadTasker {
	t := DownloadTasker{
		Tasker: godler.NewTasker(godler.NewTaskerConfig())}
	t.To = to
	return &t
}

type DownloadTasker struct {
	*godler.Tasker
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

func (m *DownloadTasker) Build() {
	if m.From != "" {
		file, err := GetFileByPath(m.From)
		if err != nil {
			panic(err)
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
				panic(err)
			}
			m.toDir = m.To
		} else {
			m.toDir = filepath.Join(m.To, filepath.Base(m.From))
		}
	}

	m.dler = &Downloader{DisableLog: true}
}

func (m *DownloadTasker) AfterRun() {
}

func (m *DownloadTasker) BuildTasks() {
	// 构建下载文件夹时候的任务集合
	if m.fromDir != "" {
		files, err := GetDirAllFiles(m.fromDir)
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			m.total++
			to := filepath.Join(m.toDir, f.GetFilename())
			info := DownloadTaskInfo{FromFile: f, To: to}
			m.AddTask(&godler.Task{Info: info})
		}
	}
}

func (m DownloadTasker) RunTask(task *godler.Task) error {
	info := task.Info.(DownloadTaskInfo)
	return m.dler.DownloadFile(info.FromFile, info.To)
}

func (m *DownloadTasker) BeforeRun() {
	if !common.DirExists(m.To) {
		panic(errors.New(fmt.Sprintf("%s 目录不存在", m.To)))
		// Log.Debugf("%s 目录不存在", m.To)
	}
	if !common.DirExists(m.toDir) {
		Log.Infof("创建目录: %s", m.toDir)
		err := os.Mkdir(m.toDir, common.PermDir)
		if err != nil {
			panic(err)
		}
	}
}

func (m *DownloadTasker) Exec() error {
	m.Build()
	m.BuildTasks()
	m.BeforeRun()
	m.Run(m.RunTask)
	m.AfterRun()
	return nil
}

func (m *DownloadTasker) SyncExec() error {
	m.Build()
	m.dler.DisableLog = false
	m.BuildTasks()
	m.BeforeRun()
	for _, task := range m.GetTasks() {
		err := m.RunTask(task)
		if err != nil {
			Log.Errorf("%s 下载失败: %v", task.Info.(DownloadTaskInfo).FromFile.Path, err)
		} else {
			m.succ++
		}
	}
	m.AfterRun()
	Log.Infof("下载完成: %d/%d", m.succ, m.total)
	return nil
}
