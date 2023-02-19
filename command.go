package main

import (
	"bdpan/common"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/wxnacy/gotool/arrays"
)

type ICommand interface {
	Happened() bool
	Init() error
	Run() error
}

func NewCommand(c *argparse.Command) *Command {
	return &Command{
		Command: c,
		AppId: c.String("", "app-id",
			&argparse.Options{Required: false, Help: "指定添加 App Id"},
		),
	}

}

type Command struct {
	Command *argparse.Command
	AppId   *string
}

func (c Command) Happened() bool {
	return c.Command.Happened()
}

func (c Command) Init() error {
	return nil
}

func NewCommands(parser *argparse.Parser) []ICommand {
	res := make([]ICommand, 0)
	res = append(res, NewLoginCommand(parser))
	res = append(res, NewQueryCommand(parser))
	res = append(res, NewDeleteCommand(parser))
	res = append(res, NewUploadCommand(parser))
	res = append(res, NewDownloadCommand(parser))
	res = append(res, NewListCommand(parser))
	res = append(res, NewTestCommand(parser))
	return res
}

// *********************************
// QueryCommand
// *********************************

func NewQueryCommand(parser *argparse.Parser) *QueryCommand {
	c := parser.NewCommand("query", "查询网盘")
	cmd := &QueryCommand{
		Command: NewCommand(c),
	}
	cmd.Dir = c.String("", "dir",
		&argparse.Options{Required: false, Help: "查询目录"},
	)
	cmd.Name = c.String("n", "name",
		&argparse.Options{Required: false, Help: "查询名称"},
	)
	cmd.FSIDS = c.StringList("", "fsid",
		&argparse.Options{Required: false, Help: "查询名称"},
	)
	return cmd
}

type QueryCommand struct {
	*Command

	Dir   *string
	Name  *string
	FSIDS *[]string
}

func (q QueryCommand) Run() error {
	dir := *q.Dir
	name := *q.Name

	if len(*q.FSIDS) > 0 {
		fmt.Println("query fsid")
		fsids := make([]uint64, 0)
		for _, fsid := range *q.FSIDS {
			id, err := strconv.Atoi(fsid)
			if err != nil {
				panic(err)
			}
			fsids = append(fsids, uint64(id))
		}

		files, err := GetFilesByFSIDS(fsids)
		if err != nil {
			panic(err)
		}
		printFileInfoList(files)
		return nil
	}

	if name != "" {
		fmt.Println("query name")
		res, err := NewFileSearchRequest(name).Dir(dir).Execute()
		if err != nil {
			panic(err)
		}
		printFileInfoList(res.List)
		return nil
	}
	if dir != "" {
		files, err := GetDirAllFiles(dir)
		if err != nil {
			panic(err)
		}
		printFileInfoList(files)
		return nil
	}
	return nil
}

// *********************************
// ListCommand
// *********************************
func NewListCommand(parser *argparse.Parser) *ListCommand {
	c := parser.NewCommand("list", "列出文件")
	cmd := &ListCommand{
		Command: NewCommand(c),
	}
	cmd.Dir = c.String("", "dir",
		&argparse.Options{Required: true, Help: "查询目录"},
	)
	cmd.IsRecursion = c.Flag("r", "recursion",
		&argparse.Options{Required: false, Help: "是否遍历子目录，默认否"},
	)
	cmd.IsDesc = c.Flag("d", "desc",
		&argparse.Options{Required: false, Help: "是否为导出，默认否"},
	)
	cmd.Limit = c.Int("l", "limit",
		&argparse.Options{Required: false, Help: "单页文件个数，默认 10", Default: 10},
	)
	cmd.Start = c.Int("s", "start",
		&argparse.Options{Required: false, Help: "查询起点，默认为0"},
	)
	cmd.Order = c.String("o", "order",
		&argparse.Options{
			Required: false, Help: "排序字段:time(修改时间)，name(文件名)，size(大小，目录无大小)，默认为文件类型",
			Validate: func(args []string) error {
				for _, arg := range args {
					if !arrays.StringContains(ORDERS, arg) {
						return errors.New(fmt.Sprintf("%s not in %v", arg, ORDERS))
					}
				}
				return nil
			},
		},
	)
	return cmd
}

