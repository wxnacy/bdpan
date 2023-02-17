package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

type Response struct {
	r             *http.Response
	JsonData      map[string]interface{}
	Err           *ErrorResponse
	ResponseModel interface{}
}

func NewResponse(resp interface{}, r *http.Response, err error) (*Response, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &m); err != nil {
		return nil, err
	}
	e := &ErrorResponse{}
	if r.StatusCode != http.StatusOK {
		if err := json.Unmarshal(bodyBytes, e); err != nil {
			return nil, err
		}

	}
	e.r = r
	return &Response{
		r: r, JsonData: m, Err: e, ResponseModel: resp,
	}, nil
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

func (f FileInfoDto) GetLink() string {
	return f.Dlink + "&access_token=" + _token.AccessToken
}

func (f FileInfoDto) GetFilename() string {
	if f.ServerFilename != "" {
		return f.ServerFilename
	}
	return f.Filename
}

func (f FileInfoDto) PrintOneLine() {
	fmt.Printf("%d\t%s\t%s\t%d\n", f.FSID, f.MD5, f.GetFilename(), f.Size)
}

func printFileInfoList(files []*FileInfoDto) {
	for _, f := range files {
		f.PrintOneLine()
	}
	fmt.Printf("Total: %d", len(files))
}
