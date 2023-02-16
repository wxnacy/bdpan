package main

import (
	"fmt"
)

type _Credential struct {
	AppId       string `json:"app_id,omitempty"`
	Credentail  string `json:"credentail,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
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

func (c _Credential) GetAccessToken() (*AccessToken, error) {
	res := &AccessToken{}
	err = decryptHexToInterface(c.AccessToken, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *_Credential) SetAccessToken(a AccessToken) error {

	str, err := encryptInterfaceToHex(a)
	if err != nil {
		return err
	}
	c.AccessToken = str
	return nil
}

type Credential struct {
	AppId     string `json:"app_id,omitempty"`
	AppKey    string `json:"app_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
	SignKey   string `json:"sign_key,omitempty"`
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
