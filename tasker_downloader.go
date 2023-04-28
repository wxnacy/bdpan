package bdpan

import (
	"bdpan/common"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/wxnacy/go-tasker"
	"github.com/wxnacy/gotool"
)

func TaskDownloadDir(file *FileInfoDto, to string, isSync bool) error {
	t := NewDownloadTasker(to)
	t.FromFile = file
	t.IsSync = isSync
	err := t.Exec()
	if err != nil {
		return err
	}
	total := len(t.GetTasks())
	succ := total - len(t.GetErrorTasks())
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
	IsSync   bool // 是否同步执行

	fromDir string      // 文件夹地址
	toDir   string      // 真实的保存目录
	dler    *Downloader // 下载器
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

	if m.IsSync {
		m.Tasker.Config.UseProgressBar = false // 不使用进度条
	} else {
		m.dler.DisableLog = true // 不输出下载日志
	}
	return nil
}

func (m *DownloadTasker) Exec() error {
	return tasker.ExecTasker(m, m.IsSync)
}

type DownloadUrlTaskInfo struct {
	index      int
	rangeStart int
	rangeEnd   int
	tempPath   string
}

func NewDownloadUrlTasker(url, path string) *DownloadUrlTasker {
	t := tasker.NewTasker()
	return &DownloadUrlTasker{
		Tasker:      t,
		url:         url,
		path:        path,
		segmentSize: 8 * (1 << 20), // 单个分片大小
	}
}

type DownloadUrlTasker struct {
	*tasker.Tasker
	// 迁移的地址
	url           string
	path          string
	contentLength int
	segmentSize   int
	to            string
	isSync        bool // 是否同步执行
	cacheDir      string
	id            string
}

func (d *DownloadUrlTasker) Build() error {
	d.id = genId()
	d.cacheDir = filepath.Join(cacheDir, "download", d.id)
	return nil
}

func (d *DownloadUrlTasker) AfterRun() error {
	// 写入总文件
	writeFile, err := os.OpenFile(d.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, common.PermFile)
	defer os.RemoveAll(d.cacheDir)
	defer writeFile.Close()
	if err != nil {
		return err
	}
	for _, task := range d.GetTasks() {
		info := task.Info.(DownloadUrlTaskInfo)
		tempFile, err := os.Open(info.tempPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(writeFile, tempFile)
		if err != nil {
			return err
		}
		tempFile.Close()
		os.Remove(info.tempPath)
	}
	return nil
}

func (d *DownloadUrlTasker) BuildTasks() error {
	length := d.contentLength
	page := int(length/d.segmentSize) + 1
	for i := 0; i < page; i++ {
		info := DownloadUrlTaskInfo{
			index:      i,
			rangeStart: i * d.segmentSize,
			rangeEnd:   (i+1)*d.segmentSize - 1,
			tempPath:   filepath.Join(d.cacheDir, fmt.Sprintf("%s-%d", d.id, i)),
		}
		d.AddTask(&tasker.Task{Info: info})
		if info.rangeStart >= length {
			break
		}
		if info.rangeEnd >= length {
			info.rangeEnd = length - 1
		}
	}
	return nil
}

func (d DownloadUrlTasker) RunTask(task *tasker.Task) error {
	info := task.Info.(DownloadUrlTaskInfo)
	bytes, err := GetUriBytes(d.url, info.rangeStart, info.rangeEnd)
	if err != nil {
		return err
	}
	return os.WriteFile(info.tempPath, bytes, common.PermFile)
}

func (d *DownloadUrlTasker) BeforeRun() error {
	gotool.DirExistsOrCreate(d.cacheDir)
	return nil
}

func (d *DownloadUrlTasker) Exec() error {
	return tasker.ExecTasker(d, d.isSync)
}
