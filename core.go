package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/gammazero/workerpool"
)

const (
	CODES         = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	CODES_LEN     = len(CODES)
	USER_AGENT    = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0"
	FILE_NAME_LEN = 20         // 文件名字符串长度
	POOL_MAXSIZE  = 16         // goroutine 池容量
	PICS_EXT      = ".jpg"     // 图片后缀
	PICS_DIR      = "pics"     // 存放图片文件夹
	URLS_DATA     = "data.txt" // url 数据来源
	RETRY_TIMES   = 3
)

// 初始化操作
func init() {
	createDir(PICS_DIR)
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// 如果文件夹不存在，则创建文件夹
func createDir(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsExist(err) {
			os.Mkdir(path, os.ModePerm)
			fmt.Println("Create pics dir.")
		}
	}
}

// 返回请求响应内容
func getResponse(url string) (*http.Response, error) {
	var ref string
	// 根据 url 确定 header Referer 字段
	if strings.HasPrefix(url, "http://i.meizitu.net/") {
		ref = "http://www.mzitu.com"
	}
	if strings.HasPrefix(url, "http://img.mmjpg.com/") {
		ref = "http://www.mmjpg.com"
	}
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", USER_AGENT)
	request.Header.Set("Referer", ref)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 下载图片
func downloadPics(url string) {
	fileName := randStr() + PICS_EXT
	localFile, err := os.Create(path.Join(PICS_DIR, fileName))
	if err != nil {
		fmt.Println("Failed opening local file: " + err.Error())
		return
	}
	defer localFile.Close()
	fmt.Println("Download pics:", fileName)
	var resp *http.Response
	for retry := 0; retry < RETRY_TIMES; retry++ {
		resp, err = getResponse(url)
		if err == nil && resp.StatusCode == 200 {
			break
		}
	}
	if resp == nil || resp.StatusCode != 200 {
		fmt.Println("Failed requesting " + url)
		return
	}
	if _, err := io.Copy(localFile, resp.Body); err != nil {
		fmt.Println(err)
	}
}

// 返回随机字符串，用作函数名
func randStr() string {
	data := make([]byte, FILE_NAME_LEN)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < FILE_NAME_LEN; i++ {
		idx := rand.Intn(CODES_LEN)
		data[i] = byte(CODES[idx])
	}

	return string(data)
}

// 主函数
func main() {
	// 创建 goroutine 池
	wp := workerpool.New(POOL_MAXSIZE)
	start := time.Now()
	f, _ := os.Open(URLS_DATA)
	scanner := bufio.NewScanner(bufio.NewReader(f))
	urlIsRead := make(chan bool)
	for scanner.Scan() {
		wp.Submit(func() {
			url := scanner.Text()
			urlIsRead <- true
			downloadPics(url)
		})
		select {
		case <-urlIsRead:
			break
		}
	}
	// 等待所有任务完成
	wp.StopWait()
	elapsed := time.Since(start)
	fmt.Println("Elapsed: ", elapsed)
}
