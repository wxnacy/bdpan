package bdpan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Errno int32 `json:"errno"`
}

func (r Response) IsError() bool {
	return r.Errno > 0
}

func (r Response) Error() string {
	// https://pan.baidu.com/union/doc/okumlx17r
	switch r.Errno {
	case -6:
		return "身份验证失败"
	case 2:
		return "参数错误"
	case 6:
		return "不允许接入用户数据"
	case 111:
		return "access token 失效"
	case 31034:
		return "接口请求过于频繁，注意控制"
	}
	return fmt.Sprintf("未知错误: %d", r.Errno)
}

func (r Response) Err() error {
	return errors.New(r.Error())
}

// type UploadDirResponse struct {
// SuccessCount int
// FailedCount  int
// }

type FileListResponse struct {
	Response
	GuidInfo string         `json:"guid_info"`
	Errmsg   string         `json:"errmsg"`
	List     []*FileInfoDto `json:"list"`
}

func (f FileListResponse) Print() {
	if f.IsError() {
		fmt.Println(f.Error())
		return
	}
	PrintFileInfoList(f.List)
}

func NewFileListResponse(r *http.Response) (*FileListResponse, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	dto := &FileListResponse{}
	if err := json.Unmarshal(bodyBytes, dto); err != nil {
		return nil, err
	}
	if dto.IsError() {
		return nil, dto.Err()
	}
	return dto, nil
}

type FileListAllResponse struct {
	Errno   int32          `json:"errno"`
	Errmsg  string         `json:"errmsg"`
	HasMore int            `json:"has_more"`
	Cursor  int            `json:"cursor"`
	List    []*FileInfoDto `json:"list"`
}

type UploadFileResponse struct {
	Response
	FileInfoDto
}

// ****************************************
// FileManagerResponse
// ****************************************

type FileManagerInfo struct {
	Errno int32  `json:"errno,omitempty"`
	Path  string `json:"path,omitempty"`
}

type FileManagerResponse struct {
	Errno     int32             `json:"errno,omitempty"`
	Info      []FileManagerInfo `json:"info,omitempty"`
	RequestId int64             `json:"request_id,omitempty"`
	Taskid    int64             `json:"taskid,omitempty"`
}
