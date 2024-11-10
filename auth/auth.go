package auth

// import (
// "context"

// "github.com/wxnacy/bdpan"
// "github.com/wxnacy/bdpan/response"
// )

// func NewAuth(credential *Credential) *Auth {
// return &Auth{
// credential: credential,
// }
// }

// type Auth struct {
// credential *Credential
// access     *Access
// }

// func (a *Auth) SetAccess(access *Access) *Auth {
// a.access = access
// return a
// }

// func (a *Auth) SetCredential(c *Credential) *Auth {
// a.credential = c
// return a
// }

// // 获取登录使用 code 二维码
// func (a *Auth) GetDeviceCode() (*GetDeviceCodeRes, error) {
// scope := "basic,netdisk" // string
// _, r, _ := bdpan.GetClient().
// AuthApi.
// OauthTokenDeviceCode(context.Background()).
// ClientId(a.credential.AppKey).
// Scope(scope).
// Execute()
// return response.ToInterface[GetDeviceCodeRes](r)
// }

// // 通过 code 获取用户身份 token
// func (a *Auth) GetDeviceToken(code string) (*GetDeviceTokenRes, error) {
// _, r, _ := bdpan.GetClient().
// AuthApi.
// OauthTokenDeviceToken(context.Background()).
// Code(code).
// ClientId(a.credential.AppKey).
// ClientSecret(a.credential.SecretKey).
// Execute()
// return response.ToInterface[GetDeviceTokenRes](r)
// }

// // 获取用户信息
// func (a *Auth) GetUserInfo() (*GetUserInfoRes, error) {
// _, r, _ := bdpan.GetClient().
// UserinfoApi.
// Xpannasuinfo(context.Background()).
// AccessToken(a.access.AccessToken).
// Execute()
// return response.ToInterface[GetUserInfoRes](r)
// }

// // 获取配额
// func (a *Auth) GetQuota() (*GetQuotaRes, error) {
// _, r, _ := bdpan.GetClient().
// UserinfoApi.
// Apiquota(context.Background()).
// AccessToken(a.access.AccessToken).
// Checkexpire(1).
// Checkfree(1).
// Execute()
// return response.ToInterface[GetQuotaRes](r)
// }
