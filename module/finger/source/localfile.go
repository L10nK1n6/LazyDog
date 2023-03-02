package source

import (
	"bufio"
	"github.com/gookit/color"
	"log"
	"os"
	"strings"
)

func LocalFile(filename string) (urls []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("本地文件读取失败:", err)
		color.RGBStyleFromString("238,99,99").Println("[error] 输入文件错误!!!")
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "http") {
			urls = append(urls, scanner.Text())
		} else {
			urls = append(urls, "https://"+scanner.Text())
		}
	}
	return
}
