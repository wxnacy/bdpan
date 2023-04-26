package bdpan

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type UploadFileRequest struct {
	// UploadRequest
	fromPath string
	toPath   string
	_type    string
	partseq  int
	// rtype
	// 文件命名策略，默认0
	// 0 为不重命名，返回冲突
	// 1 为只要path冲突即重命名
	// 2 为path冲突且block_list不同才重命名
	// 3 为覆盖，需要与预上传precreate接口中的rtype保持一致
	rtype  int32
	isDir  int32
	size   int32
	isSync bool
}

func NewUploadFileRequest(fromPath, toPath string) *UploadFileRequest {
	return &UploadFileRequest{
		fromPath: fromPath, toPath: toPath, _type: "tmpfile",
		partseq: 0, rtype: int32(3), isDir: int32(0),
	}
}

func (r UploadFileRequest) FromPath(path string) UploadFileRequest {
	r.fromPath = path
	return r
}

func (r UploadFileRequest) ToPath(path string) UploadFileRequest {
	r.toPath = path
	return r
}

func (r UploadFileRequest) Type(typ string) UploadFileRequest {
	r._type = typ
	return r
}

func (r *UploadFileRequest) RType(t int32) *UploadFileRequest {
	r.rtype = t
	return r
}

func (r *UploadFileRequest) IsSync(s bool) *UploadFileRequest {
	r.isSync = s
	return r
}

func (r UploadFileRequest) GetPartseq() string {
	return strconv.Itoa(r.partseq)
}

func (r UploadFileRequest) Execute() (*UploadFileResponse, error) {
	return uploadFile(r)
}

type FileListRequest struct {
	dir   *string
	web   int
	page  int
	start int
	limit int32
}

func NewFileListRequest() FileListRequest {
	return FileListRequest{
		web: 1, start: 0, limit: 1000,
	}
}

func (r FileListRequest) Dir(dir string) FileListRequest {
	r.dir = &dir
	return r
}

func (r FileListRequest) Web(web int) FileListRequest {
	r.web = web
	return r
}

func (r FileListRequest) GetWeb() string {
	return fmt.Sprintf("%d", r.web)
}

func (r FileListRequest) Page(page int) FileListRequest {
	r.page = page
	r.start = (page - 1) * int(r.limit)
	return r
}

func (r FileListRequest) GetStart() string {
	return fmt.Sprintf("%d", r.start)
}

func (r FileListRequest) Execute() (*FileListResponse, error) {
	return fileList(r)
}

// ****************************************
// FileListAllRequest
// ****************************************

type FileListAllRequest struct {
	path string
	web  int
	// page  int
	start     int32
	limit     int32
	recursion int32
	desc      int32
	order     string
}

func NewFileListAllRequest(path string) FileListAllRequest {
	return FileListAllRequest{
		web: 1, start: 0, limit: 1000, path: path,
	}
}

func (f FileListAllRequest) Path(path string) FileListAllRequest {
	f.path = path
	return f
}

func (r FileListAllRequest) Web(web int) FileListAllRequest {
	r.web = web
	return r
}

func (f FileListAllRequest) Recursion(r int32) FileListAllRequest {
	f.recursion = r
	return f
}

func (f FileListAllRequest) Desc(r int32) FileListAllRequest {
	f.desc = r
	return f
}

func (f FileListAllRequest) Limit(r int32) FileListAllRequest {
	f.limit = r
	return f
}

func (f FileListAllRequest) Start(r int32) FileListAllRequest {
	f.start = r
	return f
}

func (f FileListAllRequest) Order(order string) FileListAllRequest {
	f.order = order
	return f
}

func (r FileListAllRequest) GetWeb() string {
	return fmt.Sprintf("%d", r.web)
}

func (r FileListAllRequest) Execute() (*FileListAllResponse, error) {
	return fileListAll(r)
}

// ****************************************
// FileInfoRequest
// ****************************************

type FileInfoRequest struct {
	dlink int
	fsids []uint64
}

func NewFileInfoRequest(fsids []uint64) FileInfoRequest {
	return FileInfoRequest{
		dlink: 1, fsids: fsids,
	}
}

func (r FileInfoRequest) FSIDs(fsids []uint64) FileInfoRequest {
	r.fsids = fsids
	return r
}

func (r FileInfoRequest) FSID(fsid uint64) FileInfoRequest {
	r.fsids = append(r.fsids, fsid)
	return r
}

func (r FileInfoRequest) GetFSID() string {
	bytesData, _ := json.Marshal(r.fsids)
	return string(bytesData)
}

func (r FileInfoRequest) GetDlink() string {
	return fmt.Sprintf("%d", r.dlink)
}

func (r FileInfoRequest) Execute() (*FileListResponse, error) {
	return fileInfo(r)
}

// ****************************************
// FileSearchRequest
// ****************************************

type FileSearchRequest struct {
	dir       string
	key       string
	recursion int
}

func NewFileSearchRequest(key string) FileSearchRequest {
	return FileSearchRequest{
		key: key,
	}
}

func (r FileSearchRequest) Dir(dir string) FileSearchRequest {
	r.dir = dir
	return r
}

func (f FileSearchRequest) Recursion(r int) FileSearchRequest {
	f.recursion = r
	return f
}

func (f FileSearchRequest) GetRecursion() string {
	return strconv.Itoa(f.recursion)
}

func (r FileSearchRequest) Execute() (*FileListResponse, error) {
	return fileSearch(r)
}

// ****************************************
// FileDeleteRequest
// ****************************************

type FileManagerFilelist struct {
	Path    string `json:"path,omitempty"`
	Newname string `json:"newname,omitempty"`
	Dest    string `json:"dest,omitempty"`
	Ondup   string `json:"ondup,omitempty"`
}

type FileManagerRequest struct {
	Filelist []FileManagerFilelist
	Async    int32
	Ondup    string
}

func newFileManagerRequest(filelist []FileManagerFilelist) FileManagerRequest {
	return FileManagerRequest{
		Async: int32(1), Ondup: "overwrite", Filelist: filelist,
	}
}

func (fm FileManagerRequest) GetFilelistString() string {
	bytes, _ := json.Marshal(fm.Filelist)
	return string(bytes)
}

// ****************************************
// FileDeleteRequest
// ****************************************

type FileDeleteRequest struct {
	FileManagerRequest
}

func NewFileDeleteRequest(paths []string) FileDeleteRequest {
	filelist := []FileManagerFilelist{}
	for _, path := range paths {
		filelist = append(filelist, FileManagerFilelist{Path: path})
	}
	return FileDeleteRequest{
		FileManagerRequest: newFileManagerRequest(filelist),
	}
}

func (r FileDeleteRequest) Execute() (*FileManagerResponse, error) {
	return fileDelete(r)
}
