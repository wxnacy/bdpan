package file

import (
	"encoding/json"
	"fmt"

	"github.com/wxnacy/bdpan"
)

func NewGetFileListReq() *GetFileListReq {
	return &GetFileListReq{
		Dir:   "/",
		Web:   1,
		Start: 0,
		Limit: 1000,
		Order: "name",
	}
}

type GetFileListReq struct {
	Dir   string
	Web   int
	Page  int
	Start int
	Limit int32
	// 排序字段：默认为name；
	// time表示先按文件类型排序，后按修改时间排序；
	// name表示先按文件类型排序，后按文件名称排序；
	// size表示先按文件类型排序，后按文件大小排序。
	Order string
	// 默认为升序，设置为1实现降序 （注：排序的对象是当前目录下所有文件，不是当前分页下的文件）
	Desc int32
}

type GetFileListRes struct {
	GuidInfo string               `json:"guid_info"`
	Errmsg   string               `json:"errmsg"`
	List     []*bdpan.FileInfoDto `json:"list"`
}

func NewGetFileInfoReq(fsid uint64) *GetFileInfoReq {
	return &GetFileInfoReq{
		FSID: fsid,
	}
}

type GetFileInfoReq struct {
	Dlink int
	FSID  uint64
}

type GetFileInfoRes struct {
	bdpan.FileInfoDto
}

func NewBatchGetFileListReq(fsid uint64) *BatchGetFileInfoReq {
	return &BatchGetFileInfoReq{
		FSIds: []uint64{fsid},
	}
}

type BatchGetFileInfoReq struct {
	Dlink int
	FSIds []uint64
}

func (r *BatchGetFileInfoReq) AppendFSID(fsid uint64) *BatchGetFileInfoReq {
	r.FSIds = append(r.FSIds, fsid)
	return r
}

func (r *BatchGetFileInfoReq) GetFSIDString() string {
	bytesData, _ := json.Marshal(r.FSIds)
	return string(bytesData)
}

func (r *BatchGetFileInfoReq) GetDlink() string {
	return fmt.Sprintf("%d", r.Dlink)
}

type BatchGetFileInfoRes struct {
	List []*bdpan.FileInfoDto `json:"list"`
}
