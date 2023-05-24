package bdpan

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/wxnacy/go-tools"
)

// 通过本地文件来获取已存在的文件
func getFileInfosByLocal(remote, local string, hasHide bool) ([]*FileInfoDto, error) {
	if local == "" {
		return GetDirAllFiles(remote)
	}
	var limit int32
	limit = 1000
	getFiles := func(page int) ([]*FileInfoDto, error) {
		res, err := NewFileListRequest().Dir(remote).Order("time").
			Desc(1).Limit(limit).Page(page).Execute()
		if err != nil {
			return nil, err
		}
		return res.List, nil
	}
	infos, err := ioutil.ReadDir(local)
	if err != nil {
		return nil, err
	}
	localNames := make([]string, 0)
	for _, info := range infos {
		if strings.HasPrefix(info.Name(), ".") {
			if hasHide {
				localNames = append(localNames, info.Name())
			}
		} else {
			localNames = append(localNames, info.Name())

		}
	}
	page := 1
	totalList := []*FileInfoDto{}
	for {
		files, err := getFiles(page)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			if tools.ArrayContainsString(localNames, f.GetFilename()) {

				totalList = append(totalList, f)
			}
		}
		if len(files) <= 0 || len(files) < int(limit) {
			break
		}
		if len(totalList) == len(localNames) {
			break
		}
		page++
	}
	return totalList, nil
}

func CleanCache() error {
	return os.RemoveAll(cacheDir)
}
