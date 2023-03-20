package source

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gookit/color"
	"github.com/spf13/viper"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type HunterResult struct {
	Code int `json:"code"`
	Data struct {
		AccountType string `json:"account_type"`
		Total       int    `json:"total"` //资产总数
		Time        int    `json:"time"`
		Arr         []struct {
			IsRisk         string      `json:"is_risk"`
			URL            string      `json:"url"` //带端口URL
			IP             string      `json:"ip"`
			Port           int         `json:"port"`      //端口
			WebTitle       string      `json:"web_title"` //标题
			Domain         string      `json:"domain"`    //域名
			IsRiskProtocol string      `json:"is_risk_protocol"`
			Protocol       string      `json:"protocol"` //协议
			BaseProtocol   string      `json:"base_protocol"`
			StatusCode     int         `json:"status_code"` //网站状态码
			Component      interface{} `json:"component"`
			Os             string      `json:"os"`
			Company        string      `json:"company"` //备案单位
			Number         string      `json:"number"`  //备案号
			Country        string      `json:"country"`
			Province       string      `json:"province"` //省份
			City           string      `json:"city"`     //市区
			UpdatedAt      string      `json:"updated_at"`
			IsWeb          string      `json:"is_web"`
			AsOrg          string      `json:"as_org"`
			Isp            string      `json:"isp"` //运营商信息
			Banner         string      `json:"banner"`
		} `json:"arr"`
		ConsumeQuota string `json:"consume_quota"` //消耗积分
		RestQuota    string `json:"rest_quota"`    //剩余积分
		SyntaxPrompt string `json:"syntax_prompt"`
	} `json:"data"`
	Message string `json:"message"`
}

// 获取当前执行程序所在的绝对路径
func GetCurrentAbPathByExecutable() string {
	//获取项目的执行路径
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}

type Config struct {
	Hunter struct {
		Apikey string //定义配置文件参数
	}
}

// 配置数据变量
var ConfigData *Config

func init() {

	//导入配置文件
	viper.AddConfigPath(GetCurrentAbPathByExecutable()) //设置读取的文件路径
	viper.SetConfigName("config")                       //设置读取的文件名
	viper.SetConfigType("yaml")

	//读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("配置文件错误：", err.Error())
	}
	err = viper.Unmarshal(&ConfigData)
	if err != nil {
		fmt.Println("配置文件数据错误", err.Error())
	}

}

// GetCOnfigData 返回配置数据方法
func GetCOnfigData() *Config {
	return ConfigData
}

