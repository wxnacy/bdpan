package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	accessToken string
)

func init() {
	accessToken = os.Getenv("BDPAN_ACCESS_TOKEN")
}

func TestGetFileList(t *testing.T) {
	req := NewGetFileListReq()
	res, err := GetFileList(accessToken, req)
	assert.NoError(t, err)
	fileList := res.List
	assert.GreaterOrEqual(t, len(fileList), 1)
}

func TestGetFileInfo(t *testing.T) {
	listRes, err := GetFileList(accessToken, NewGetFileListReq())
	reqInfo := listRes.List[0]
	fsid := reqInfo.FSID

	req := NewGetFileInfoReq(fsid)
	res, err := GetFileInfo(accessToken, req)
	assert.NoError(t, err)
	assert.Equal(t, res.Path, reqInfo.Path)
}

func TestBatchGetFileInfo(t *testing.T) {
	listRes, err := GetFileList(accessToken, NewGetFileListReq())
	reqInfo := listRes.List[0]
	fsid := reqInfo.FSID

	req := NewBatchGetFileListReq(fsid)
	for _, f := range listRes.List {
		req.AppendFSID(f.FSID)
	}
	res, err := BatchGetFileInfo(accessToken, req)
	assert.NoError(t, err)
	assert.Equal(t, len(res.List), len(listRes.List))
}
