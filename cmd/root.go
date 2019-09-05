package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	// config file path
	cfgFile     string
	versionFlag bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "p2p-sharer",
	Short: "a p2p sharer",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}