//请求hunter-api过程
func HunterAPI(searchkeyword string, searhtype string) (urls []string) {
	color.RGBStyleFromString("255,165,0").Println("请耐心等待Hunter搜索......")
	huntercfg := GetCOnfigData() //初始化配置文件

	fmt.Println("Hunter APIKEY为：" + huntercfg.Hunter.Apikey) //打印APIKEY

	var keywords string
	var searchkeyworddir string
	switch searhtype {
	case "domain":
		keywords = `domain.suffix="` + searchkeyword + `"` //hunter查询域名的语法
		searchkeyworddir = searchkeyword
	case "ip":
		keywords = `ip="` + searchkeyword + `"` //查询ip的语法
		StrContainers := strings.Contains(searchkeyword, "/24")
		if StrContainers {
			searchkeyworddir = strings.Replace(searchkeyword, "/24", "C", -1)
		}else {
			searchkeyworddir = searchkeyword
		}


		//预留未来备案号和备案单位的类型
	}

	b64query := base64.StdEncoding.EncodeToString([]byte(keywords))           //base64编码
	month := time.Now().Format("01")                                          //年月日格式
	day := time.Now().Format("02")                                            //年月日格式
	start_time := strconv.Itoa(time.Now().Year()-1) + "-" + month + "-" + day //开始日期
	end_time := time.Now().Format("2006") + "-" + month + "-" + day           //结束日期

	//初始值
	page := 1
	maxsize := 100

	//打印请求URL
	fmt.Println("请求URL：https://hunter.qianxin.com/openApi/search?api-key=" + huntercfg.Hunter.Apikey + "&search=" + b64query + "&page=" + strconv.Itoa(page) + "&page_size=100&is_web=3&port_filter=false&start_time=" + start_time + "&end_time=" + end_time)
	//resty发送请求包
	client := resty.New()
	HunterResultS := &HunterResult{}
	client.R().SetResult(HunterResultS).Get("https://hunter.qianxin.com/openApi/search?api-key=" + huntercfg.Hunter.Apikey + "&search=" + b64query + "&page=" + strconv.Itoa(page) + "&page_size=100&is_web=3&port_filter=false&start_time=" + start_time + "&end_time=" + end_time)

	//创建目录、创建标记文件
	flagtext := "flag"
	RealAssetsfile := "/flag.txt"
	filepath := "./Result/" + searchkeyworddir
	filepathfull := filepath + RealAssetsfile
	if HunterResultS.Code == 200 {
		CPath(filepath, filepathfull, flagtext)
	} else {
		color.RGBStyleFromString("238,99,99").Println("\nHunter API错误：" + HunterResultS.Message)
		os.Exit(0)
	}

	//保存到文件
	f, err := os.Create(filepath + "/HunterResult.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，防止csv中文乱码
	writer := csv.NewWriter(f)
	writer.Write([]string{"URL", "IP", "端口", "协议", "标题", "域名", "网站状态码", "备案单位", "备案号", "省份", "市区", "运营商信息"})

	//打印数据
	for _, lib := range HunterResultS.Data.Arr {
		fmt.Println("IP为：", lib.IP, "\n端口为：", lib.Port, "\nURL为：", lib.URL, "\n---------") //信息打印
		//保存hunter数据
		writer.Write([]string{lib.URL, lib.IP, strconv.Itoa(lib.Port), lib.Protocol, lib.WebTitle, lib.Domain, strconv.Itoa(lib.StatusCode), lib.Company, lib.Number, lib.Province, lib.City, lib.Isp})
		//加入到urls后面丢给指纹识别
		urls = append(urls, lib.URL)
	}

	//翻页
	for HunterResultS.Data.Total > maxsize {
		//控制页数
		page++
		maxsize = maxsize + 100

		time.Sleep(1 * time.Second)
		//打印请求URL
		fmt.Println("请求URL：https://hunter.qianxin.com/openApi/search?api-key=" + huntercfg.Hunter.Apikey + "&search=" + b64query + "&page=" + strconv.Itoa(page) + "&page_size=100&is_web=3&port_filter=false&start_time=" + start_time + "&end_time=" + end_time)
		//resty发送请求包
		client.R().SetResult(HunterResultS).Get("https://hunter.qianxin.com/openApi/search?api-key=" + huntercfg.Hunter.Apikey + "&search=" + b64query + "&page=" + strconv.Itoa(page) + "&page_size=100&is_web=3&port_filter=false&start_time=" + start_time + "&end_time=" + end_time)

		if HunterResultS.Code != 200 {
			color.RGBStyleFromString("238,99,99").Println("\nHunter API错误：" + HunterResultS.Message)
		}

		for _, lib := range HunterResultS.Data.Arr {
			fmt.Println("IP为：", lib.IP, "\n端口为：", lib.Port, "\nURL为：", lib.URL, "\n---------") //信息打印
			//保存hunter数据
			writer.Write([]string{lib.URL, lib.IP, strconv.Itoa(lib.Port), lib.Protocol, lib.WebTitle, lib.Domain, strconv.Itoa(lib.StatusCode), lib.Company, lib.Number, lib.Province, lib.City, lib.Isp})
			//加入到urls后面丢给指纹识别
			urls = append(urls, lib.URL)
		}
	}

	fmt.Println("资产总数为：", HunterResultS.Data.Total, "\n", HunterResultS.Data.RestQuota) //积分打印

	writer.Flush() // 此时才会将缓冲区数据写入

	return
}

//判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//检查文件是否存在
func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

//保存文件操作
func Sfile(filepath string, writeString string) {
	var f *os.File
	var err1 error
	if checkFileIsExist(filepath) { //如果文件存在
		f, err1 = os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0777) //打开文件
		fmt.Println("文件已存在，正在覆盖...")
	} else {
		f, err1 = os.Create(filepath) //创建文件
		fmt.Println("文件不存在，正在创建...")
	}
	defer f.Close()
	n, err1 := io.WriteString(f, writeString) //写入文件(字符串)
	if err1 != nil {
		panic(err1)
	}
	fmt.Printf("写入 %d 个字节，文件路径：%s \n", n, filepath)
}

//创建文件夹操作
func CPath(filepath string, filepathfull string, savefile string) {
	//判断目录是否存在
	exist, err := PathExists(filepath)
	if err != nil {
		fmt.Printf("目录错误![%v]\n", err)
		return
	}
	if exist {
		fmt.Printf("目录存在![%v]\n", filepath)
		Sfile(filepathfull, savefile)
	} else {
		fmt.Printf("目录不存在![%v]\n", filepath)
		// 创建文件夹
		err := os.Mkdir(filepath, os.ModePerm)
		if err != nil {
			color.RGBStyleFromString("238,99,99").Printf("目录创建失败![%v]\n", err)
		} else {
			color.RGBStyleFromString("144,238,144").Printf("目录创建成功!\n")
		}
		Sfile(filepathfull, savefile)
	}
}
