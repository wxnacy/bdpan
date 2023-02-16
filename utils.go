package main

import (
	"bdpan/common"
	"encoding/hex"
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

func encryptInterfaceToHex(i interface{}) (string, error) {
	str, err := common.ToMapString(i)
	if err != nil {
		return "", err
	}
	bytes, err := encrypt([]byte(str))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func decrypt(src []byte) ([]byte, error) {
	key, err := GetKey()
	if err != nil {
		return nil, err
	}
	return common.AesDecrypt(src, key)
}

func decryptHexToInterface(src string, i interface{}) error {
	str, err := hex.DecodeString(src)
	if err != nil {
		return err
	}
	bytes, err := decrypt(str)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, i)
	if err != nil {
		return err
	}
	return nil
}
