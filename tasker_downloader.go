package bdpan

import (
	"bdpan/common"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wxnacy/dler"
	"github.com/wxnacy/go-tasker"
	"github.com/wxnacy/go-tools"
)

func NewDownloadTasker(file *FileInfoDto) *DownloadTasker {
	t := DownloadTasker{
		Tasker:   tasker.NewTasker(),
		FromFile: file,
		dler:     NewDownloader(),
		Dir:      GetDefaultDownloadDir(),
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
	Path        string
	Dir         string
	IsSync      bool // 是否同步执行
	IsRecursion bool

	fromDir string      // 文件夹地址
	toDir   string      // 真实的保存目录
	dler    *Downloader // 下载器
	logFile *os.File
}

func (d *DownloadTasker) buildToDir() error {
	if d.Path != "" {
		d.toDir = d.Path
	} else {
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
	m.dler.IsNotCover = true
	logFile, err := os.OpenFile(
		logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend|os.ModePerm)
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

func (d *DownloadTasker) BuildTasks() error {
	Log.Info("开始查找文件")
	Log.Debugf("是否递归下载文件 %v", d.IsRecursion)
	// 构建下载文件夹时候的任务集合
	if d.fromDir != "" {
		err := WalkDir(d.fromDir, d.IsRecursion, func(file *FileInfoDto) error {
			info := DownloadTaskInfo{FromFile: file}
			d.AddTask(&tasker.Task{Info: info})
			return nil
		})
		if err != nil {
			return err
		}
	}
	Log.Infof("找到文件个数 %d", len(d.GetTasks()))
	return nil
}

func (d DownloadTasker) RunTask(task *tasker.Task) error {
	info := task.Info.(DownloadTaskInfo)
	fromPath := info.FromFile.Path
	// 计算保存地址
	toPath := strings.TrimLeft(strings.TrimLeft(fromPath, d.FromFile.Path), "/")
	toPath = filepath.Join(d.toDir, toPath)
	// 目录的生成
	err := tools.DirExistsOrCreate(filepath.Dir(toPath))
	if err != nil {
		return err
	}
	// 下载
	err = d.dler.DownloadFileWithPath(info.FromFile, toPath)
	if err == dler.ErrFileExists {
		return nil
	}
	if err != nil {
		Log.Errorf("RunTask %v", err)
	}
	return err
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
