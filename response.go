package bdpan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func NewRespError(r *http.Response) error {
	return NewErrorResponse(r).Err()
}

func NewErrorResponse(r *http.Response) *ErrorResponse {
	errResp := &ErrorResponse{}
	httpResponseToInterface(r, errResp)
	return errResp
}

type ErrorResponse struct {
	Errno            int32  `json:"errno,omitempty"`
	Erro             string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCode        int    `json:"error_code"`
	ErrorMsg         string `json:"error_msg"`
	r                *http.Response
}

func (e ErrorResponse) Error() string {
	return e.Err().Error()
}

func (e ErrorResponse) Print() {
	fmt.Fprintf(os.Stderr, "Error Code: %d\n", e.r.StatusCode)
	fmt.Fprintf(os.Stderr, "Error: %s\n", e.Erro)
	fmt.Fprintf(os.Stderr, "Error Desc: %s\n", e.ErrorDescription)
}

func (e ErrorResponse) Err() error {
	// https://pan.baidu.com/union/doc/okumlx17r
	// https://openauth.baidu.com/doc/appendix.html#_4-openapi%E9%94%99%E8%AF%AF%E7%A0%81%E5%88%97%E8%A1%A8

	if e.ErrorCode > 0 {
		return fmt.Errorf("%d[%s]", e.ErrorCode, e.ErrorMsg)
	} else if e.Erro != "" {
		return fmt.Errorf("%s[%s]", e.Erro, e.ErrorDescription)
	} else {
		switch e.Errno {
		case -9:
			return ErrPathNotFound
		case -6:
			return ErrAccessFail
		case 2:
			return ErrParamError
		case 6:
			return ErrUserNoUse
		case 111:
			return ErrAccessTokenFail
		case 31034, 9013, 9019:
			return ErrApiFrequent
		default:
			return fmt.Errorf("未知错误: %d", e.Errno)
		}
	}
}

type Response struct {
	ErrorResponse
}

func (r Response) IsError() bool {
	return r.Errno != 0
}

type FileListResponse struct {
	Response
	GuidInfo string         `json:"guid_info"`
	Errmsg   string         `json:"errmsg"`
	List     []*FileInfoDto `json:"list"`
}

func (f FileListResponse) Print() {
	if f.IsError() {
		fmt.Println(f.Err())
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
