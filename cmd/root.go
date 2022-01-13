package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	aw "github.com/deanishe/awgo"
)

const CmdFlag = "cmd"

func init() {
	// Create a new Workflow using default settings.
	// Critical settings are provided by Alfred via environment variables,
	// so this *will* die in flames if not run in an Alfred-like environment.
	wf = aw.New()
	rootCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list kit",
		Run: func(cmd *cobra.Command, args []string) {
			wf.Run(func() {
				for _, action := range actionMap {
					action.ActionItem().Var(CmdFlag, action.Use())
				}
				wf.SendFeedback()
			})
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "exec",
		Run: func(cmd *cobra.Command, args []string) {
			wf.Run(func() {
				cmd := wf.Config.Get(CmdFlag)
				if len(cmd) == 0 {
					wf.NewWarningItem("Invalid CMD", fmt.Sprintf("not found [%s] in vars", CmdFlag))
				} else {
					actionMap[cmd].Action(args)
				}
				wf.SendFeedback()
			})
		},
	})
}

var rootCmd = &cobra.Command{
	Use:   "awpark",
	Short: "alfred workflow park",
}

var wf *aw.Workflow

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var actionMap = make(map[string]Command)

func Register(cmd Command) {
	if _, ok := actionMap[cmd.Use()]; ok {
		panic(fmt.Sprintf("%s: is duplicate"))
	}
	actionMap[cmd.Use()] = cmd
}

type Command interface {
	Use() string
	ActionItem() *aw.Item
	Action(args []string)
}
