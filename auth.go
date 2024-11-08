package bdpan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/wxnacy/bdpan/openapi"
)

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
		// convertErrorResponse(r).Print()
		Log.Errorf("OauthTokenDeviceToken %v %v", err, convertErrorResponse(r))
		return err
	}
	token := AccessToken{}
	httpResponseToInterface(r, &token)
	saveAccessToken(credentail.AppId, token)
	return err
}

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

func UserInfo() (*UserInfoDto, error) {
	token, err := GetConfigAccessToken()
	if err != nil {
		return nil, err
	}
	resp, _, err := GetClient().UserinfoApi.Xpannasuinfo(context.Background()).AccessToken(
		token.AccessToken).Execute()
	if err != nil {
		return nil, err
	}
	return &UserInfoDto{
		Uinforesponse: resp,
	}, nil
}

func PanInfo() (*openapi.Quotaresponse, error) {
	token, err := GetConfigAccessToken()
	if err != nil {
		return nil, err
	}
	resp, _, err := GetClient().UserinfoApi.Apiquota(context.Background()).AccessToken(
		token.AccessToken).Checkexpire(1).Checkfree(1).Execute()
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
