/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	testCommand *TestCommand
)

type TestCommand struct {
	Cmd *cobra.Command
	ID  string
}

func (t TestCommand) Init() {
	t.Cmd.Flags().StringVar(&t.ID, "id", "", "")
}

func (t TestCommand) Run(cmd *cobra.Command, args []string) {
	fmt.Println(t.ID)
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	testCommand = &TestCommand{Cmd: testCmd}
	testCommand.Init()
	testCmd.Run = testCommand.Run
	rootCmd.AddCommand(testCmd)

}
