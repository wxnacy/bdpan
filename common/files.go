// Package godler  provides ...
package common

import (
	"io/fs"
	"os"
	"time"

	"github.com/djherbis/times"
	"github.com/wxnacy/gotool"
)

const (
	PermFile fs.FileMode = 0666
	PermDir              = 0755
)

// 判断地址是否存在
func FileExists(filepath string) bool {
	return gotool.FileExists(filepath)
}

// 判断地址是否为目录
func DirExists(dirpath string) bool {
	return gotool.DirExists(dirpath)
}

// 获取目录大小
func DirSize(path string) (int64, error) {
	return gotool.DirSize(path)
}

func ReadFileToMap(path string) (map[string]interface{}, error) {
	return gotool.FileReadToMap(path)
}

func ReadFileToInterface(path string, i interface{}) error {
	return gotool.FileReadForInterface(path, i)
}

func WriteMapToFile(path string, data map[string]interface{}) error {
	return gotool.FileWriteWithInterface(path, data)
}

func WriteInterfaceToFile(path string, data interface{}) error {
	return gotool.FileWriteWithInterface(path, data)
}

// 获取文件的时间
func GetFileTimes(path string) (ctime, mtime time.Time, err error) {
	file, err := os.Stat(path)
	if err != nil {
		return
	}
	ctime, mtime = GetFileInfoTimes(file)
	return
}

// 获取文件的时间
func GetFileInfoTimes(f os.FileInfo) (ctime, mtime time.Time) {
	fileTime := times.Get(f)
	if fileTime.HasBirthTime() {
		ctime = fileTime.BirthTime()
	} else {
		if fileTime.HasChangeTime() {
			ctime = fileTime.ChangeTime()
		} else {
			ctime = fileTime.ModTime()
		}
	}
	mtime = fileTime.ModTime()
	return
}

func GetFileSize(path string) (int64, error) {
	file, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return file.Size(), nil
}
