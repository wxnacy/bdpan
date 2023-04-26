package bdpan

import (
	"bdpan/openapi"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/wxnacy/gotool"
)

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
	LocalCTime     int64             `json:"local_ctime"`
	LocalMTime     int64             `json:"local_mtime"`
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

func (f FileInfoDto) GetLocalCTime() string {
	return f.formatTime(f.LocalCTime)
}

func (f FileInfoDto) GetLocalMTime() string {
	return f.formatTime(f.LocalMTime)
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
	fmt.Printf("------------ %s ---------------\n", f.Path)
	fmt.Print(f.GetPretty())
}
func (f FileInfoDto) GetPretty() string {
	tpl := `    FSID: {{.FSID}}
    Name: {{.GetFilename}}
Filetype: {{.GetFileTypeIcon}} {{.GetFileType}}
    Size: {{.GetSize}}
    Path: {{.Path}}
     MD5: {{.MD5}}
   CTime: {{.GetServerCTime}}
   MTime: {{.GetServerMTime}}
  LCTime: {{.GetLocalCTime}}
  LMTime: {{.GetLocalMTime}}
`
	tmpl, _ := template.New("FileInfoDtoPrettyPrint").Parse(tpl)
	buf := new(strings.Builder)
	_ = tmpl.Execute(buf, f)
	return buf.String()
}

func (f FileInfoDto) BuildPrintData() []PrettyData {
	var data = make([]PrettyData, 0)
	data = append(data, PrettyData{
		Name:       "FSID",
		Value:      strconv.Itoa(int(f.FSID)),
		IsFillLeft: true})
	data = append(data, PrettyData{Name: "Name", Value: f.GetFilename()})
	data = append(data, PrettyData{Name: "Filetype", Value: f.GetFileType()})
	data = append(data, PrettyData{Name: "Size", Value: f.GetSize()})
	// data = append(data, PrettyData{Name: "CTime", Value: f.GetServerCTime()})
	data = append(data, PrettyData{Name: "MTime", Value: f.GetServerMTime()})
	// data = append(data, PrettyData{Name: "LCTime", Value: f.GetLocalCTime()})
	data = append(data, PrettyData{Name: "LMTime", Value: f.GetLocalMTime()})
	return data
}

func PrintFileInfoList(files []*FileInfoDto) {
	// return
	prettyList := make([]Pretty, 0)
	for _, f := range files {
		prettyList = append(prettyList, f)
	}
	PrettyPrintList(PrettyList(prettyList))
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
