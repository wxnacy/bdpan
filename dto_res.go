package bdpan

import (
	"bdpan/openapi"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/wxnacy/gotool"
)

func NewErrorResponse(r *http.Response) *ErrorResponse {
	errResp := &ErrorResponse{}
	httpResponseToInterface(r, errResp)
	return errResp
}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorCode        int    `json:"error_code"`
	ErrorMsg         string `json:"error_msg"`
	ErrorDescription string `json:"error_description"`
	r                *http.Response
}

func (e ErrorResponse) Err() error {
	if e.ErrorCode > 0 {
		return errors.New(fmt.Sprintf("%d[%s]", e.ErrorCode, e.ErrorMsg))
	}
	return errors.New(fmt.Sprintf("%s[%s]", e.Error, e.ErrorDescription))
}

func (e ErrorResponse) Print() {
	fmt.Fprintf(os.Stderr, "Error Code: %d\n", e.r.StatusCode)
	fmt.Fprintf(os.Stderr, "Error: %s\n", e.Error)
	fmt.Fprintf(os.Stderr, "Error Desc: %s\n", e.ErrorDescription)
}

type FileInfoDto struct {
	FSID           uint64 `json:"fs_id"`
	Path           string `json:"path"`
	Size           int    `json:"size"`
	FileType       int    `json:"isdir"`
	Filename       string `json:"filename"`
	ServerFilename string `json:"server_filename"`
	Category       int    `json:"category"`
	ServerCTime    int64  `json:"server_ctime"`
	ServerMTime    int64  `json:"server_mtime"`
	Dlink          string `json:"dlink"`
	MD5            string `json:"md5"`
}

func (f FileInfoDto) GetFilename() string {
	if f.ServerFilename != "" {
		return f.ServerFilename
	}
	return f.Filename
}

func (f FileInfoDto) formatTime(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

func (f FileInfoDto) GetServerCTime() string {
	return f.formatTime(f.ServerCTime)
}

func (f FileInfoDto) GetServerMTime() string {
	return f.formatTime(f.ServerMTime)
}

func (f FileInfoDto) IsDir() bool {
	if f.FileType == 1 {
		return true
	} else {
		return false
	}
}

func (f FileInfoDto) GetFileType() string {
	if f.IsDir() {
		return "目录"
	} else {
		return f.GetCategory()
	}
}

func (f FileInfoDto) GetCategory() string {
	// 文件类型，1 视频、2 音频、3 图片、4 文档、5 应用、6 其他、7 种子
	switch f.Category {
	case 1:
		return "视频"
	case 2:
		return "音频"
	case 3:
		return "图片"
	case 4:
		return "文档"
	case 5:
		return "应用"
	case 6:
		return "其他"
	case 7:
		return "种子"
	}
	return "未知"
}

func (f FileInfoDto) PrintOneLine() {
	// fmt.Printf("%d\t%s\t%s\t%d\n", f.FSID, f.MD5, f.GetFilename(), f.Size)
	fmt.Printf("%d\t%s\t%s\n", f.FSID, f.GetFilename(), gotool.FormatSize(int64(f.Size)))
}

type UserInfoDto struct {
	openapi.Uinforesponse
}

func (u UserInfoDto) GetVipName() string {
	switch u.GetVipType() {
	case 0:
		return "普通用户"
	case 1:
		return "普通会员"
	case 2:
		return "超级会员"
	}
	return "未知身份"
}

func printFileInfoList(files []*FileInfoDto) {
	for _, f := range files {
		f.PrintOneLine()
	}
	fmt.Printf("Total: %d\n", len(files))
}
