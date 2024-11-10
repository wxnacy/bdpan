package auth

type GetDeviceCodeRes struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	QrcodeURL       string `json:"qrcode_url"`
	ExpiresIn       int64  `json:"expires_in"`
	Interval        int64  `json:"interval"`
}

type GetDeviceTokenRes struct {
	ExpiresIn     int32  `json:"expires_in,omitempty"`
	RefreshToken  string `json:"refresh_token,omitempty"`
	AccessToken   string `json:"access_token,omitempty"`
	SessionSecret string `json:"session_secret,omitempty"`
	SessionKey    string `json:"session_key,omitempty"`
	Scope         string `json:"scope,omitempty"`
}

type GetUserInfoRes struct {
	Uk          int    `json:"uk,omitempty"`
	RequestId   string `json:"request_id,omitempty"`
	AvatarUrl   string `json:"avatar_url,omitempty"`
	BaiduName   string `json:"baidu_name,omitempty"`
	NetdiskName string `json:"netdisk_name,omitempty"`
	VipType     int32  `json:"vip_type,omitempty"`
}

func (u GetUserInfoRes) GetVipName() string {
	switch u.VipType {
	case 0:
		return "普通用户"
	case 1:
		return "普通会员"
	case 2:
		return "超级会员"
	}
	return "未知身份"
}

type GetQuotaRes struct {
	Total     int64 `json:"total,omitempty"`
	Free      int64 `json:"free,omitempty"`
	RequestId int64 `json:"request_id,omitempty"`
	Expire    bool  `json:"expire,omitempty"`
	Used      int64 `json:"used,omitempty"`
}
