package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

var (
	err error
)

func main() {
	// 注册命令工具
	parser := argparse.NewParser("bdpan", "百度网盘命令行工具")
	commands := NewCommands(parser)
	err = parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	for _, cmd := range commands {
		if cmd.Happened() {
			err = cmd.Init()
			if err != nil {
				panic(err)
			}
			err = cmd.Run()
			if err != nil {
				panic(err)
			}
			return
		}
	}
	fmt.Println("没有此命令")
}
