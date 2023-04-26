package bdpan

import (
	"bdpan/common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/wxnacy/go-tasker"
	"github.com/wxnacy/gotool"
)

const (
	// FRAGMENT_MAX_SIZE int64 = 4 * 1024 * 1024
	FRAGMENT_MAX_SIZE int64 = 4 * (1 << 20) // 4 M
)

func UploadFile(fromPath, toPath string) (*FileInfoDto, error) {
	res, err := NewUploadFileRequest(fromPath, toPath).Execute()
	if err != nil {
		return nil, err
	}
	return &res.FileInfoDto, err
}

func uploadFile(req UploadFileRequest) (*UploadFileResponse, error) {
	Log.Debugf("uploadFile %#v", req)
	fromPath := req.fromPath
	// 检查文件是否有内容
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() == 0 {
		return nil, errors.New("不能上传空文件")
	}
	// 获取 access_token
	token, err := GetConfigAccessToken()
	if err != nil {
		return nil, err
	}
	// 获取参数
	toPath := req.toPath
	rtype := req.rtype
	isDir := req.isDir
	size := fileInfo.Size()
	localCTime, localMTime := common.GetFileInfoTimes(fileInfo)
	// 创建临时目录
	rand.Seed(time.Now().UnixNano())
	tmpdir := filepath.Join(cacheDir, fmt.Sprintf("upload_tmp_%d", rand.Intn(1000)))
	gotool.DirExistsOrCreate(tmpdir)
	// 分割文件
	paths, err := SplitFile(fromPath, tmpdir, FRAGMENT_MAX_SIZE)
	Log.Debugf("SplitFile paths: %v", paths)
	Log.Debugf("SplitFile error: %v", err)
	if err != nil {
		return nil, err
	}
	blocklist := make([]string, 0)
	for _, _path := range paths {

		md5Str, _ := common.Md5File(_path)
		blocklist = append(blocklist, md5Str)
	}
	blocklistBytes, _ := json.Marshal(blocklist)
	blocklistStr := string(blocklistBytes)
	// 预上传 https://pan.baidu.com/union/doc/3ksg0s9r7
	resp, r, err := GetClient().FileuploadApi.
		Xpanfileprecreate(context.Background()).
		AccessToken(token.AccessToken).Path(toPath).Isdir(isDir).
		Size(int32(size)).Autoinit(1).BlockList(blocklistStr).Rtype(rtype).
		Execute()
	Log.Debugf("Xpanfileprecreate error: %v", err)
	Log.Debugf("Xpanfileprecreate resp: %v", r)
	if err != nil {
		return nil, NewRespError(r)
	}
	uploadId := resp.GetUploadid()
	// 上传分片
	ufTasker := &UploadFragmentTasker{
		Tasker:  tasker.NewTasker(),
		Request: &req, Token: token, UploadId: uploadId, TmpDir: tmpdir,
		Fragments: paths, To: toPath,
	}
	if req.isSync {
		ufTasker.Tasker.Config.UseProgressBar = false
	}
	ufTasker.IsSync = req.isSync
	err = ufTasker.Exec()
	if err != nil {
		return nil, err
	}
	// 创建文件 https://pan.baidu.com/union/doc/rksg0sa17
	_, r, err = GetClient().FileuploadApi.
		Xpanfilecreate(context.Background()).AccessToken(token.AccessToken).
		Path(toPath).Isdir(isDir).Size(int32(size)).Uploadid(uploadId).
		BlockList(blocklistStr).Rtype(rtype).LocalCTime(localCTime.Unix()).
		LocalMTime(localMTime.Unix()).Execute()
	Log.Debugf("Xpanfilecreate resp: %v", r)
	Log.Debugf("Xpanfilecreate error: %v", err)
	if err != nil {
		return nil, NewRespError(r)
	}

	dto := &UploadFileResponse{}
	httpResponseToInterface(r, dto)
	if dto.IsError() {
		return nil, dto.Err()
	}
	return dto, nil
}
