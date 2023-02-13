package main

import (
	"bdpan/common"
	"fmt"
)

func Login() {
	credential := Credential{}
	fmt.Println("请先完善秘钥信息")
	fmt.Print("App Id: ")
	fmt.Scanln(&credential.AppId)
	fmt.Print("App Key: ")
	fmt.Scanln(&credential.AppKey)
	fmt.Print("Secret Key: ")
	fmt.Scanln(&credential.SecretKey)
	fmt.Print("Sign Key: ")
	fmt.Scanln(&credential.SignKey)
	err = common.WriteInterfaceToFile(credentialsPath, credential)
	if err != nil {
		panic(err)
	}

}
