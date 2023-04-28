package bdpan

import (
	"bdpan/common"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wxnacy/go-tasker"
)

func TaskUploadDir(from, to string, isSync bool, isRecursion bool, isIncludeHide bool) error {
	// 获取准确上传目录
	toFile, err := GetFileByPath(to)
	// 目标不存在时不会报错
	if err != nil && err != ErrPathNotFound {
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
	localKey := strings.TrimLeft(strings.Replace(info.From, m.From, "", 1), "/")
	existFile, exist := m.existFileMap[localKey]
	// fmt.Println(localKey, exist)
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
	req := NewUploadFileRequest(info.From, info.To).IsSync(true)
	_, err := req.Execute()
	return err
}

func (m *UploadTasker) BeforeRun() error {
	var err error
	m.existFileMap, err = getDirFileInfoMap(m.toDir, true)
	if err != nil && err != ErrPathNotFound {
		return err
	}
	return nil
}

func (u *UploadTasker) Exec() error {
	return tasker.ExecTasker(u, u.IsSync)
}

type UFTaskInfo struct {
	Index int
	From  string
	To    string
}

type UploadFragmentTasker struct {
	*tasker.Tasker
	Request   *UploadFileRequest
	Token     *AccessToken
	UploadId  string
	TmpDir    string
	Fragments []string
	To        string
	IsSync    bool
}

func (u *UploadFragmentTasker) Build() error {
	// 创建临时目录
	var err error
	return err
}

func (u *UploadFragmentTasker) BuildTasks() error {
	for i, _path := range u.Fragments {
		taskInfo := UFTaskInfo{From: _path, To: u.To, Index: i}
		Log.Debugf("add Task: %#v", taskInfo)
		u.AddTask(&tasker.Task{Info: taskInfo})
	}
	return nil
}

func (u UploadFragmentTasker) RunTask(task *tasker.Task) error {

	after := func(from string) {
		if strings.HasPrefix(from, u.TmpDir) {
			os.Remove(from)
		}
	}

	info := task.Info.(UFTaskInfo)
	from := info.From
	file, _ := os.Open(from)
	defer after(from)
	defer file.Close()
	// 分片上传 https://pan.baidu.com/union/doc/nksg0s9vi
	_, r, err := GetClient().FileuploadApi.Pcssuperfile2(
		context.Background()).AccessToken(u.Token.AccessToken).Partseq(
		strconv.Itoa(info.Index)).Path(
		u.To).Uploadid(u.UploadId).Type_(u.Request._type).File(file).Execute()
	Log.Debugf("Pcssuperfile2 path: %s", from)
	Log.Debugf("Pcssuperfile2 resp: %v", r)
	Log.Debugf("Pcssuperfile2 error: %v", err)
	if err != nil {
		return NewRespError(r)
	}
	return nil
}

func (u *UploadFragmentTasker) BeforeRun() error {
	var err error
	return err
}

func (u *UploadFragmentTasker) AfterRun() error {
	var err error
	err = os.RemoveAll(u.TmpDir)
	return err
}

func (u *UploadFragmentTasker) Exec() error {
	return tasker.ExecTasker(u, u.IsSync)
}
