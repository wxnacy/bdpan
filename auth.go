package main

import (
	sdk "bdpan/openapi"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func init() {
	// buildConfig()
	// buildApiClient()
	// buildAccessToken()
	// _config = config
}

var (
	apiClient *sdk.APIClient
	// config      AppConfig
	// _config     AppConfig
	_token      AccessToken
	CONFIG_PATH string = joinConfigPath("baidupan.json")
	TOKEN_PATH  string = joinConfigPath("baidupan_access_token.json")
)

type AppConfig struct {
	AppId     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	SecretKey string `json:"secret_key"`
	SignKey   string `json:"sign_key"`
}

// func buildConfig() {
// m, err := common.ReadFileToMap(CONFIG_PATH)
// if err != nil {
// panic(err)
// }

// b, err := json.MarshalIndent(m, "", "")
// if err != nil {
// panic(err)
// }
// err = json.Unmarshal(b, &config)
// if err != nil {
// panic(err)
// }
// }

// func buildAccessToken() {
// if !common.FileExists(TOKEN_PATH) {
// fmt.Fprintf(os.Stderr, "配置: %s 不存在\n", TOKEN_PATH)
// return
// }
// if _token != nil {
// return
// }
// _token = &AccessToken{}
// err := common.ReadFileToInterface(TOKEN_PATH, _token)
// if err != nil {
// panic(err)
// }
// }

func convertErrorResponse(r *http.Response) *ErrorResponse {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", r)
		panic(err)
	}
	var res ErrorResponse
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		fmt.Println(err)
		panic(err)
	}
	res.r = r
	return &res
}

// func getAccessTokenByDeviceCode(code string) (sdk.OauthTokenDeviceTokenResponse, *ErrorResponse) {
// resp, r, err := apiClient.AuthApi.OauthTokenDeviceToken(
// context.Background()).Code(
// code).ClientId(
// config.AppKey).ClientSecret(
// config.SecretKey).Execute()
// if err != nil {
// return resp, convertErrorResponse(r)
// }
// return resp, nil
// }

// func openDeviceCodeQrCode(appId string) string {
// scope := "basic,netdisk" // string
// c, err := GetCredentail(appId)
// if err != nil {
// return err
// }
// resp, r, err := apiClient.AuthApi.OauthTokenDeviceCode(
// context.Background()).ClientId(config.AppKey).Scope(scope).Execute()
// if err != nil {
// convertErrorResponse(r).Print()
// panic(err)
// }
// code := *resp.DeviceCode
// qrcode := *resp.QrcodeUrl
// fmt.Printf("DeviceCode: %s\n", code)
// cmd := exec.Command("open", qrcode)
// err = cmd.Run()
// if err != nil {
// panic(err)
// }
// return code
// }

func CreateAccessTokenByDeviceCode() error {
	fmt.Println("请求 device_code")
	scope := "basic,netdisk" // string
	credentail, err := GetConfigCredentail()
	if err != nil {
		return err
	}
	resp, r, err := GetClient().AuthApi.OauthTokenDeviceCode(
		context.Background()).ClientId(credentail.AppKey).Scope(scope).Execute()
	if err != nil {
		convertErrorResponse(r).Print()
		return err
	}
	code := *resp.DeviceCode
	qrcode := *resp.QrcodeUrl
	fmt.Printf("DeviceCode: %s\n", code)
	cmd := exec.Command("open", qrcode)
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("等待扫码")
	time.Sleep(time.Duration(10) * time.Second)
	fmt.Println("请求 access_token")
	// 请求 token
	_, r, err = GetClient().AuthApi.OauthTokenDeviceToken(
		context.Background()).Code(
		code).ClientId(
		credentail.AppKey).ClientSecret(
		credentail.SecretKey).Execute()
	if err != nil {
		convertErrorResponse(r).Print()
		return err
	}
	// fmt.Println(errRes)
	// for errRes != nil {
	// time.Sleep(time.Duration(2) * time.Second)
	// tokenRes, errRes = getAccessTokenByDeviceCode(code)
	// errRes.Print()
	// fmt.Fprintln(os.Stderr, "等待重试...")
	// }
	// fmt.Println(*tokenRes.AccessToken)
	// fmt.Println(*tokenRes.RefreshToken)
	// fmt.Println(*tokenRes.ExpiresIn)

	// var _token AccessToken
	// _token.AccessToken = *tokenRes.AccessToken
	// _token.RefreshToken = *tokenRes.RefreshToken
	// _token.ExpiresIn = *tokenRes.ExpiresIn
	token := AccessToken{}
	httpResponseToInterface(r, &token)
	saveAccessToken(credentail.AppId, token)
	return err
}

// func HttpResponseToAccessToken(r *http.Response, t *AccessToken) error {
// bodyBytes, err := ioutil.ReadAll(r.Body)
// if err != nil {
// return err
// }
// if err := json.Unmarshal(bodyBytes, t); err != nil {
// return err
// }
// t.RefreshTimestamp = time.Now().Unix()
// return nil
// }

func RefreshAccessToken() error {
	fmt.Println("开始刷新 access_token")
	credentail, err := GetConfigCredentail()
	if err != nil {
		return err
	}
	token, err := credentail.GetAccessToken()
	if err != nil {
		return err
	}
	if token == nil || token.AccessToken == "" {
		fmt.Println("初始 access_token 不存在，重新走申请流程")
		return CreateAccessTokenByDeviceCode()
	}
	fmt.Println("当前信息")
	token.Print()
	_, r, err := GetClient().AuthApi.OauthTokenRefreshToken(
		context.Background()).RefreshToken(token.RefreshToken).ClientId(
		credentail.AppKey).ClientSecret(credentail.SecretKey).Execute()
	if err != nil {
		errRes := convertErrorResponse(r)
		if errRes.ErrorDescription == "refresh token has been used" {
			fmt.Println("refresh_token 已被使用，重新走申请流程")
			return CreateAccessTokenByDeviceCode()
		}
		errRes.Print()
		return errors.New(errRes.ErrorDescription)
	}
	err = httpResponseToInterface(r, token)
	if err != nil {
		errRes := convertErrorResponse(r)
		errRes.Print()
		return errors.New(errRes.ErrorDescription)
	}
	token.Print()
	saveAccessToken(credentail.AppId, *token)
	fmt.Println("access_token 刷新完成")
	return nil
}

func ScheRefreshAccessToken() error {
	credentail, err := GetConfigCredentail()
	if err != nil {
		return err
	}
	token, err := credentail.GetAccessToken()
	if err != nil {
		return RefreshAccessToken()
	}
	expiresSecond := 7 * 24 * 3600
	refreshTime := time.Unix(token.RefreshTimestamp, 0)
	if time.Now().Sub(refreshTime).Seconds() < float64(expiresSecond) {
		fmt.Println("当前 access_token 已经是最新，无需刷新")
		return nil
	}

	return RefreshAccessToken()
}
