package bdpan

func getDirFileInfoMap(dir string) (map[string]FileInfoDto, error) {
	files, err := GetDirAllFiles(dir)
	if err != nil {
		return nil, err
	}

	m := map[string]FileInfoDto{}
	for _, file := range files {
		m[file.GetFilename()] = *file
	}
	return m, nil
}
