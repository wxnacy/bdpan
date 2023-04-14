package bdpan

import (
	"bdpan/common"
	"errors"
	"fmt"
)

type _Credential struct {
	AppId      string `json:"app_id,omitempty"`
	Credentail string `json:"credentail,omitempty"`
}

func newCredentail(c Credential) *_Credential {
	res := &_Credential{}
	res.AppId = c.AppId
	res.SetCredentail(c)
	return res
}

func (c _Credential) GetCredentail() (*Credential, error) {
	res := &Credential{}
	err = decryptHexToInterface(c.Credentail, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *_Credential) SetCredentail(cre Credential) error {

	str, err := encryptInterfaceToHex(cre)
	if err != nil {
		return err
	}
	c.Credentail = str
	return nil
}

type Credential struct {
	AppId       string `json:"app_id,omitempty"`
	AppKey      string `json:"app_key,omitempty"`
	SecretKey   string `json:"secret_key,omitempty"`
	SignKey     string `json:"sign_key,omitempty"`
	accessToken *AccessToken
}

func (c *Credential) GetAccessToken() (*AccessToken, error) {
	if c.accessToken != nil {
		// fmt.Println("return AccessToken from field")
		return c.accessToken, nil
	}
	m, err := common.ReadFileToMap(tokenPath)
	if err != nil {
		return nil, err
	}

	tokenStr, ok := m[c.AppId]
	if !ok {
		return nil, errors.New(fmt.Sprintf(
			"AppId: %s AccessToken not found", c.AppId))
	}

	token := &AccessToken{}
	err = decryptHexToInterface(tokenStr.(string), token)
	if err != nil {
		return nil, err
	}
	// fmt.Println("return AccessToken from file")
	c.accessToken = token
	return token, nil
}

type AccessToken struct {
	ExpiresIn        int32  `json:"expires_in,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	RefreshTimestamp int64  `json:"refresh_timestamp,omitempty"`
}

func (t AccessToken) Print() {
	fmt.Printf("AccessToken: %s\n", t.AccessToken)
	fmt.Printf("RefreshToken: %s\n", t.RefreshToken)
	fmt.Printf("RefreshTimestamp: %d\n", t.RefreshTimestamp)
}

type Config struct {
	LoginAppId string `json:"login_app_id,omitempty"`
}

func NewConfig(loginAppId string) *Config {
	return &Config{LoginAppId: loginAppId}
}
