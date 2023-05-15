package bdpan

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

const (
	DefaultUploadDir = "/apps/bdpan/"
	EnvDownloadDir   = "BDPAN_DOWNLOAD_DIR"
)

var (
	// 当前目录
	PwdDir, _ = os.Getwd()
	// 默认下载目录
	DefaultDownloadDir, _ = homedir.Expand("~/Downloads")
	// 存储目录
	stoageDir, _ = homedir.Expand("~/.local/share/bdpan")
	// 配置目录
	configDir, _ = homedir.Expand("~/.config/bdpan")
	// 缓存目录
	cacheDir, _ = homedir.Expand("~/.cache/bdpan")
	// 日志地址
	logPath = JoinCache("bdpan.log")
	// 同步存储地址
	syncPath = JoinStoage("sync.json")
)

func GetDefaultDownloadDir() string {
	dir := os.Getenv(EnvDownloadDir)
	if dir == "" {
		dir = PwdDir
	}
	return dir
}

func JoinCache(elem ...string) string {
	return join(cacheDir, elem...)
}

func JoinDownload(elem ...string) string {
	return join(DefaultDownloadDir, elem...)
}

func JoinStoage(elem ...string) string {
	return join(stoageDir, elem...)
}

func join(root string, elem ...string) string {
	elem = append([]string{root}, elem...)
	return filepath.Join(elem...)
}
