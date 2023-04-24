package bdpan

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/wxnacy/go-tasker"
)

func TaskUploadDir(from, to string, isSync bool) error {
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
	t := UploadTasker{Tasker: tasker.NewTasker()}
	t.From = from
	t.To = to
	_, err := os.Stat(from)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	fromBaseName := filepath.Base(t.From)
	t.toDir = filepath.Join(t.To, fromBaseName)

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
		return nil
	}
	_, err := UploadFile(info.From, info.To)
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
