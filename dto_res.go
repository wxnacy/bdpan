package bdpan

import (
	"bdpan/openapi"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/mattn/go-runewidth"
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
		switch f.GetCategory() {
		case "其他":
			return "文件"
		default:
			return f.GetCategory()
		}
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

func (f FileInfoDto) PrettyPrint() {
	tpl := `
------------ {{.Path}} ---------------
    FSID: {{.FSID}}
    Name: {{.GetFilename}}
Filetype: {{.GetFileTypeIcon}} {{.GetFileType}}
    Size: {{.GetSize}}
    Path: {{.Path}}
     MD5: {{.MD5}}
   CTime: {{.GetServerCTime}}
   MTime: {{.GetServerMTime}}
`
	tmpl, _ := template.New("FileInfoDtoPrettyPrint").Parse(tpl)
	_ = tmpl.Execute(os.Stdout, f)
}

func PrintFileInfoList(files []*FileInfoDto) {
	fmt.Println()
	idMaxLen := len("fsid")
	filenameMaxLen := len("name")
	sizeLen := len("Size")
	for _, f := range files {
		var length int
		length = runewidth.StringWidth(f.GetFilename())
		if length > filenameMaxLen {
			filenameMaxLen = length
		}
		length = len(strconv.Itoa(int(f.FSID)))
		if length > idMaxLen {
			idMaxLen = length
		}
		length = len(gotool.FormatSize(int64(f.Size)))
		if length > sizeLen {
			sizeLen = length
		}
	}
	idFmt := fmt.Sprintf("%%%ds", idMaxLen+1)
	sizeFmt := fmt.Sprintf(" %%-%ds", sizeLen+1)
	format := fmt.Sprintf("%s %%s %%s %-s %%-19s %%-19s\n", idFmt, sizeFmt)
	fmt.Printf(
		format,
		"FSID",
		runewidth.FillRight("Name", filenameMaxLen+1),
		"Filetype",
		"Size",
		"CTime",
		"MTime",
	)
	for _, f := range files {
		fmt.Printf(
			format,
			strconv.Itoa(int(f.FSID)),
			runewidth.FillRight(f.GetFilename(), filenameMaxLen+1),
			runewidth.FillRight(f.GetFileType(), 8),
			gotool.FormatSize(int64(f.Size)),
			f.GetServerCTime(),
			f.GetServerMTime(),
		)
	}
	fmt.Printf("Total: %d\n", len(files))
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
