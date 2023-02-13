package main

import (
	"bdpan/common"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	conifg_dir, _   = common.ExpandUser("~/.config/bdpan")
	keyPath         = joinConfigPath("key")
	credentialsPath = joinConfigPath("credentials")
)

func initConfigDir() error {
	return os.MkdirAll(conifg_dir, common.PermDir)
}

func initCryptoKey() error {
	info, err := os.Stat(keyPath)
	if err == nil {
		if !info.IsDir() {
			return nil
		}
		os.RemoveAll(keyPath)
	}

	key := common.Md5(strconv.Itoa(int(time.Now().Unix())))
	return os.WriteFile(keyPath, []byte(key), common.PermFile)
}

func joinConfigPath(name string) string {
	return filepath.Join(conifg_dir, name)
}

func GetKey() ([]byte, error) {
	return os.ReadFile(keyPath)
}

func saveCredentail(c Credential) error {
	credentials := make([]_Credential, 0)
	err = common.ReadFileToInterface(credentialsPath, &credentials)
	if err != nil {
		return err
	}
	credentials = append(credentials, *newCredentail(c))
	return common.WriteInterfaceToFile(credentialsPath, credentials)
}

func getCredentail(appId string) (*Credential, error) {
	credentials := make([]_Credential, 0)
	err = common.ReadFileToInterface(credentialsPath, &credentials)
	if err != nil {
		return nil, err
	}
	for _, c := range credentials {
		if c.AppId == appId {
			return c.GetCredentail()
		}
	}
	return nil, errors.New(fmt.Sprintf("AppId %s credentail not found", appId))
}

func GetCredentail() (*Credential, error) {

	return nil, nil
}

func AddCredentail(c Credential) {

}
