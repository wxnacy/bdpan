package file

import "github.com/wxnacy/bdpan"

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
	bdpan.FileInfoDto
	GuidInfo string               `json:"guid_info"`
	Errmsg   string               `json:"errmsg"`
	List     []*bdpan.FileInfoDto `json:"list"`
}
