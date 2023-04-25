package bdpan

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

const (
	DefaultUploadDir = "/apps/bdpan/"
)

var (
	// 默认下载目录
	DefaultDownloadDir, _ = homedir.Expand("~/Downloads")
	// 存储目录
	stoageDir, _ = homedir.Expand("~/.local/share/bdpan")
	// 配置目录
	configDir, _ = homedir.Expand("~/.config/bdpan")
	// 缓存目录
	cacheDir, _ = homedir.Expand("~/.cache/bdpan")
)

func JoinStoage(elem ...string) string {
	return join(stoageDir, elem...)
}

func join(root string, elem ...string) string {
	elem = append([]string{root}, elem...)
	return filepath.Join(elem...)
}
