package common

import (
	"github.com/wxnacy/gotool"
)

func Md5(str string) string {
	return gotool.Md5(str)
}

func Md5File(path string) (string, error) {
	return gotool.Md5File(path)
}
