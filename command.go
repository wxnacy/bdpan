package main

// import (
// "fmt"
// "website/site"

// "github.com/akamensky/argparse"
// )

// type Command struct {
// site.Command
// // InitCommand          *argparse.Command
// ListCommand *argparse.Command
// TestCommand *argparse.Command

// ID         *string
// QueryArgId *string
// }

// type ListArg struct {
// Dir *string
// }

// var (
// listArg = ListArg{}
// )

// func NewCommand(parser *argparse.Parser) *Command {
// c := &Command{Command: *site.NewCommand(parser, "bdpan")}
// c.ListCommand = c.NewCommand("list", "列出文件")

// listArg.Dir = c.ListCommand.String("", "dir",
// &argparse.Options{Required: true, Help: "文件名"},
// )

// return c
// }

// func (c *Command) Run() {

// if c.ListCommand.Happened() {
// dto, err := GetDirAllFiles(*listArg.Dir)
// if err != nil {
// panic(err)
// }
// for _, file := range dto {
// fmt.Printf("%s\t%d\n", file.Path, file.FSID)
// }
// } else if c.TestCommand.Happened() {
// }
// }
