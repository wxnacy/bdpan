package bdpan

import (
	"bdpan/common"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wxnacy/go-tasker"
	"github.com/wxnacy/go-tools"
)

// func TaskDownloadDir(file *FileInfoDto, to string, isSync bool) error {
// t := NewDownloadTasker(to)
// t.FromFile = file
// t.IsSync = isSync
// err := t.Exec()
// if err != nil {
// return err
// }
// total := len(t.GetTasks())
// succ := total - len(t.GetErrorTasks())
// Log.Infof("下载完成: %d/%d", succ, total)
// return nil
// }

func NewDownloadTasker(file *FileInfoDto) *DownloadTasker {
	t := DownloadTasker{
		Tasker:   tasker.NewTasker(),
		FromFile: file,
		dler:     NewDownloader(),
	}
	return &t
}

type DownloadTaskInfo struct {
	From     string
	FromFile *FileInfoDto
	To       string
}

type DownloadTasker struct {
	*tasker.Tasker
	// 迁移的地址
	From     string
	FromFile *FileInfoDto
	Froms    []string
	// To       string
	Path   string
	Dir    string
	IsSync bool // 是否同步执行

	fromDir string      // 文件夹地址
	toDir   string      // 真实的保存目录
	dler    *Downloader // 下载器
	logFile *os.File
}

func (d *DownloadTasker) buildToDir() error {
	if d.Path != "" {
		d.toDir = d.Path
	} else {
		if d.Dir == "" {
			pwd, _ := os.Getwd()
			d.Dir = pwd
		}
		d.toDir = filepath.Join(d.Dir, filepath.Base(d.FromFile.Path))
	}

	if !tools.PathDirExists(d.toDir) {
		return fmt.Errorf("%s 目录不存在", d.toDir)
	}
	if tools.FileExists(d.toDir) {
		return fmt.Errorf("%s 是已存在的文件", d.toDir)
	}
	return tools.DirExistsOrCreate(d.toDir)
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
		err := m.buildToDir()
		if err != nil {
			return err
		}
	}
	m.dler.Dir = m.toDir
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	m.dler.Out = logFile
	m.logFile = logFile
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
			to := filepath.Join(m.toDir, f.GetFilename())
			// FIX: 需要修改 to
			info := DownloadTaskInfo{FromFile: f, To: to}
			m.AddTask(&tasker.Task{Info: info})
		}
	}
	return nil
}

func (m DownloadTasker) RunTask(task *tasker.Task) error {
	info := task.Info.(DownloadTaskInfo)
	m.dler.Path = info.To
	return m.dler.DownloadFile(info.FromFile, "")
}

func (m *DownloadTasker) BeforeRun() error {
	if !common.DirExists(m.toDir) {
		Log.Debugf("创建目录: %s", m.toDir)
		err := os.Mkdir(m.toDir, common.PermDir)
		if err != nil {
			return err
		}
	}

	if m.IsSync {
		m.Tasker.Config.UseProgressBar = false // 不使用进度条
	} else {
		m.dler.DisableLog = true // 不输出下载日志
	}
	return nil
}

func (m *DownloadTasker) Exec() error {
	defer m.logFile.Close()
	return tasker.ExecTasker(m, m.IsSync)
}