const (
	ORDER_TIME = "time"
	ORDER_NAME = "name"
	ORDER_SIZE = "size"
)

var (
	ORDERS = []string{ORDER_NAME, ORDER_TIME, ORDER_SIZE}
)

type ListCommand struct {
	*Command

	Dir         *string
	IsRecursion *bool
	IsDesc      *bool
	Start       *int
	Limit       *int
	Order       *string
}

func (l ListCommand) Run() error {
	dir := *l.Dir
	var recursion int32
	if *l.IsRecursion {
		recursion = 1
	}
	var desc int32
	if *l.IsDesc {
		desc = 1
	}
	res, err := NewFileListAllRequest(dir).Recursion(recursion).Limit(
		int32(*l.Limit)).Order(*l.Order).Desc(desc).Start(int32(*l.Start)).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "查询失败 %s", err.Error())
		return err
	}
	printFileInfoList(res.List)
	if res.HasMore == 1 {
		fmt.Printf("下一页查询命令 %d\n", res.Cursor)

	}
	return nil
}

// *********************************
// LoginCommand
// *********************************
func NewLoginCommand(parser *argparse.Parser) *LoginCommand {
	c := parser.NewCommand("login", "登录网盘")
	cmd := &LoginCommand{
		Command: NewCommand(c),
	}
	return cmd
}

type LoginCommand struct {
	*Command
}

func (l LoginCommand) buildCredentail() Credential {
	appId := *l.AppId

	credential := Credential{}
	fmt.Println("请先完善秘钥信息")
	if appId == "" {

		fmt.Print("App Id: ")
		fmt.Scanln(&credential.AppId)
	} else {
		credential.AppId = appId
	}
	fmt.Print("App Key: ")
	fmt.Scanln(&credential.AppKey)
	fmt.Print("Secret Key: ")
	fmt.Scanln(&credential.SecretKey)
	fmt.Print("Sign Key: ")
	fmt.Scanln(&credential.SignKey)
	return credential
}

func (l LoginCommand) getVipName(vipType int32) string {
	switch vipType {
	case 0:
		return "普通用户"
	case 1:
		return "普通会员"
	case 2:
		return "超级会员"
	}
	return "未知身份"
}

func (l LoginCommand) Run() error {
	appId := *l.AppId
	// var cres []*Credential
	if appId == "" {
		_, err = GetCredentails()
	} else {
		_, err = GetCredentail(appId)
	}
	// 当查询不到用户时进行添加流程
	if err != nil {
		credential := l.buildCredentail()
		err := AddCredentail(credential)
		if err != nil {
			fmt.Fprintf(os.Stderr, "登录失败 %s\n", err.Error())
			return err
		}
		// 获取授权
		err = CreateAccessTokenByDeviceCode()
		if err != nil {
			fmt.Fprintf(os.Stderr, "登录失败 %s\n", err.Error())
			return err
		}
	}

	if appId != "" {

		// 设置当前需要使用的 appId
		config, err := GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "登录失败 %s\n", err.Error())
			return err
		}
		config.LoginAppId = appId
	}

	user, err := userInfo()
	// 获取用户信息失败，可能是授权过期则需要进行授权操作
	if err != nil {
		err = CreateAccessTokenByDeviceCode()
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取用户信息失败 %s\n", err.Error())
			return err
		}
	}
	fmt.Printf("Hello, %s(%s)\n", user.GetNetdiskName(), l.getVipName(user.GetVipType()))
	pan, err := panInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取网盘信息失败 %s\n", err.Error())
		return err
	}
	fmt.Printf("网盘容量 %s/%s\n", formatSize(pan.GetUsed()), formatSize(pan.GetTotal()))
	// fmt.Printf("网盘总容量 %d", pan.GetTotal())
	return nil
}

