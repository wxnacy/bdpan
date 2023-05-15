package bdpan

import "errors"

var ErrPathNotFound error = errors.New("文件或目录不存在")

var ErrAccessFail error = errors.New("身份验证失败")

var ErrParamError error = errors.New("参数错误")

var ErrUserNoUse error = errors.New("不允许接入用户数据")

var ErrAccessTokenFail error = errors.New("access token 失效")

var ErrApiFrequent error = errors.New("接口请求过于频繁，注意控制")

var ErrPathExists error = errors.New("文件已存在")
