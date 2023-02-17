package main

import (
	"fmt"
	"os"
)

func buildCredentail(arg LoginArg) Credential {
	appId := *arg.AppId

	credential := Credential{}
	fmt.Println("请先完善秘钥信息")
	if appId == "" {

		fmt.Print("App Id: ")
		fmt.Scanln(&credential.AppId)
	} else {
		credential.AppId = appId
	}
	fmt.Print("App Key: ")
	fmt.Scanln(&credential.AppKey)
	fmt.Print("Secret Key: ")
	fmt.Scanln(&credential.SecretKey)
	fmt.Print("Sign Key: ")
	fmt.Scanln(&credential.SignKey)
	return credential
}

func Login(arg LoginArg) {
	appId := *arg.AppId
	// var cres []*Credential
	if appId == "" {
		_, err = GetCredentails()
	} else {
		_, err = GetCredentail(appId)

	}
	if err != nil {
		credential := buildCredentail(arg)
		// credential.AppId = "1"
		// credential.AppKey = "1"
		// credential.SecretKey = "1"
		// credential.SignKey = "1"
		err := AddCredentail(credential)
		if err != nil {
			fmt.Fprintf(os.Stderr, "登录失败 %s\n", err.Error())
		}
	}

	// items, err := GetCredentails()
	// if err != nil {
	// panic(err)
	// }
	// for _, v := range items {
	// fmt.Println(*v)
	// }
	var c *Credential
	if appId != "" {
		c, err = GetCredentail(appId)
	} else {
		c, err = defaultCredentail()
	}
	if err != nil {
		panic(err)
	}
	t := &AccessToken{}
	t.AccessToken = "1"
	t.RefreshToken = "1"
	saveAccessToken(c.AppId, *t)
	fmt.Println(c)
	fmt.Println(c.GetAccessToken())
	fmt.Println(c.GetAccessToken())
	fmt.Println(c.GetAccessToken())

}
