package bdpan

func NewFileSortSlice(files []*FileInfoDto) fileSortSlice {
	return fileSortSlice(files)
}

type fileSortSlice []*FileInfoDto

func (s fileSortSlice) Len() int {
	return len(s)
}

func (s fileSortSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// func (s SyncModelSlice) Less(i, j int) bool { return s[i].CreateTime.Before(s[j].CreateTime) }
