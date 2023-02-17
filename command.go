package main

import (
	"bdpan/common"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/akamensky/argparse"
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
		fmt.Println("query dir")
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

func (l LoginCommand) Run() error {
	appId := *l.AppId
	// var cres []*Credential
	if appId == "" {
		_, err = GetCredentails()
	} else {
		_, err = GetCredentail(appId)

	}
	if err != nil {
		credential := l.buildCredentail()
		err := AddCredentail(credential)
		if err != nil {
			fmt.Fprintf(os.Stderr, "登录失败 %s\n", err.Error())
		}
	}

	config, err := GetConfig()
	if err != nil {
		panic(err)
	}

	var c *Credential
	if appId != "" {
		config.LoginAppId = appId
		c, err = GetCredentail(appId)
	} else {
		c, err = GetConfigCredentail()
	}
	if err != nil {
		panic(err)
	}
	// err = CreateAccessTokenByDeviceCode()
	// if err != nil {
	// panic(err)
	// }
	// kt := &AccessToken{}
	// t.AccessToken = "1"
	// t.RefreshToken = "1"
	// saveAccessToken(c.AppId, *t)
	token, err := c.GetAccessToken()
	if err != nil {
		panic(err)
	}
	fmt.Println(*token)
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
	return cmd
}

type UploadCommand struct {
	*Command

	From *string
	To   *string
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
		res, err := UploadDir(from, to)
		if err != nil {
			return err
		}
		fmt.Printf("Success: %d\n", res.SuccessCount)
		fmt.Printf("Failed: %d\n", res.FailedCount)
	}
	return nil
}
