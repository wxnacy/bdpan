package bdpan

import (
	"bdpan/common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
		return nil, NewErrorResponse(r).Err()
	}
	uploadId := resp.GetUploadid()
	for i, _path := range paths {

		file, _ := os.Open(_path)
		defer file.Close()
		// 分片上传 https://pan.baidu.com/union/doc/nksg0s9vi
		_, r, err := GetClient().FileuploadApi.Pcssuperfile2(
			context.Background()).AccessToken(token.AccessToken).Partseq(
			strconv.Itoa(i)).Path(
			toPath).Uploadid(uploadId).Type_(req._type).File(file).Execute()
		if strings.HasPrefix(_path, tmpdir) {
			os.Remove(_path)
		}
		Log.Debugf("Pcssuperfile2 path: %s", _path)
		Log.Debugf("Pcssuperfile2 resp: %v", r)
		Log.Debugf("Pcssuperfile2 error: %v", err)
		if err != nil {
			return nil, NewErrorResponse(r).Err()
		}
	}
	// 清理缓存目录
	os.RemoveAll(tmpdir)
	// 创建文件 https://pan.baidu.com/union/doc/rksg0sa17
	createRes, r, err := GetClient().FileuploadApi.
		Xpanfilecreate(context.Background()).AccessToken(token.AccessToken).
		Path(toPath).Isdir(isDir).Size(int32(size)).Uploadid(uploadId).
		BlockList(blocklistStr).Rtype(rtype).LocalCTime(localCTime.Unix()).
		LocalMTime(localMTime.Unix()).Execute()
	Log.Debugf("Xpanfilecreate resp: %v", r)
	Log.Debugf("Xpanfilecreate error: %v", err)
	if err != nil {
		return nil, NewErrorResponse(r).Err()
	}
	if *createRes.Errno > 0 {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Xpanfilecreate error: %s\n", err.Error())
			fmt.Fprintf(os.Stderr, "Xpanfilecreate http response: %v\n", r)
			return nil, err
		}
		return nil, errors.New(string(bytes))
	}

	dto := &UploadFileResponse{}
	httpResponseToInterface(r, dto)
	if dto.IsError() {
		return nil, dto.Err()
	}
	return dto, nil
}

// func UploadDir(fromPath, toPath string) (*UploadDirResponse, error) {
// res, err := NewUploadDirRequest(fromPath, toPath).Execute()
// if err != nil {
// return nil, err
// }
// return res, err
// }

// func uploadDir(req UploadDirRequest) (*UploadDirResponse, error) {
// if req.fromPath == nil {
// return nil, errors.New("fromPath 不能为空")
// }
// if req.toPath == nil {
// return nil, errors.New("toPath 不能为空")
// }
// fromDir := *req.fromPath
// fileInfo, err := os.Stat(fromDir)
// if err != nil {
// return nil, err
// }
// if !fileInfo.IsDir() {
// return nil, errors.New(fmt.Sprintf("%s 不是目录", fromDir))
// }
// fromBaseName := filepath.Base(fromDir)
// toDir := filepath.Join(*req.toPath, fromBaseName)
// infos, err := ioutil.ReadDir(fromDir)
// if err != nil {
// return nil, err
// }
// res := &UploadDirResponse{FailedCount: 0, SuccessCount: 0}
// total := len(infos)
// for _, info := range infos {
// if info.IsDir() {
// continue
// }
// if strings.HasPrefix(info.Name(), ".") {
// continue
// }
// fromPath := filepath.Join(fromDir, info.Name())
// toPath := filepath.Join(toDir, info.Name())
// _, uploadErr := NewUploadFileRequest(fromPath, toPath).Execute()
// if uploadErr != nil {
// res.FailedCount++
// } else {
// res.SuccessCount++
// }
// progress := fmt.Sprintf("Process %d / %d (%d)\n", res.SuccessCount, total, res.FailedCount)
// fmt.Printf("%s \033[K\n", progress) // 输出一行结果
// fmt.Printf("\033[%dA\033[K", 1)     // 将光标向上移动一行
// // fmt.Printf("%s \033[K\n", progress) // 输出第二行结果

// }

// return res, nil
// }
