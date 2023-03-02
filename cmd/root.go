package cmd

import (
	"github.com/spf13/cobra"
)

var cfgFile string
var outputFile string

var rootCmd = &cobra.Command{
	Use:   "lazydog",
	Short: "LazyDog是一款通过Hunter-API快速获取目标单位的网络资产并使用Ehole进行指纹分析的工具。",
	Long: " _                     _____              \n" +
		"| |                   |  __ \\             \n" +
		"| |     __ _ _____   _| |  | | ___   __ _  \n" +
		"| |    / _` |_  / | | | |  | |/ _ \\ / _` | \n" +
		"| |___| (_| |/ /| |_| | |__| | (_) | (_| | \n" +
		"|______\\__,_/___|\\__, |_____/ \\___/ \\__, | \n" +
		"                  __/ |              __/ |\n" +
		" by:M@tri(x_o)   |___/              |___/ ",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
