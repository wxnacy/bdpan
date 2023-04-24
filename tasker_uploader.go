package bdpan

import (
	"bdpan/common"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/wxnacy/go-tasker"
)

func TaskUploadDir(from, to string, isSync bool) error {
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
	return tasker.ExecTasker(t, isSync)
}

type UploadTaskInfo struct {
	From string
	To   string
}

type UploadTasker struct {
	*tasker.Tasker
	// 迁移的地址
	From   string
	To     string
	IsSync bool

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
	infos, err := ioutil.ReadDir(m.From)
	if err != nil {
		return err
	}
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		if strings.HasPrefix(info.Name(), ".") {
			continue
		}
		from := filepath.Join(m.From, info.Name())
		to := filepath.Join(m.toDir, info.Name())
		info := UploadTaskInfo{From: from, To: to}
		m.AddTask(&tasker.Task{Info: info})
	}
	return nil
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
