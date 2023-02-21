package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/wxnacy/gotool"
)

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	r                *http.Response
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
	Filename       string `json:"filename"`
	ServerFilename string `json:"server_filename"`
	Category       int    `json:"category"`
	Dlink          string `json:"dlink"`
	MD5            string `json:"md5"`
}

func (f FileInfoDto) GetFilename() string {
	if f.ServerFilename != "" {
		return f.ServerFilename
	}
	return f.Filename
}

func (f FileInfoDto) PrintOneLine() {
	// fmt.Printf("%d\t%s\t%s\t%d\n", f.FSID, f.MD5, f.GetFilename(), f.Size)
	fmt.Printf("%d\t%s\t%s\n", f.FSID, f.GetFilename(), gotool.FormatSize(int64(f.Size)))
}

func printFileInfoList(files []*FileInfoDto) {
	for _, f := range files {
		f.PrintOneLine()
	}
	fmt.Printf("Total: %d\n", len(files))
}
