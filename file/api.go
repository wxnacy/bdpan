package file

import (
	"context"
	"fmt"

	"github.com/wxnacy/bdpan"
	"github.com/wxnacy/bdpan/response"
)

// 获取文件列表
func GetFileList(accessToken string, req *GetFileListReq) (*GetFileListRes, error) {
	_, r, _ := bdpan.GetClient().
		FileinfoApi.Xpanfilelist(context.Background()).
		AccessToken(accessToken).
		Dir(req.Dir).
		Web(fmt.Sprintf("%d", req.Web)).
		Start(fmt.Sprintf("%d", req.Start)).
		Order(req.Order).
		Desc(req.Desc).
		Limit(req.Limit).
		Execute()
	return response.ToInterface[GetFileListRes](r)
}
