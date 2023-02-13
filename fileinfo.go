package main

import (
	"context"
	"errors"
	"fmt"
	"io"
)

func GetDirAllFiles(dir string) ([]*FileInfoDto, error) {
	req := NewFileListRequest().Dir(dir)
	totalList := []*FileInfoDto{}
	fileList := []*FileInfoDto{}
	page := 1
	for {
		res, err := req.Page(page).Execute()
		if err != nil {
			return nil, err
		}
		fileList = res.List
		totalList = append(totalList, fileList...)

		if len(fileList) <= 0 {
			break
		}
		page++
	}
	return totalList, nil
}

func GetFilesByFSIDS(fsids []uint64) ([]*FileInfoDto, error) {
	res, err := NewFileInfoRequest(fsids).Execute()
	if err != nil {
		return nil, err
	}
	return res.List, nil
}

func GetFileByFSID(fsid uint64) (*FileInfoDto, error) {
	res, err := NewFileInfoRequest([]uint64{fsid}).Execute()
	if err != nil {
		return nil, err
	}
	if len(res.List) > 0 {
		return res.List[0], nil
	}
	return nil, errors.New(fmt.Sprintf("fsid %d 找不到", fsid))
}

func fileInfo(req FileInfoRequest) (*FileListResponse, error) {
	_, r, err := apiClient.MultimediafileApi.Xpanmultimediafilemetas(
		context.Background()).AccessToken(_token.AccessToken).Dlink(
		req.GetDlink()).Fsids(req.GetFSID()).Execute()
	if err != nil {
		return nil, err
	}
	return NewFileListResponse(r)
}

func fileList(req FileListRequest) (*FileListResponse, error) {
	dir := *req.dir
	_, r, err := apiClient.FileinfoApi.Xpanfilelist(
		context.Background()).AccessToken(
		_token.AccessToken).Dir(dir).Web(req.GetWeb()).Start(
		req.GetStart()).Limit(req.limit).Execute()
	if err != nil {
		return nil, err
	}

	return NewFileListResponse(r)
}

func SearchFiles(dir, key string) ([]*FileInfoDto, error) {
	res, err := NewFileSearchRequest(key).Execute()
	if err != nil {
		return nil, err
	}
	return res.List, nil
}

func SearchFirstFile(dir, key string) (*FileInfoDto, error) {
	res, err := NewFileSearchRequest(key).Dir(dir).Execute()
	if err != nil {
		return nil, err
	}
	if len(res.List) > 0 {
		return res.List[0], nil
	}
	return nil, errors.New(fmt.Sprintf("%s 找不到", key))
}

func fileSearch(req FileSearchRequest) (*FileListResponse, error) {
	_, r, err := apiClient.FileinfoApi.Xpanfilesearch(
		context.Background()).AccessToken(_token.AccessToken).Dir(req.dir).Key(
		req.key).Execute()
	if err != nil {
		return nil, err
	}

	return NewFileListResponse(r)
}

func GetFileBytes(fsid uint64) ([]byte, error) {
	fileDto, err := GetFileByFSID(fsid)
	if err != nil {
		return nil, err
	}
	uri := fileDto.GetLink()
	headers := map[string]string{
		"User-Agent": "pan.baidu.com",
	}
	var postBody io.Reader
	body, statusCode, err := Do2HTTPRequestToBytes(uri, postBody, headers)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, errors.New("get bytes failed")
	}
	return body, nil
}

func DeleteFile(path string) error {
	res, err := NewFileDeleteRequest([]string{path}).Execute()
	if err != nil {
		return err
	}
	for _, info := range res.Info {
		if info.Path == path && info.Errno > 0 {
			return errors.New(fmt.Sprintf("%s delete failed", path))
		}
	}
	return nil
}

func fileDelete(req FileDeleteRequest) (*FileManagerResponse, error) {
	r, err := apiClient.FilemanagerApi.Filemanagerdelete(
		context.Background()).AccessToken(_token.AccessToken).Async(
		req.Async).Ondup(req.Ondup).Filelist(req.GetFilelistString()).Execute()
	if err != nil {
		return nil, err
	}

	res := &FileManagerResponse{}
	err = httpResponseToInterface(r, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
