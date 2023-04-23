package bdpan

import (
	"encoding/json"
	"path/filepath"
)

var (
	iconExtMap map[string]Icon
)

func init() {
	// iconExtMap = make(map[string]Icon, 0)

	err := json.Unmarshal([]byte(iconText), &iconExtMap)
	if err != nil {
		Log.Error(err)
	}
}

type Icon struct {
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

func GetDefaultFileIcon() Icon {
	icon, _ := GetIconByExt("txt")
	return icon
}

func GetIconByExt(ext string) (Icon, bool) {
	icon, ok := iconExtMap[ext]
	if !ok {
		return Icon{}, false
	}
	return icon, true
}

func GetIconByPath(path string) (Icon, bool) {
	ext := filepath.Ext(path)
	if ext == "" {
		return Icon{}, false
	}
	return GetIconByExt(ext[1:])
}
