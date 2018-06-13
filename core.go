package main

import (
	"bufio"
	"fmt"
	"github.com/gammazero/workerpool"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

const (
	CODES         = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	CODES_LEN     = len(CODES)
	USER_AGENT    = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0"
	FILE_NAME_LEN = 20         // 文件名字符串长度
	POOL_MAXSIZE  = 128        // goroutine 池容量
	PICS_EXT      = ".jpg"     // 图片后缀
	PICS_DIR      = "pics"     // 存放图片文件夹
	URLS_DATA     = "data.txt" // url 数据来源
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
func getResponse(url string) *http.Response {
	var ref string
	// 根据 url 确定 header Referer 字段
	if strings.HasPrefix(url, "http://i.meizitu.net/") {
		ref = "http://www.mzitu.com"
	}
	if strings.HasPrefix(url, "http://img.mmjpg.com/") {
		ref = "http://www.mmjpg.com"
	}
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", USER_AGENT)
	request.Header.Set("Referer", ref)
	response, _ := client.Do(request)
	return response
}

// 下载图片
func downloadPics(url string) {
	fileName := randStr() + PICS_EXT
	localFile, _ := os.Create(path.Join(PICS_DIR, fileName))
	fmt.Println("Download pics:", fileName)
	if _, err := io.Copy(localFile, getResponse(url).Body); err != nil {
		fmt.Println(err)
	}
	defer localFile.Close()
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
	for scanner.Scan() {
		wp.SubmitWait(func() { downloadPics(scanner.Text()) })
	}
	// 等待所有任务完成
	wp.Stop()
	elapsed := time.Since(start)
	fmt.Println("Elapsed: ", elapsed)
}
