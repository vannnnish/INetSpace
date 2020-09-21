package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"netspace/db"
	"netspace/meta"
	"netspace/util"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//  返回上传HTML页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 接受文件流
		file, head, err := r.FormFile("file")
		defer file.Close()
		if err != nil {
			fmt.Println("Failed to get data ,err:", err.Error())
			panic(err)
		}
		fileMeta := meta.FileMeta{
			FileSha1: "",
			FileName: head.Filename,
			FileSize: 0,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		fmt.Println("文件:", fileMeta.Location)
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file ,err:%s\n", err.Error())
			panic(err)
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Println("Failed to save data into file , err:", err.Error())
		}
		// Seek设置下一次读/写的位置。offset为相对偏移量，而whence决定相对位置：0为相对文件开头，1为相对当前位置，2为相对文件结尾。它返回新的偏移量（相对开头）和可能的错误。
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		// TODO 更新用户文件记录表
		username := r.Form.Get("username")
		ok := db.OnUserUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if !ok {
			fmt.Println("上传失败")
			return
		}
		//meta.UpdateFileMeta(fileMeta)
		meta.UpdateFileMetaDB(fileMeta)
		io.WriteString(w, "上传成功啦")
	}
}

// UploadSucHandler:上传完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload Success")
}

// GetFileMetaHandler:获取上传的列表
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fileMeta, err := meta.GetFileMetaDb(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 批量查询文件元信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	userFiles, err := db.QueryUserFileMeta(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func Downloadhandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(fsha1)
	//
	file, err := os.Open(fileMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
	w.Write(data)
}

// f
func FileUpdateMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opType := r.Form.Get("op")
	fileHash := r.Form.Get("filehash")
	newFName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var (
		oldPath string
	)
	curFileMeta := meta.GetFileMeta(fileHash)
	oldPath = curFileMeta.Location
	curFileMeta.FileName = newFName
	curFileMeta.Location = filepath.Dir(oldPath) + "/" + curFileMeta.FileName
	meta.UpdateFileMeta(curFileMeta)
	//
	err := os.Rename(oldPath, curFileMeta.Location)
	if err != nil {
		fmt.Println("err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileDeleteHandler：删除文件
func FileDelHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	err := meta.RemoveFileMeta(filehash)
	if err != nil {
		fmt.Println("err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("删除成功"))
}

func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))
	//
	fileMeta, err := db.GetFileMeta(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	if fileMeta == nil {
		msg := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败",
			Data: nil,
		}
		w.Write(msg.JSONBytes())
		return
	}
	finished := db.OnUserUploadFinished(username, fileMeta.FileHash, filename, int64(filesize))
	if finished {
		res := util.RespMsg{
			Code: 0,
			Msg:  "妙传成功",
			Data: nil,
		}
		w.Write(res.JSONBytes())
		return
	} else {
		res := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败,稍后重试",
			Data: nil,
		}
		w.Write(res.JSONBytes())
	}

}
