package bdpan

import (
	"bdpan/common"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/wxnacy/dler"
	"github.com/wxnacy/go-tools"
)

func NewDownloader() *Downloader {
	return &Downloader{
		Out: os.Stdout,
	}

}

type Downloader struct {
	// File           *FileInfoDto
	// To             string
	Path           string
	Dir            string
	Out            io.Writer
	IsNotCover     bool
	DisableLog     bool
	UseProgressBar bool
	isVerbose      bool
}

func (d *Downloader) Exec() error {
	// if d.File != nil {
	// return d.DownloadFile(d.File)
	// }
	return nil
}

func (d *Downloader) SetPath(path string) *Downloader {
	d.Path = path
	fmt.Println("SetPath", path, d.Path)
	return d
}

func (d *Downloader) DownloadFileWithPath(file *FileInfoDto, path string) error {
	downloadPath := d.Path
	if path != "" {
		downloadPath = path
	}
	from := file.Path
	if !d.DisableLog {
		Log.Infof("获取文件内容: %s", from)
	}

	dlink, err := GetFileDLink(file.FSID)
	Log.Debugf("%s DLink: %s", file.Path, dlink)
	if err != nil {
		return err
	}
	downloadCacheDir := JoinCache("download")
	tools.DirExistsOrCreate(downloadCacheDir)
	t := dler.NewFileDownloadTasker(dlink).
		SetDownloadPath(downloadPath).SetDownloadDir(d.Dir).
		SetCacheDir(downloadCacheDir)
	if d.isVerbose {
		t.Request.EnableVerbose()
	}
	t.Out = d.Out
	t.IsNotCover = d.IsNotCover
	t.OutputFunc = LogInfoString

	t.Config.UseProgressBar = d.UseProgressBar
	err = t.Exec()
	if err != nil {
		// 将具体任务错误信息打印
		errTasks := t.GetErrorTasks()
		for _, t := range errTasks {
			Log.Errorf("task %v: %v", t.Info, t.Err)
		}
		return err
	}
	return nil
}

func (d *Downloader) DownloadFile(file *FileInfoDto) error {
	return d.DownloadFileWithPath(file, "")
}

func getToFilePath(from, to string) (string, error) {
	var path string
	if common.DirExists(to) {
		// 只给目标目录时，自动指定保存名并处理重名问题
		path = filepath.Join(to, filepath.Base(from))
		path = tools.FileAutoReDownloadName(path)
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

func (d *Downloader) EnableVerbose() *Downloader {
	d.isVerbose = true
	return d
}
