// Package godler  provides ...
package common

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	PermFile fs.FileMode = 0666
	PermDir              = 0755
)

// 判断地址是否存在
func FileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// 判断地址是否为目录
func DirExists(dirpath string) bool {
	info, err := os.Stat(dirpath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// 获取目录大小
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func ReadFileToMap(path string) (map[string]interface{}, error) {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var fileData map[string]interface{}
	err = json.Unmarshal(bytes, &fileData)
	if err != nil {
		return nil, err
	}
	return fileData, nil
}

func ReadFileToInterface(path string, i interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, i)
	if err != nil {
		return err
	}
	return nil

}

func WriteMapToFile(path string, data map[string]interface{}) error {
	writeBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, writeBytes, PermFile)
	if err != nil {
		return err
	}
	return nil
}

func WriteInterfaceToFile(path string, data interface{}) error {
	writeBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, writeBytes, PermFile)
	if err != nil {
		return err
	}
	return nil
}
