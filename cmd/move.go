/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var moveCommand = &ManageCommand{opera: bdpan.OperaMove}

type ManageCommand struct {
	path  string
	to    string
	opera bdpan.FileManageOpera
}

func (m ManageCommand) Exec(args []string) error {
	if len(args) > 0 {
		m.path = args[0]
	} else {
		if m.path == "" {
			return fmt.Errorf("缺少 path 参数")
		}
	}
	_, err := bdpan.GetFileByPath(m.path)
	if err != nil {
		return err
	}
	switch m.opera {
	case bdpan.OperaDelete:
		err = bdpan.DeleteFile(m.path)
	case bdpan.OperaCopy:
		err = m.handleOpera(m.path, m.to, bdpan.CopyFile)
	case bdpan.OperaMove:
		err = m.handleOpera(m.path, m.to, bdpan.MoveFile)
	}
	if err != nil {
		return err
	}
	return nil
}

func (m ManageCommand) handleOpera(from, to string, fn func(string, string) error) error {
	err := fn(from, to)
	if err == bdpan.ErrPathExists {
		var input string
		fmt.Printf("%s 已存在，是否重命名(y/N): ", to)
		fmt.Scanln(&input)
		if input == "y" {
			var newPath string
			fmt.Print("请输入新地址: ")
			fmt.Scanln(&newPath)
			if newPath == "" {
				return fmt.Errorf("输入地址有误")
			}
			if !strings.HasPrefix(newPath, "/") {
				newPath = filepath.Join(filepath.Dir(to), newPath)
			}
			return fn(from, newPath)
		} else {
			fmt.Println("操作取消")
			return nil
		}
	} else {
		return err
	}
}

func (m ManageCommand) Copy() error {
	return nil
	// 判定源文件是否存在
	// _, err := bdpan.GetFileByPath(m.path)
	// if err != nil {
	// return err
	// }
	// return m.handleOpera(m.path, m.to, bdpan.CopyFile)
	// toPath := m.to
	// err = bdpan.CopyFile(m.path, toPath)
	// if err == bdpan.ErrPathExists {
	// var input string
	// fmt.Printf("%s 已存在，是否重命名(y/N): ", toPath)
	// fmt.Scanln(&input)
	// if input == "y" {
	// fmt.Print("请输入新地址: ")
	// fmt.Scanln(&input)
	// }
	// }
	// return err
	// if err == bdpan.ErrPathExists {
	// tools.
	// }
	// 目标文件已存在，则报错
	// toFile, err := bdpan.GetFileByPath(m.to)
	// if err != nil && err != bdpan.ErrPathNotFound {
	// return err
	// }
	// if toFile == nil {
	// return bdpan.CopyFile(m.path, m.to)
	// } else {
	// if !toFile.IsDir() {
	// Log.Warn(bdpan.ErrPathExists.Error())
	// return nil
	// }
	// // 判定是否需要重命名
	// toPath := filepath.Join(m.to, filepath.Base(m.path))
	// toFile, err = bdpan.GetFileByPath(toPath)
	// if err != nil && err != bdpan.ErrPathNotFound {
	// return err
	// }
	// if toFile != nil {
	// toPath = tools.FileAutoReDownloadName(toPath)
	// }
	// return bdpan.CopyFile(m.path, toPath)
	// }
}

// moveCmd represents the move command
var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "移动文件",
	Run: func(cmd *cobra.Command, args []string) {
		handleCmdErr(moveCommand.Exec(args))
	},
}

func init() {
	moveCmd.Flags().StringVarP(&moveCommand.path, "path", "p", "", "操作地址")
	moveCmd.Flags().StringVarP(&moveCommand.to, "to", "o", "", "移动目标")
	rootCmd.AddCommand(moveCmd)
}
