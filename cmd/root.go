package cmd

import (
	"fmt"
	"os"

	"github.com/fliropp/gokube/cmd/server"
	"github.com/spf13/cobra"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use: "root",
}

var TestCmd = &cobra.Command{
	Use:   "are2",
	Short: "Ares World cmd 2",
	Run: func(cmd *cobra.Command, args []string) {
		for _, b := range args {
			fmt.Println(b)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(TestCmd)
	RootCmd.AddCommand(server.ServerCmd)
}
