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

func init() {
	_credentailMap = make(map[string]*Credential, 0)
}

var (
	conifg_dir, _   = common.ExpandUser("~/.config/bdpan")
	keyPath         = joinConfigPath("key")
	credentialsPath = joinConfigPath("credentials")
	tokenPath       = joinConfigPath("access_token")

	_credentailMap map[string]*Credential
	_config        *Config
)

func defaultCredentail() (*Credential, error) {
	items, err := GetCredentails()
	if err != nil {
		return nil, err
	}

	return items[0], nil
}

func GetConfigCredentail() (*Credential, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}
	fmt.Println(config.LoginAppId)
	return GetCredentail(config.LoginAppId)
}

func GetConfig() (*Config, error) {
	if _config != nil {
		fmt.Println("Get config from cache")
		return _config, nil
	}

	c, err := defaultCredentail()
	if err != nil {
		return nil, err
	}

	_config = NewConfig(c.AppId)
	fmt.Println("Get config from new")
	return _config, nil
}

// func defaultAccessToken() (*AccessToken, error) {
// items, err := GetCredentails()
// if err != nil {
// return nil, err
// }
// return items[0]., nil
// }

func initConfigDir() error {
	return os.MkdirAll(conifg_dir, common.PermDir)
}

func initConfig() error {
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

func saveCredentail(credentials []_Credential, c Credential) error {
	credentials = append(credentials, *newCredentail(c))
	return common.WriteInterfaceToFile(credentialsPath, credentials)
}

func GetCredentail(appId string) (*Credential, error) {
	if c, ok := _credentailMap[appId]; ok {
		fmt.Println("Get credentail from cache")
		return c, nil
	}
	credentials, err := GetCredentails()
	if err != nil {
		return nil, err
	}
	for _, c := range credentials {
		if c.AppId == appId {
			_credentailMap[appId] = c
			fmt.Println("Get credentail from file")
			return c, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("AppId %s credentail not found", appId))
}

func GetCredentails() ([]*Credential, error) {
	credentials := make([]_Credential, 0)
	err = common.ReadFileToInterface(credentialsPath, &credentials)
	if err != nil {
		return nil, err
	}
	res := make([]*Credential, 0)
	for _, c := range credentials {
		cre, err := c.GetCredentail()
		if err != nil {
			return nil, err
		}
		res = append(res, cre)
	}

	return res, nil
}

func AddCredentail(c Credential) error {
	credentials := make([]_Credential, 0)
	err = common.ReadFileToInterface(credentialsPath, &credentials)
	if err != nil {
		// TODO: 增加错误日志
		credentials = make([]_Credential, 0)
	}
	return saveCredentail(credentials, c)
}

func saveAccessToken(appId string, t AccessToken) error {
	m, err := common.ReadFileToMap(tokenPath)
	if err != nil {
		m = make(map[string]interface{}, 0)
	}
	tokenStr, err := encryptInterfaceToHex(t)
	if err != nil {
		return err
	}
	m[appId] = tokenStr

	err = common.WriteMapToFile(tokenPath, m)
	if err != nil {
		return err
	}
	return nil
}
