package bdpan

import (
	"bdpan/common"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Downloader struct {
	File       *FileInfoDto
	To         string
	DisableLog bool
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
	bytes, err := GetFileBytes(file.FSID)
	if err != nil {
		return err
	}

	if !d.DisableLog {
		Log.Infof("开始写入文件: %s", path)
	}
	err = os.WriteFile(path, bytes, common.PermFile)
	if err != nil {
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
		path = filepath.Join(to, filepath.Base(from))
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