package auth

import (
	"context"

	"github.com/wxnacy/bdpan"
	"github.com/wxnacy/bdpan/response"
)

// 获取配额
func GetQuota(accessToken string) (*GetQuotaRes, error) {
	_, r, _ := bdpan.GetClient().
		UserinfoApi.
		Apiquota(context.Background()).
		AccessToken(accessToken).
		Checkexpire(1).
		Checkfree(1).
		Execute()
	return response.ToInterface[GetQuotaRes](r)
}

// 获取用户信息
func GetUserInfo(accessToken string) (*GetUserInfoRes, error) {
	_, r, _ := bdpan.GetClient().
		UserinfoApi.
		Xpannasuinfo(context.Background()).
		AccessToken(accessToken).
		Execute()
	return response.ToInterface[GetUserInfoRes](r)
}

// 通过 code 获取用户身份 token
func GetDeviceToken(appKey, secretKey, code string) (*GetDeviceTokenRes, error) {
	_, r, _ := bdpan.GetClient().
		AuthApi.
		OauthTokenDeviceToken(context.Background()).
		Code(code).
		ClientId(appKey).
		ClientSecret(secretKey).
		Execute()
	return response.ToInterface[GetDeviceTokenRes](r)
}

// 获取登录使用 code 二维码
func GetDeviceCode(appKey, scope string) (*GetDeviceCodeRes, error) {
	// scope := "basic,netdisk" // string
	_, r, _ := bdpan.GetClient().
		AuthApi.
		OauthTokenDeviceCode(context.Background()).
		ClientId(appKey).
		Scope(scope).
		Execute()
	return response.ToInterface[GetDeviceCodeRes](r)
}
