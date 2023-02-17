package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

func init() {
	fmt.Println("main init")
	initArgparse()
	initAll()
}

var (
	err           error
	testCommand   *argparse.Command
	loginCommand  *argparse.Command
	queryCommand  *argparse.Command
	configCommand *argparse.Command

	loginArg LoginArg
	queryArg QueryArg
)

func initArgparse() {
	parser := argparse.NewParser("bdpan", "网站服务")

	configCommand = parser.NewCommand("config", "修改和获取配置")
	testCommand = parser.NewCommand("test", "测试程序")
	initLoginArgparse(parser)
	initQueryArgparse(parser)
	// // Parse input
	err = parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		// panic(parser.Usage(err))
	}
}

func initLoginArgparse(parser *argparse.Parser) {
	loginCommand = parser.NewCommand("login", "登录网盘")
	loginArg = LoginArg{}
	loginArg.AppId = loginCommand.String("", "app-id",
		&argparse.Options{Required: false, Help: "指定添加 App Id"},
	)

}

type LoginArg struct {
	AppId *string
	// AppKey    *string
	// SecretKey *string
	// SignKey   *string
}

type QueryArg struct {
	AppId *string
	Dir   *string
	Name  *string
	FSIDS *[]string
	// AppKey    *string
	// SecretKey *string
	// SignKey   *string
}

func initQueryArgparse(parser *argparse.Parser) {
	queryCommand = parser.NewCommand("query", "查询网盘")
	queryArg = QueryArg{}
	queryArg.AppId = queryCommand.String("", "app-id",
		&argparse.Options{Required: false, Help: "指定添加 App Id"},
	)
	queryArg.Dir = queryCommand.String("", "dir",
		&argparse.Options{Required: false, Help: "查询目录"},
	)
	queryArg.Name = queryCommand.String("n", "name",
		&argparse.Options{Required: false, Help: "查询名称"},
	)
	queryArg.FSIDS = queryCommand.StringList("", "fsid",
		&argparse.Options{Required: false, Help: "查询名称"},
	)
}

func main() {
	fmt.Println("main")
	if !testCommand.Happened() {
		_, err = defaultCredentail()
		if err != nil {
			fmt.Println("请先执行 bdpan login 进行登录")
			return
		}
	}
	if testCommand.Happened() {
		fmt.Println("test")
	} else if loginCommand.Happened() {
		Login(loginArg)
	} else if queryCommand.Happened() {
		Query(queryArg)
	} else {
		fmt.Println("no command")
	}
}