// *********************************
// DeleteCommand
// *********************************
func NewDeleteCommand(parser *argparse.Parser) *DeleteCommand {
	c := parser.NewCommand("del", "删除文件")
	cmd := &DeleteCommand{
		Command: NewCommand(c),
	}

	cmd.Path = c.String("p", "path",
		&argparse.Options{Required: false, Help: "文件地址"},
	)
	return cmd
}

type DeleteCommand struct {
	*Command

	Path *string
}

func (d DeleteCommand) Run() error {
	path := *d.Path
	if path != "" {
		err = DeleteFile(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// *********************************
// UploadCommand
// *********************************
func NewUploadCommand(parser *argparse.Parser) *UploadCommand {
	c := parser.NewCommand("upload", "上传文件")
	cmd := &UploadCommand{
		Command: NewCommand(c),
	}

	cmd.From = c.String("f", "from",
		&argparse.Options{Required: true, Help: "文件来源"},
	)
	cmd.To = c.String("t", "to",
		&argparse.Options{
			Required: false, Help: "保存地址", Default: DEFAULT_UPLOAD_DIR},
	)
	cmd.IsSync = c.Flag("", "sync",
		&argparse.Options{
			Required: false, Help: "是否同步上传"},
	)
	return cmd
}

type UploadCommand struct {
	*Command

	From   *string
	To     *string
	IsSync *bool
}

func (u UploadCommand) Run() error {
	from := *u.From
	to := *u.To
	if common.FileExists(from) {
		if strings.HasSuffix(to, "/") {
			to = filepath.Join(to, filepath.Base(from))
		}
		fmt.Printf("Upload %s to %s\n", from, to)
		_, err = UploadFile(from, to)
		if err != nil {
			return err
		}
		fmt.Printf("File: %s upload success\n", from)
	}

	if common.DirExists(from) {
		if *u.IsSync {

			res, err := UploadDir(from, to)
			if err != nil {
				return err
			}
			fmt.Printf("Success: %d\n", res.SuccessCount)
			fmt.Printf("Failed: %d\n", res.FailedCount)
		} else {
			TaskUploadDir(from, to)
		}
	}
	return nil
}

// *********************************
// TestCommand
// *********************************
func NewTestCommand(parser *argparse.Parser) *TestCommand {
	c := parser.NewCommand("test", "测试程序")
	cmd := &TestCommand{
		Command: NewCommand(c),
	}
	return cmd
}

type TestCommand struct {
	*Command
}

func (t TestCommand) Run() error {
	fmt.Println(float64(17599702237186) / (1 << 30))
	return nil
}

// *********************************
// DownloadCommand
// *********************************
func NewDownloadCommand(parser *argparse.Parser) *DownloadCommand {
	c := parser.NewCommand("download", "下载文件")
	cmd := &DownloadCommand{
		Command: NewCommand(c),
	}

	cmd.From = c.String("f", "from",
		&argparse.Options{Required: true, Help: "网盘文件"},
	)
	downloadDir, _ := common.ExpandUser("~/Downloads")
	cmd.To = c.String("t", "to",
		&argparse.Options{
			Required: false, Help: "下载位置", Default: downloadDir},
	)
	cmd.IsSync = c.Flag("", "sync",
		&argparse.Options{
			Required: false, Help: "是否同步进程"},
	)
	return cmd
}

type DownloadCommand struct {
	*Command

	From   *string
	To     *string
	IsSync *bool
}

func (d DownloadCommand) Run() error {
	from := *d.From
	to := *d.To
	path := filepath.Join(to, filepath.Base(from))
	if common.FileExists(path) {
		fmt.Fprintf(os.Stderr, "下载失败 %s 已存在\n", path)
		return err
	}
	file, err := GetFileByPath(from)
	if err != nil {
		fmt.Fprintf(os.Stderr, "下载失败 %s\n", err.Error())
		return err
	}
	// TODO: 判定 to 的类型
	bytes, err := GetFileBytes(file.FSID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "下载失败 %s\n", err.Error())
		return err
	}

	err = os.WriteFile(path, bytes, common.PermFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "下载失败 %s\n", err.Error())
		return err
	}
	fmt.Printf("%s 下载成功\n", path)

	return nil
}
