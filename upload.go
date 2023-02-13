package main

import (
	"bdpan/common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	FRAGMENT_MAX_SIZE int32 = 4 * 1024 * 1024
)

func UploadFile(fromPath, toPath string) (*FileInfoDto, error) {
	res, err := NewUploadFileRequest(fromPath, toPath).Execute()
	if err != nil {
		return nil, err
	}
	return &res.FileInfoDto, err
}

func uploadFile(req UploadFileRequest) (*UploadFileResponse, error) {
	fromPath := req.fromPath
	toPath := req.toPath
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return nil, err
	}
	rtype := req.rtype
	isDir := req.isDir
	size := int32(fileInfo.Size())
	if size > FRAGMENT_MAX_SIZE {
		return nil, errors.New("不能上传大于 4M 的文件")
	}
	md5Str, _ := common.Md5File(fromPath)
	blocklist := []string{md5Str}
	blocklistBytes, _ := json.Marshal(blocklist)
	blocklistStr := string(blocklistBytes)
	resp, r, err := apiClient.FileuploadApi.Xpanfileprecreate(
		context.Background()).AccessToken(_token.AccessToken).Path(
		toPath).Isdir(isDir).Size(size).Autoinit(
		1).BlockList(blocklistStr).Rtype(rtype).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Xpanfileprecreate error: %s\n", err.Error())
		fmt.Fprintf(os.Stderr, "Xpanfileprecreate http response: %v\n", r)
		return nil, err
	}
	uploadId := resp.GetUploadid()
	file, _ := os.Open(fromPath)
	_, r, err = apiClient.FileuploadApi.Pcssuperfile2(
		context.Background()).AccessToken(_token.AccessToken).Partseq(
		req.GetPartseq()).Path(
		toPath).Uploadid(uploadId).Type_(req._type).File(file).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Pcssuperfile2 error: %s\n", err.Error())
		fmt.Fprintf(os.Stderr, "Pcssuperfile2 http response: %v\n", r)
		return nil, err
	}
	_, r, err = apiClient.FileuploadApi.Xpanfilecreate(
		context.Background()).AccessToken(_token.AccessToken).Path(toPath).Isdir(
		isDir).Size(size).Uploadid(uploadId).BlockList(blocklistStr).Rtype(rtype).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Xpanfilecreate error: %s\n", err.Error())
		fmt.Fprintf(os.Stderr, "Xpanfilecreate http response: %v\n", r)
		return nil, err
	}

	dto := &UploadFileResponse{}
	httpResponseToInterface(r, dto)
	return dto, nil
}

func UploadDir(req UploadDirRequest) (*UploadDirResponse, error) {
	if req.fromPath == nil {
		return nil, errors.New("fromPath 不能为空")
	}
	if req.toPath == nil {
		return nil, errors.New("toPath 不能为空")
	}
	fromDir := *req.fromPath
	fileInfo, err := os.Stat(fromDir)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, errors.New(fmt.Sprintf("%s 不是目录", fromDir))
	}
	fromBaseName := filepath.Base(fromDir)
	toDir := filepath.Join(*req.toPath, fromBaseName)
	infos, err := ioutil.ReadDir(fromDir)
	if err != nil {
		return nil, err
	}
	res := &UploadDirResponse{FailedCount: 0, SuccessCount: 0}
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		if strings.HasPrefix(info.Name(), ".") {
			continue
		}
		fromPath := filepath.Join(fromDir, info.Name())
		toPath := filepath.Join(toDir, info.Name())
		_, uploadErr := NewUploadFileRequest(fromPath, toPath).Execute()
		if uploadErr != nil {
			res.FailedCount++
		} else {
			res.SuccessCount++
		}
	}

	return res, nil
}
