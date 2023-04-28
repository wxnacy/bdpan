package bdpan

import (
	"bdpan/common"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Downloader struct {
	File           *FileInfoDto
	To             string
	DisableLog     bool
	UseProgressBar bool
}

func (d *Downloader) Exec() error {
	if d.File != nil {
		return d.DownloadFile(d.File, d.To)
	}
	return nil
}

func (d *Downloader) DownloadFile(file *FileInfoDto, to string) error {

	from := file.Path
	path, err := getToFilePath(from, to)
	if err != nil {
		return err
	}
	if path == "" {
		return errors.New("保存地址获取失败")
	}
	// Log.Infof("保存地址: %s", path)
	if common.FileExists(path) {
		if !d.DisableLog {
			Log.Warnf("文件已存在: %s", path)
		}
		return nil
	}
	if !d.DisableLog {
		Log.Infof("获取文件内容: %s", from)
	}

	if !d.DisableLog {
		Log.Infof("开始写入文件: %s", path)
	}
	dlink, err := GetFileDLink(file.FSID)
	Log.Debugf("%s DLink: %s", file.Path, dlink)
	if err != nil {
		return err
	}
	t := NewDownloadUrlTasker(dlink, path)
	t.contentLength = file.Size
	t.Config.UseProgressBar = d.UseProgressBar
	err = t.Exec()
	if err != nil {
		// 将具体任务错误信息打印
		errTasks := t.GetErrorTasks()
		for _, t := range errTasks {
			Log.Errorf("task %v: %v", t.Info, t.Err)
		}
		// 确保报错时也删除临时文件
		os.RemoveAll(t.cacheDir)
		return err
	}
	if !d.DisableLog {
		Log.Info("下载成功")
	}
	return nil
}

func getToFilePath(from, to string) (string, error) {
	var path string
	if common.DirExists(to) {
		// 只给目标目录时，自动指定保存名并处理重名问题
		path = filepath.Join(to, filepath.Base(from))
		path = AutoReDownloadName(path)
	} else {
		toDir := filepath.Dir(to)
		if !common.DirExists(toDir) {
			return "", errors.New(fmt.Sprintf("%s 目录不存在", toDir))
		} else {
			path = to
		}
	}
	return path, nil
}
