// Package godler  provides ...
package common

import (
	"io/fs"

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
