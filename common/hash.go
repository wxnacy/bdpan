package common

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Md5File(path string) (string, error) {
	h := md5.New()
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", nil
	}
	h.Write(bytes)
	return hex.EncodeToString(h.Sum(nil)), nil
}
