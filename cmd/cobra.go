package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:               "exporter-center",
	Short:             "-v",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Long:              `exporter-center`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		usageStr := `欢迎使用 exporter-center，可以使用 -h 查看命令`
		fmt.Println("%s\n", usageStr)
	},
}

func init() {
	// 启动服务
	rootCmd.AddCommand(StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
