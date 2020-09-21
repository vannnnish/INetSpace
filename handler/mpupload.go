package handler

import (
	"fmt"
	redis2 "github.com/garyburd/redigo/redis"
	"math"
	"net/http"
	"netspace/cache/redis"
	"netspace/db"
	"netspace/util"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// MultipartUpload
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

func InitialMultipartUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 获得一个redis连接
	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	//	生成一个分块上传的初始化信息
	uploadInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, // 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}
	// 将初始化的信息写入redis
	rConn.Do("HSET", "MP_"+uploadInfo.UploadID, "chunkcount", uploadInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+uploadInfo.UploadID, "filehash", uploadInfo.FileHash)
	rConn.Do("HSET", "MP_"+uploadInfo.UploadID, "filesize", uploadInfo.FileSize)
	// 将响应数据返回给客户端
	w.Write((&util.RespMsg{Code: 0, Msg: "OK", Data: uploadInfo}).JSONBytes())
}

// 上传文件分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 解析用户参数
	r.ParseForm()
	//username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")
	// 获取redis连接
	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	//	先创建目录
	fPath := "./tmp/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fPath), 0744)
	// 获取文件句柄,用于存储分块内容
	fd, err := os.Create(fPath)
	if err != nil {
		fmt.Println("创建文件错误:", err)
		w.Write(util.NewRespMsg(-1, "upload fail", nil).JSONBytes())
		return
	}
	defer fd.Close()
	buf := make([]byte, 1024*1024)
	// TODO: 可以加CRC 校验
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 更新redis 缓存数据  这里是,根据分块的id, 判断该分块的数据是否完成.
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)
	// 返回处理结果给客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// 通知上传合并接口

func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	r.ParseForm()
	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	// 获取redis连接
	rConn := redis.RedisPool().Get()
	defer rConn.Close()
	// 是否所有分块上传完成
	data, err := redis2.Values(rConn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	length := len(data)
	for i := 0; i < length; i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-1, "invalid request", nil).JSONBytes())
		return
	}
	// TODO: 所有的分块上传完成后,合并分块
	// 更新唯一文件表, 更新用户文件表
	fsize, _ := strconv.Atoi(filesize)
	db.OnfileUploadFinished(filehash, filename, int64(fsize), "")
	db.OnUserUploadFinished(username, filehash, filename, int64(fsize))
	// 响应客户端处理结果
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// 取消分块上传
func CancelUploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 删除已存在的分块文件
	// 删除redis缓存状态
	// 更新mysql status文件

}

// 查看分块上传整体状态
func MultipartUploadStatusHandler(w http.ResponseWriter, r *http.Request) {
	// 检查分块上传状态是否有效
	// 获取分块初始化信息
	// 获取已上传的分块信息
}
