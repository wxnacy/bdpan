package bdpan

import (
	"bdpan/openapi"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
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
	FSID           uint64            `json:"fs_id"`
	Path           string            `json:"path"`
	Size           int               `json:"size"`
	FileType       int               `json:"isdir"`
	Filename       string            `json:"filename"`
	ServerFilename string            `json:"server_filename"`
	Category       int               `json:"category"`
	Dlink          string            `json:"dlink"`
	MD5            string            `json:"md5"`
	Thumbs         map[string]string `json:"thumbs"`
	ServerCTime    int64             `json:"server_ctime"`
	ServerMTime    int64             `json:"server_mtime"`
}

func (f FileInfoDto) GetFilename() string {
	if f.ServerFilename != "" {
		return f.ServerFilename
	}
	return f.Filename
}

func (f FileInfoDto) GetSize() string {
	return gotool.FormatSize(int64(f.Size))

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

func (f FileInfoDto) GetFileTypeIcon() string {
	if f.IsDir() {
		return ""
	}
	icon, ok := GetIconByPath(f.GetFilename())
	if !ok {
		icon = GetDefaultFileIcon()
	}
	return icon.Icon
}

func (f FileInfoDto) GetFileTypeEmoji() string {
	if f.IsDir() {
		// 🗂️
		return "\U0001f5c2"
	} else {
		switch f.Category {
		case 1:
			// 📹
			return "\U0001f4f9"
		case 2:
			// 🎵
			return "\U0001f3b5"
		case 3:
			// 🖼️
			return "\U0001f5bc"
		case 4:
			// 📄
			return "\U0001f4c4"
		case 5:
			// 🚀
			return "\U0001f680"
		case 6:
			// 其他 🤷
			return "\U0001f937"
		case 7:
			// 种子 🤷
			return "\U0001f937"
		}
		// 🤷
		return "\U0001f937"
	}
}

func (f FileInfoDto) GetFileType() string {
	if f.IsDir() {
		return "文件夹"
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
	fmt.Printf("%d\t%s\t%s\n", f.FSID, f.GetFilename(), f.GetSize())
}

func (f FileInfoDto) PrettyPrint() {
	fields := [][]string{
		[]string{"FSID", strconv.Itoa(int(f.FSID))},
		[]string{"Name", f.GetFilename()},
		[]string{"Filetype", f.GetFileTypeIcon() + " " + f.GetFileType()},
		[]string{"Size", f.GetSize()},
		[]string{"Path", f.Path},
		[]string{"MD5", f.MD5},
		[]string{"Dlink", f.Dlink},
		[]string{"Ctime", f.GetServerCTime()},
		[]string{"Mtime", f.GetServerMTime()},
	}
	for _, f := range fields {
		fmt.Printf("%8s: %s\n", f[0], f[1])
	}
	fmt.Println("")
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
