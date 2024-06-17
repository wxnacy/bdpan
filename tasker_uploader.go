package bdpan

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/wxnacy/bdpan/common"
	"github.com/wxnacy/go-tasker"
	"github.com/wxnacy/go-tools"
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
	t, err := NewUploadTasker(from, to)
	if err != nil {
		return err
	}
	t.IsSync = isSync
	t.IsRecursion = isRecursion
	t.IsIncludeHide = isIncludeHide
	return tasker.ExecTasker(t, isSync)
}

type UploadTaskInfo struct {
	From string
	To   string
}

func NewUploadDirTasker(from, to string) (*UploadTasker, error) {
	// 获取准确上传目录
	toFile, err := GetFileByPath(to)
	// 目标不存在时不会报错
	if err != nil && err != ErrPathNotFound {
		return nil, err
	}
	if toFile != nil {
		if toFile.IsDir() {
			fromBaseName := filepath.Base(from)
			to = filepath.Join(to, fromBaseName)
		} else {
			return nil, fmt.Errorf("%s 已存在", to)
		}
	}
	// 构建上传任务
	return NewUploadTasker(from, to)
}

func NewUploadTasker(from, to string) (*UploadTasker, error) {
	Log.Debugf("NewUploadTasker from: %s, to: %s", from, to)
	t := UploadTasker{Tasker: tasker.NewTasker()}
	t.Tasker.Config.RetryMaxTime = 5
	t.From = from
	t.toDir = to
	t.To = to
	_, err := os.Stat(from)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

type UploadTasker struct {
	*tasker.Tasker
	// 迁移的地址
	From          string
	To            string
	IsSync        bool
	IsRecursion   bool // 是否递归子文件夹文件
	IsIncludeHide bool // 是否上传隐藏文件
	regexp        string

	existFileMap map[string]FileInfoDto
	toDir        string
}

func (u *UploadTasker) SetRegexp(pattern string) *UploadTasker {
	u.regexp = pattern
	return u
}

func (m *UploadTasker) AfterRun() error {
	return nil
}

func (m *UploadTasker) BuildTasks() error {
	var err error
	filter := tools.NewFileFilter(m.From, func(paths []string) error {
		for _, path := range paths {
			from := path
			subDir := m.getSubDir(path)
			info, _ := os.Stat(path)
			to := filepath.Join(m.toDir, subDir, info.Name())
			taskInfo := UploadTaskInfo{From: from, To: to}
			Log.Infof("add Task: %#v", taskInfo)
			m.AddTask(&tasker.Task{Info: taskInfo})
		}
		return nil
	})
	if m.IsIncludeHide {
		filter.WithHide()
	}
	if m.IsRecursion {
		filter.EnableRecursion()
	}
	if m.regexp != "" {
		filter.SetFilter(func(path string, info os.FileInfo, err error) (bool, error) {
			return regexp.MatchString(m.regexp, path)
		}).EnableConfirm()
	}
	err = filter.Run()
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
			Log.Infof("Upload already %s", info.From)
			return nil
		}
	}
	Log.Infof("Upload file %s to %s", info.From, info.To)
	req := NewUploadFileRequest(info.From, info.To).IsSync(true)
	_, err := req.Execute()
	return err
}

func (u *UploadTasker) BeforeRun() error {
	var err error
	u.existFileMap, err = u.GetExistFileMap()
	if err != nil && err != ErrPathNotFound {
		return err
	}
	return nil
}

func (u *UploadTasker) Exec() error {
	return tasker.ExecTasker(u, u.IsSync)
}

func (u *UploadTasker) GetExistFileMap() (map[string]FileInfoDto, error) {
	if u.existFileMap != nil {
		return u.existFileMap, nil
	}
	var err error
	u.existFileMap, err = u.getExistFileMap()
	if err != nil && err != ErrPathNotFound {
		return nil, err
	}
	return u.existFileMap, nil
}
func (u *UploadTasker) getExistFileMap() (map[string]FileInfoDto, error) {
	files, err := u.getExistFiles(u.toDir, u.From)
	if err != nil {
		return nil, err
	}

	m := map[string]FileInfoDto{}
	for _, file := range files {
		key := strings.TrimLeft(strings.Replace(file.Path, u.toDir, "", 1), "/")
		m[key] = file
	}
	return m, nil
}

func (u *UploadTasker) getExistFiles(dir, local string) ([]FileInfoDto, error) {
	files, err := getFileInfosByLocal(dir, local, u.IsIncludeHide)
	if err != nil {
		return nil, err
	}

	resFiles := make([]FileInfoDto, 0)
	for _, f := range files {
		if f.IsDir() && u.IsRecursion {
			dirFiles, err := u.getExistFiles(f.Path, "")
			if err != nil {
				return nil, err
			}
			resFiles = append(resFiles, dirFiles...)
		} else {
			resFiles = append(resFiles, *f)
		}
	}
	return resFiles, nil
}

// -------------------------------------
// UploadFragmentTasker
// -------------------------------------

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
