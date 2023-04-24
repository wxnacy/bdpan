package bdpan

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
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

		if len(fileList) <= 0 || len(fileList) < int(req.limit) {
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
	token, err := GetConfigAccessToken()
	if err != nil {
		return nil, err
	}
	_, r, err := GetClient().MultimediafileApi.Xpanmultimediafilemetas(
		context.Background()).AccessToken(token.AccessToken).Dlink(
		req.GetDlink()).Fsids(req.GetFSID()).Execute()
	if err != nil {
		return nil, err
	}
	return NewFileListResponse(r)
}

// https://pan.baidu.com/union/doc/nksg0sat9
func fileList(req FileListRequest) (*FileListResponse, error) {
	dir := *req.dir
	token, err := GetConfigAccessToken()
	if err != nil {
		return nil, err
	}
	_, r, err := GetClient().FileinfoApi.Xpanfilelist(
		context.Background()).AccessToken(
		token.AccessToken).Dir(dir).Web(req.GetWeb()).Start(
		req.GetStart()).Limit(req.limit).Execute()
	Log.Debugf("Xpanfilelist resp: %v", r)
	Log.Debugf("Xpanfilelist err: %v", err)
	if err != nil {
		return nil, err
	}

	return NewFileListResponse(r)
}

// func Search(dir, key string) ([]*FileInfoDto, error) {
// res, err := NewFileSearchRequest(key).Execute()
// if err != nil {
// return nil, err
// }
// return res.List, nil
// res, err := NewFileSearchRequest(name).Dir(dir).Execute()
// Log.Debugf("search resp: %v", res)
// Log.Debugf("search error: %v", err)
// if err != nil {
// return nil, err
// }
// if res.IsError() {
// return nil, res.Err()
// }
// return res, nil
// }

// TODO: 暂时用遍历的方式查找文件，后期需要改为搜索
func SearchFiles(dir, key string) ([]*FileInfoDto, error) {
	files, err := GetDirAllFiles(dir)
	if err != nil {
		return nil, err
	}
	res := make([]*FileInfoDto, 0)
	for _, f := range files {
		if strings.Contains(f.GetFilename(), key) {
			res = append(res, f)
		}
	}
	return res, nil
}

func GetFileByPath(path string) (*FileInfoDto, error) {
	path = strings.TrimRight(path, "/")
	Log.Infof("开始查询文件: %s", path)
	dir, name := filepath.Split(path)
	files, err := SearchFiles(dir, name)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.GetFilename() == name {
			Log.Infof("查询到文件类型为: %s", f.GetFileType())
			return f, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("%s 找不到", path))
}

// https://pan.baidu.com/union/doc/zksg0sb9z
func fileSearch(req FileSearchRequest) (*FileListResponse, error) {
	token, err := GetConfigAccessToken()
	if err != nil {
		return nil, err
	}
	_, r, err := GetClient().FileinfoApi.Xpanfilesearch(
		context.Background()).AccessToken(token.AccessToken).Key(
		req.key).Recursion(req.GetRecursion()).Execute()
	Log.Debugf("Xpanfilesearch resp: %v", r)
	Log.Debugf("Xpanfilesearch error: %v", err)
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
	token, err := GetConfigAccessToken()
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s&access_token=%s", fileDto.Dlink, token.AccessToken)
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
	token, err := GetConfigAccessToken()
	if err != nil {
		return nil, err
	}
	r, err := GetClient().FilemanagerApi.Filemanagerdelete(
		context.Background()).AccessToken(token.AccessToken).Async(
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

// https://pan.baidu.com/union/doc/Zksg0sb73
func fileListAll(req FileListAllRequest) (*FileListAllResponse, error) {
	token, err := GetConfigAccessToken()
	if err != nil {
		return nil, err
	}
	_, r, err := GetClient().MultimediafileApi.Xpanfilelistall(
		context.Background()).AccessToken(
		token.AccessToken).Path(req.path).Web(req.GetWeb()).Start(
		req.start).Limit(req.limit).Order(req.order).Recursion(
		req.recursion).Desc(req.desc).Execute()
	Log.Debugf("Xpanfilelistall resp: %v", r)
	Log.Debugf("Xpanfilelistall error: %v", err)
	if err != nil {
		return nil, NewErrorResponse(r).Err()
	}

	res := &FileListAllResponse{}
	err = httpResponseToInterface(r, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
