package main

import (
	"bdpan/common"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func httpResponseToInterface(r *http.Response, i interface{}) error {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bodyBytes, i); err != nil {
		return err
	}
	return nil
}

func encrypt(src []byte) ([]byte, error) {
	key, err := GetKey()
	if err != nil {
		return nil, err
	}
	return common.AesEncrypt(src, key)
}

func decrypt(src []byte) ([]byte, error) {
	key, err := GetKey()
	if err != nil {
		return nil, err
	}
	return common.AesDecrypt(src, key)
}
