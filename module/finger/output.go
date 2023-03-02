package finger

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gookit/color"
	"os"
	"strconv"
	"strings"
)

func outjson(kwpath string, filename string, data []byte) {
	filepath := "./result/" + kwpath + "/"
	f, err := os.Create(filepath + filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		defer f.Close()
		_, err = f.Write(data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func outxlsx(kwpath string, filename string, msg []Outrestul) {
	filepath := "./result/" + kwpath + "/"
	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "url")
	xlsx.SetCellValue("Sheet1", "B1", "cms")
	xlsx.SetCellValue("Sheet1", "C1", "server")
	xlsx.SetCellValue("Sheet1", "D1", "statuscode")
	xlsx.SetCellValue("Sheet1", "E1", "length")
	xlsx.SetCellValue("Sheet1", "F1", "title")
	for k, v := range msg {
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(k+2), v.Url)
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(k+2), v.Cms)
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(k+2), v.Server)
		xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(k+2), v.Statuscode)
		xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(k+2), v.Length)
		xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(k+2), v.Title)
	}
	err := xlsx.SaveAs(filepath + filename)
	if err != nil {
		fmt.Println(err)
	}
}

func outcsv(kwpath string, filename string, msg []Outrestul) {
	filepath := "./result/" + kwpath + "/"
	f, err := os.Create(filepath + filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，防止csv中文乱码
	writer := csv.NewWriter(f)
	writer.Write([]string{"url", "cms", "server", "statuscode", "length", "title"})

	for _, v := range msg {
		writer.Write([]string{v.Url, v.Cms, v.Server, strconv.Itoa(v.Statuscode), strconv.Itoa(v.Length), v.Title})
	}

	writer.Flush() // 此时才会将缓冲区数据写入

}

func outfile(filename string, allresult []Outrestul, kwpath string) {
	file := strings.Split(filename, ".")
	if len(file) == 2 {
		switch file[1] {
		case "csv":
			outcsv(kwpath, filename, allresult)
		case "json":
			buf, err := json.MarshalIndent(allresult, "", " ")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			outjson(kwpath, filename, buf)
		case "xlsx":
			outxlsx(kwpath, filename, allresult)
		default:
			color.RGBStyleFromString("238,99,99").Println("\n文件未保存，不支持该类型！")
		}
	}

}
