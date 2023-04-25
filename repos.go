package bdpan

import (
	"strings"
)

func getDirFileInfoMap(dir string, isRecursion bool) (map[string]FileInfoDto, error) {
	files, err := getFileInfos(dir, isRecursion)
	if err != nil {
		return nil, err
	}

	m := map[string]FileInfoDto{}
	for _, file := range files {
		key := strings.TrimLeft(strings.Replace(file.Path, dir, "", 1), "/")
		m[key] = file
	}
	return m, nil
}

func getFileInfos(dir string, isRecursion bool) ([]FileInfoDto, error) {
	files, err := GetDirAllFiles(dir)
	if err != nil {
		return nil, err
	}

	resFiles := make([]FileInfoDto, 0)
	for _, f := range files {
		if f.IsDir() && isRecursion {
			dirFiles, err := getFileInfos(f.Path, isRecursion)
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
