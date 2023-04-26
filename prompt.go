package bdpan

import (
	"github.com/manifoldco/promptui"
)

func PromptConfirm(confirm string) bool {
	prompt := promptui.Prompt{
		Label:     confirm,
		IsConfirm: true,
	}
	_, err := prompt.Run()
	if err != nil {
		return false
	}
	return true
}
