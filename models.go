package main

import (
	"bdpan/common"
	"encoding/hex"
	"encoding/json"
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
	str, err := hex.DecodeString(c.Credentail)
	if err != nil {
		return nil, err
	}
	bytes, err := decrypt(str)
	if err != nil {
		return nil, err
	}
	res := &Credential{}
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *_Credential) SetCredentail(cre Credential) error {
	str, err := common.ToMapString(cre)
	if err != nil {
		return err
	}
	bytes, err := encrypt([]byte(str))
	if err != nil {
		return err
	}
	c.Credentail = hex.EncodeToString(bytes)
	return nil
}

type Credential struct {
	AppId     string `json:"app_id,omitempty"`
	AppKey    string `json:"app_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
	SignKey   string `json:"sign_key,omitempty"`
}

func (c Credential) Encrypt() string {
	str, _ := common.ToMapString(c)
	return common.Md5(str)
}

type AccessToken struct {
	Response
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
