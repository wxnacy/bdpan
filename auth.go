package main

import (
	// "context"
	// "encoding/json"
	// "fmt"
	// "io/ioutil"
	// "net/http"
	// "os"
	// "os/exec"
	// "time"
	sdk "bdpan/openapi"
	// "website/common"
)

func init() {
	// buildConfig()
	// buildApiClient()
	// buildAccessToken()
	_client = apiClient
	_config = config
}

var (
	apiClient   *sdk.APIClient
	_client     *sdk.APIClient
	config      AppConfig
	_config     AppConfig
	_token      *AccessToken
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
// err := common.ReadFileToModel(TOKEN_PATH, _token)
// if err != nil {
// panic(err)
// }
// }

// func buildApiClient() {
// configuration := sdk.NewConfiguration()
// apiClient = sdk.NewAPIClient(configuration)
// }

// func convertErrorResponse(r *http.Response) *ErrorResponse {
// bodyBytes, err := ioutil.ReadAll(r.Body)
// if err != nil {
// fmt.Fprintf(os.Stderr, "err: %v\n", r)
// panic(err)
// }
// var res ErrorResponse
// if err := json.Unmarshal(bodyBytes, &res); err != nil {
// fmt.Println(err)
// panic(err)
// }
// res.r = r
// return &res
// }

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

// func openDeviceCodeQrCode() string {
// scope := "basic,netdisk" // string
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

// func CreateAccessTokenByDeviceCode() {
// fmt.Println("请求 device_code")
// code := openDeviceCodeQrCode()
// fmt.Println("等待扫码")
// time.Sleep(time.Duration(10) * time.Second)
// fmt.Println("请求 access_token")
// // 请求 token
// _, r, err := apiClient.AuthApi.OauthTokenDeviceToken(
// context.Background()).Code(
// code).ClientId(
// config.AppKey).ClientSecret(
// config.SecretKey).Execute()
// if err != nil {
// convertErrorResponse(r).Print()
// return
// }
// // fmt.Println(errRes)
// // for errRes != nil {
// // time.Sleep(time.Duration(2) * time.Second)
// // tokenRes, errRes = getAccessTokenByDeviceCode(code)
// // errRes.Print()
// // fmt.Fprintln(os.Stderr, "等待重试...")
// // }
// // fmt.Println(*tokenRes.AccessToken)
// // fmt.Println(*tokenRes.RefreshToken)
// // fmt.Println(*tokenRes.ExpiresIn)

// // var _token AccessToken
// // _token.AccessToken = *tokenRes.AccessToken
// // _token.RefreshToken = *tokenRes.RefreshToken
// // _token.ExpiresIn = *tokenRes.ExpiresIn
// HttpResponseToAccessToken(r, _token)
// saveAccessToken(*_token)
// }

// func saveAccessToken(t AccessToken) {
// tokenMap, err := t.ToMap()
// if err != nil {
// panic(err)
// }
// err = common.WriteMapToFile(TOKEN_PATH, tokenMap)
// if err != nil {
// panic(err)
// }

// }

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

// func RefreshAccessToken() {
// fmt.Println("开始刷新 access_token")
// if _token == nil || _token.AccessToken == "" {
// fmt.Println("初始 access_token 不存在，重新走申请流程")
// CreateAccessTokenByDeviceCode()
// return
// }
// fmt.Println("当前信息")
// _token.Print()
// _, r, err := apiClient.AuthApi.OauthTokenRefreshToken(context.Background()).RefreshToken(_token.RefreshToken).ClientId(config.AppKey).ClientSecret(config.SecretKey).Execute()
// if err != nil {
// errRes := convertErrorResponse(r)
// if errRes.ErrorDescription == "refresh token has been used" {
// fmt.Println("refresh_token 已被使用，重新走申请流程")
// CreateAccessTokenByDeviceCode()
// return
// }
// errRes.Print()
// return
// }
// err = HttpResponseToAccessToken(r, _token)
// if err != nil {
// convertErrorResponse(r).Print()
// panic(err)
// }
// _token.Print()
// saveAccessToken(*_token)
// fmt.Println("access_token 刷新完成")
// }

// func ScheRefreshAccessToken() {
// expiresSecond := 7 * 24 * 3600
// refreshTime := time.Unix(_token.RefreshTimestamp, 0)
// if time.Now().Sub(refreshTime).Seconds() < float64(expiresSecond) {
// fmt.Println("当前 access_token 已经是最新，无需刷新")
// return
// }

// RefreshAccessToken()
// }
