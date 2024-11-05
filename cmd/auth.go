/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/wxnacy/bdpan"

	"github.com/spf13/cobra"
)

var authCommond = &AuthCommand{}

type AuthCommand struct {
}

func (a *AuthCommand) handleAction(item *bdpan.SelectItem) error {
	switch item.Action {
	case bdpan.ActionSystem:
		return a.promptAuths()
	}
	return nil
}

func (a AuthCommand) promptAuths() error {
	credentails, err := bdpan.GetCredentails()
	if err != nil {
		return err
	}
	items := make([]*bdpan.SelectItem, 0)
	handle := func(item *bdpan.SelectItem) error {
		return a.handleAction(item)
	}
	for _, c := range credentails {
		item := &bdpan.SelectItem{
			Name:   c.AppId,
			Info:   c,
			Desc:   c.AppKey,
			Action: bdpan.ActionSystem,
		}
		items = append(items, item)
	}
	return bdpan.PromptSelect("全部账户", items, true, 10, func(index int, s string) error {
		return handle(items[index])
	})
}

func (a AuthCommand) Run() error {
	return a.promptAuths()
}

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "账户操作",
	Run: func(cmd *cobra.Command, args []string) {
		handleCmdErr(authCommond.Run())
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
