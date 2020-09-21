package main

import (
	"encoding/json"
	"fmt"
	"netspace/config"
	"netspace/mq"
	"os"
)

func ProcessTransfer(msg []byte) bool {
	// 解析msg
	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		fmt.Println("err:", err)
		return false
	}
	// 根据临时存储路径,创建文件句柄
	_, err = os.Open(pubData.CurLocation)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// 根据文件句柄将文件内容读出来并上传到oss
	/*
		oss.Bucket().PutObject(pubData.DesLocation,bufio.NewReader(file))
		if err!=nil{
			panic(err)
		}
	*/
	// TODO: 上传到OSS
	// 更新文件的存储路径到文件表
	// TODO: UpdateFIleLocation
	return true
}

func main() {
	fmt.Println("开始监听转移任务...")
	mq.StartConsumer(config.TransExchangeName, "transfer_oss", ProcessTransfer)
}
