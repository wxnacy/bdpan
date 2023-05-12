package bdpan

import (
	"fmt"
	"strings"

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

type SelectAction int

const (
	ActionSystem SelectAction = iota + 100
	ActionDelete
)

type SelectItem struct {
	Name   string
	Icon   string
	Action SelectAction
	Desc   string
	Info   interface{}
}

type SelectFunc func(index int, s string) error

func PromptSelect(
	label string, actions []*SelectItem, hasSearch bool, limit int,
	selectFunc SelectFunc,
) (err error) {
	activeTempleFmt := "%s {{ .Icon }} {{ .Name | %s}}"
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   fmt.Sprintf(activeTempleFmt, promptui.IconSelect, "blue"),
		Inactive: fmt.Sprintf(activeTempleFmt, " ", "cyan"),
		// Selected: "\U0001f449  {{ .GetFileTypeIcon }}  {{ .Path | red | cyan }}",
		Details: `{{ .Desc }}`,
	}
	searcher := func(input string, index int) bool {
		pepper := actions[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:        label,
		Items:        actions,
		Templates:    templates,
		Size:         limit,
		IsVimMode:    true,
		HideSelected: true, // 隐藏选择后输出的内容
	}
	if hasSearch {
		prompt.Searcher = searcher
	}
	index, s, err := prompt.Run()
	if err != nil {
		return
	}
	err = selectFunc(index, s)
	if err != nil {
		return
	}
	return
}
