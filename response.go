package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type UploadDirResponse struct {
	SuccessCount int
	FailedCount  int
}

type FileListResponse struct {
	Errno    int32          `json:"errno"`
	GuidInfo string         `json:"guid_info"`
	Errmsg   string         `json:"errmsg"`
	List     []*FileInfoDto `json:"list"`
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
	FileInfoDto
	Errno int32 `json:"errno"`
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
