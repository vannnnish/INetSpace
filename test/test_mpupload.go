package main

import (
	"bufio"
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func multipartUpload(filename string, targetURL string, chunkSize int) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	bfRd := bufio.NewReader(f)
	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize)
	for {
		n, err := bfRd.Read(buf)
		if n <= 0 {
			break
		}
		index++
		bufCopied := make([]byte, 5*1048576)
		copy(bufCopied, buf)
		go func(b []byte, curIdx int) {
			fmt.Printf("upload_size%d\n", len(b))

			resp, err := http.Post(targetURL+"&index="+strconv.Itoa(curIdx), "multipart/form-data", bytes.NewReader(b))
			if err != nil {
				fmt.Println(err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			fmt.Printf("%+v%+v\n", string(body), err)
			resp.Body.Close()
			ch <- curIdx

		}(bufCopied[:n], index)

		// 遇到任何错误,立刻返回
		if err != nil {
			if err == io.EOF {
			}
		}
	}
	for idx := 0; idx < index; idx++ {
		select {
		case res := <-ch:
			fmt.Println(res)
		}
	}
	return nil
}

// 测试分块上传脚本
func main() {
	username := "admin"
	token := "3e3c525c0e097e73c99411c9f53cd8045f6563aa"
	filehash := "81bfaf7e1c9f1f35a36d07b6137b456397d394f5"
	fileSize := "121149509"
	uploadAPI := "http://localhost:8080/file/mpupload/init"
	// 请求初始化上传接口
	resp, err := http.PostForm(uploadAPI, url.Values{
		"username": {username},
		"token":    {token},
		"filehash": {filehash},
		"filesize": {fileSize},
	})

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 得到uploadID 以及服务端指定的分块大小 chunkSize
	uploadID := jsoniter.Get(body, "data").Get("UploadID").ToString()
	chunkSize := jsoniter.Get(body, "data").Get("ChunkSize").ToInt()
	fmt.Printf("uploadid:%s  chunksize:%d\n", uploadID, chunkSize)

	// 请求分块上传接口
	filename := "./go1.15.2.linux-amd64.tar.gz"
	tURL := "http://localhost:8080/file/mpupload/uppart?" + "username=" + username + "&token=" + token + "&uploadid=" + uploadID
	multipartUpload(filename, tURL, chunkSize)

	// 请求分块完成接口
	resp, err = http.PostForm("http://localhost:8080/file/mpupload/complete", url.Values{
		"username": {username},
		"token":    {token},
		"filehash": {filehash},
		"filesize": {fileSize},
		"filename": {"go1.15.2.linux-amd64.tar.gz"},
		"uploadid": {uploadID},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))

}
