package bdpan

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/wxnacy/bdpan/common"
	"github.com/wxnacy/gotool"
)

var (
	keyPath         = joinConfigPath("key")
	credentialsPath = joinConfigPath("credentials")
	tokenPath       = joinConfigPath("access_token")
	configPath      = joinConfigPath("config.json")

	credentailMap = make(map[string]*Credential, 0)
	_config       *Config
)

func init() {
	var err error
	err = gotool.DirExistsOrCreate(stoageDir)
	panicErr(err)
	err = gotool.DirExistsOrCreate(configDir)
	panicErr(err)
	err = gotool.DirExistsOrCreate(cacheDir)
	panicErr(err)
	err = initCryptoKey()
	panicErr(err)
	initLogger()
}

// 默认鉴权账户
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
	return GetCredentail(config.LoginAppId)
}

func GetConfigAccessToken() (*AccessToken, error) {
	c, err := GetConfigCredentail()
	if err != nil {
		return nil, err
	}
	return c.GetAccessToken()
}

func GetConfig() (*Config, error) {
	if _config != nil {
		return _config, nil
	}

	if common.FileExists(configPath) {
		config := &Config{}
		err := common.ReadFileToInterface(configPath, config)
		if err != nil {
			return nil, err
		}
		_config = config
		return config, nil
	}

	c, err := defaultCredentail()
	if err != nil {
		return nil, err
	}

	_config = NewConfig(c.AppId)
	return _config, nil
}

// 初始化加密 key
func initCryptoKey() error {
	if gotool.FileExists(keyPath) {
		return nil
	}
	key := common.Md5(strconv.Itoa(int(time.Now().Unix())))
	return gotool.FileWriteWithInterface(keyPath, key)
}

func joinConfigPath(name string) string {
	return filepath.Join(configDir, name)
}

func GetKey() ([]byte, error) {
	return os.ReadFile(keyPath)
}

func saveCredentail(credentials []_Credential, c Credential) error {
	credentials = append(credentials, *newCredentail(c))
	return common.WriteInterfaceToFile(credentialsPath, credentials)
}

func GetCredentail(appId string) (*Credential, error) {
	if c, ok := credentailMap[appId]; ok {
		return c, nil
	}
	credentials, err := GetCredentails()
	if err != nil {
		return nil, err
	}
	for _, c := range credentials {
		if c.AppId == appId {
			credentailMap[appId] = c
			return c, nil
		}
	}
	return nil, fmt.Errorf("AppId %s credentail not found", appId)
}

func GetCredentails() ([]*Credential, error) {
	credentials := make([]_Credential, 0)
	err := common.ReadFileToInterface(credentialsPath, &credentials)
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
	err := common.ReadFileToInterface(credentialsPath, &credentials)
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
