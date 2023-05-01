package common

import (
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"os"
)

var InputTargetsFileName string
var TokenFileName string
var OutputFileName string
var Threads int
var ClientVersion string

func Banner() {
	banner := `
  __             ____                _             
 / _|_ __ _ __  / ___|_ __ __ _  ___| | _____ _ __ 
| |_| '__| '_ \| |   | '__/ _` + "`" + `|/ __| |/ / _ \ '__|
|  _| |  | |_) | |___| | | (_| | (__|   <  __/ |
|_| |_|  | .__/ \____|_|  \__,_|\___|_|\_\___|_|
         |_|                  by: SleepingBag945
`
	print(banner)
}

func Flag() {
	Banner()
	flag.StringVar(&InputTargetsFileName, "l", "", "输入文件，支持IP:Port格式，一行一个")
	flag.StringVar(&TokenFileName, "tl", "", "Token的输入文件，一行一个")
	flag.StringVar(&OutputFileName, "o", "result.txt", "输出文件")
	flag.IntVar(&Threads, "t", 20, "线程数量")
	flag.StringVar(&ClientVersion, "v", "0.48.0", "指定爆破时客户端的版本")
	flag.Parse()
}

func WriteFile(result string, filename string) {
	var text = []byte(result)
	fl, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Open %s error, %v\n", filename, err)
		return
	}
	_, err = fl.Write(text)
	fl.Close()
	if err != nil {
		fmt.Printf("Write %s error, %v\n", filename, err)
	}
}

func WriteResult(result string) {
	fmt.Print(aurora.BrightRed(result))
	WriteFile(result, OutputFileName)
}
