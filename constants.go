package bdpan

import "github.com/mitchellh/go-homedir"

var (
	DefaultDownloadDir, _ = homedir.Expand("~/Downloads")
)
