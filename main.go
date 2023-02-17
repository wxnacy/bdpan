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
	fmt.Println("main")
	parser := argparse.NewParser("bdpan", "网站服务")
	commands := NewCommands(parser)
	err = parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	for _, cmd := range commands {
		if cmd.Happened() {
			err = cmd.Run()
			if err != nil {
				panic(err)
			}
			return
		}
	}
	fmt.Println("没有此命令")
}
