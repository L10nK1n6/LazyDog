package cmd

import (
	"LazyDog/module/finger"
	"LazyDog/module/finger/source"
	"github.com/spf13/cobra"
	"os"
)

var (
	hunterip     string
	hunterdomain string
	localfile    string
	urla         string
	thread       int
	output       string
	proxy        string
)

var fingerCmd = &cobra.Command{
	Use:   "hunter",
	Short: "Hunter查询模块",
	Long:  `通过Hunter或本地文件导入获取资产进行指纹识别，支持单条url识别。`,
	Run: func(cmd *cobra.Command, args []string) {
		if localfile != "" {
			urls := removeRepeatedElement(source.LocalFile(localfile))
			s := finger.NewScan(urls, thread, output, proxy)
			s.StartScan("")
			os.Exit(1)
		}
		if hunterip != "" {
			urls := removeRepeatedElement(source.HunterAPI(hunterip, "ip"))
			s := finger.NewScan(urls, thread, output, proxy)
			s.StartScan(hunterip)
			os.Exit(1)
		}
		if hunterdomain != "" {
			urls := removeRepeatedElement(source.HunterAPI(hunterdomain, "domain"))
			s := finger.NewScan(urls, thread, output, proxy)
			s.StartScan(hunterdomain)
			os.Exit(1)
		}
		if urla != "" {
			s := finger.NewScan([]string{urla}, thread, output, proxy)
			s.StartScan("")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(fingerCmd)
	fingerCmd.Flags().StringVarP(&hunterip, "ip", "i", "", "从Hunter提取IP资产进行指纹识别，支持ip或者ip段，例如：192.168.1.1 | 192.168.1.0/24")
	fingerCmd.Flags().StringVarP(&hunterdomain, "domain", "d", "", "从Hunter提取域名资产进行指纹识别")
	fingerCmd.Flags().StringVarP(&localfile, "local", "l", "", "从本地文件读取资产，进行指纹识别，支持无协议，例如：192.168.1.1:9090 | http://192.168.1.1:9090")
	fingerCmd.Flags().StringVarP(&urla, "url", "u", "", "识别单个目标。")
	fingerCmd.Flags().StringVarP(&output, "output", "o", "", "输出所有结果，支持csv、json和xlsx后缀的文件。")
	fingerCmd.Flags().IntVarP(&thread, "thread", "t", 100, "指纹识别线程大小。")
	fingerCmd.Flags().StringVarP(&proxy, "proxy", "p", "", "指定访问目标时的代理，支持http代理和socks5，例如：http://127.0.0.1:8080 | socks5://127.0.0.1:8080")
}

//去重
func removeRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
