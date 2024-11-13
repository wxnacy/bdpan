package file

import (
	"context"
	"fmt"
	"strconv"

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

// 获取文件详情
func GetFileInfo(accessToken string, req *GetFileInfoReq) (*GetFileInfoRes, error) {
	batchReq := NewBatchGetFileListReq(req.FSID)
	batchReq.Dlink = req.Dlink
	res, err := BatchGetFileInfo(accessToken, batchReq)
	if err != nil {
		return nil, err
	}
	return &GetFileInfoRes{
		FileInfoDto: *res.List[0],
	}, nil
}

// 批量获取文件详情
func BatchGetFileInfo(accessToken string, req *BatchGetFileInfoReq) (*BatchGetFileInfoRes, error) {
	_, r, _ := bdpan.GetClient().
		MultimediafileApi.Xpanmultimediafilemetas(
		context.Background()).AccessToken(accessToken).
		Dlink(req.GetDlink()).
		Fsids(req.GetFSIDString()).
		Execute()
	return response.ToInterface[BatchGetFileInfoRes](r)
}

// 搜索文件
func SearchFile(accessToken string, req *SearchFileReq) (*SearchFileRes, error) {
	_, r, _ := bdpan.GetClient().
		FileinfoApi.Xpanfilesearch(context.Background()).
		AccessToken(accessToken).
		Key(req.Key).
		Recursion(strconv.Itoa(req.Recursion)).
		Dir(req.Dir).
		Execute()
	return response.ToInterface[SearchFileRes](r)
}
