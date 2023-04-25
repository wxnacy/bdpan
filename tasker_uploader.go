package bdpan

import (
	"bdpan/common"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wxnacy/go-tasker"
)

func TaskUploadDir(from, to string, isSync bool, isRecursion bool, isIncludeHide bool) error {
	// 获取准确上传目录
	toFile, err := GetFileByPath(to)
	if err != nil && !strings.Contains(err.Error(), "找不到") {
		return err
	}
	if toFile != nil {
		if toFile.IsDir() {
			fromBaseName := filepath.Base(from)
			to = filepath.Join(to, fromBaseName)
		} else {
			return fmt.Errorf("%s 已存在", to)
		}
	}
	// 构建上传任务
	t := NewUploadTasker(from, to)
	t.IsSync = isSync
	t.IsRecursion = isRecursion
	t.IsIncludeHide = isIncludeHide
	return tasker.ExecTasker(t, isSync)
}

type UploadTaskInfo struct {
	From string
	To   string
}

type UploadTasker struct {
	*tasker.Tasker
	// 迁移的地址
	From          string
	To            string
	IsSync        bool
	IsRecursion   bool // 是否递归子文件夹文件
	IsIncludeHide bool // 是否上传隐藏文件

	existFileMap map[string]FileInfoDto
	toDir        string
}

func NewUploadTasker(from, to string) *UploadTasker {
	Log.Debugf("NewUploadTasker from: %s, to: %s", from, to)
	t := UploadTasker{Tasker: tasker.NewTasker()}
	t.From = from
	t.toDir = to
	_, err := os.Stat(from)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	return &t
}

func (m *UploadTasker) AfterRun() error {
	return nil
}

func (m *UploadTasker) BuildTasks() error {
	var err error
	err = filepath.Walk(m.From,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// 不处理文件夹
			if info.IsDir() {
				return nil
			}
			dirName := filepath.Base(filepath.Dir(path))
			// 判断是否处理隐藏文件
			if (strings.HasPrefix(info.Name(), ".") || strings.HasPrefix(dirName, ".")) && !m.IsIncludeHide {
				return nil
			}
			// 判断是否递归处理
			if path != filepath.Join(m.From, info.Name()) && !m.IsRecursion {
				return nil
			}
			from := path
			subDir := m.getSubDir(path)
			to := filepath.Join(m.toDir, subDir, info.Name())
			taskInfo := UploadTaskInfo{From: from, To: to}
			Log.Debugf("add Task: %#v", taskInfo)
			m.AddTask(&tasker.Task{Info: taskInfo})
			return nil
		})
	Log.Debugf("BuildTasks error: %v", err)
	return err
}

// 获取子目录
func (u *UploadTasker) getSubDir(path string) string {
	newPath := strings.Replace(path, u.From, "", 1)
	newPath = strings.TrimLeft(newPath, "/")
	return filepath.Dir(newPath)
}

func (m UploadTasker) RunTask(task *tasker.Task) error {
	info := task.Info.(UploadTaskInfo)
	existFile, exist := m.existFileMap[filepath.Base(info.From)]
	if exist && existFile.Size > 0 {
		// 对比已存在文件的修改时间是否相同，否则重新上传
		_, mtime, err := common.GetFileTimes(info.From)
		if err != nil {
			return err
		}
		if existFile.LocalMTime == mtime.Unix() {
			Log.Debugf("%s upload already", info.From)
			return nil
		}
	}
	req := NewUploadFileRequest(info.From, info.To)
	_, err := req.Execute()
	return err
}

func (m *UploadTasker) BeforeRun() error {
	var err error
	m.existFileMap, err = getDirFileInfoMap(m.toDir)
	if err != nil {
		return err
	}
	return nil
}

func (u *UploadTasker) Exec() error {
	return tasker.ExecTasker(u, u.IsSync)
}
