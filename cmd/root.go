/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bdpan"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	argAppId string
	argDebug bool

	Log = bdpan.Log
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bdpan",
	Short: "百度网盘命令行工具",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&argAppId, "app-id", "", "指定添加 App Id")
	rootCmd.PersistentFlags().BoolVarP(&argDebug, "debug", "D", false, "debug 模式")

	// 运行前全局命令
	cobra.OnInitialize(func() {
		if argDebug {
			bdpan.SetLogLevel(logrus.DebugLevel)
		}
	})
}
