package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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

func UploadHandler(c *gin.Context) {
	//  返回上传HTML页面
	c.Redirect(http.StatusFound, "./static/view/index.html")
}

func DoUploadHandler(c *gin.Context) {
	// 接受文件流
	file, head, err := c.Request.FormFile("file")
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
	username := c.Request.FormValue("username")
	ok := db.OnUserUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
	if !ok {
		fmt.Println("上传失败")
		return
	}
	//meta.UpdateFileMeta(fileMeta)
	meta.UpdateFileMetaDB(fileMeta)
	c.JSON(http.StatusOK, "上传成功啦")
}

// UploadSucHandler:上传完成
func UploadSucHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "Upload Success")
}

// GetFileMetaHandler:获取上传的列表
func GetFileMetaHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	fileMeta, err := meta.GetFileMetaDb(filehash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fileMeta)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

// 批量查询文件元信息
func FileQueryHandler(c *gin.Context) {
	limitCnt, _ := strconv.Atoi(c.Request.FormValue("limit"))
	username := c.Request.FormValue("username")
	userFiles, err := db.QueryUserFileMeta(username, limitCnt)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(userFiles)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

func Downloadhandler(c *gin.Context) {
	fsha1 := c.Request.FormValue("filehash")
	fileMeta := meta.GetFileMeta(fsha1)
	//
	file, err := os.Open(fileMeta.Location)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	//
	c.Header("Content-Type", "application/octect-stream")
	c.Header("content-disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
	c.JSON(http.StatusOK, data)
}

// f
func FileUpdateMetaHandler(c *gin.Context) {
	opType := c.Request.FormValue("op")
	fileHash := c.Request.FormValue("filehash")
	newFName := c.Request.FormValue("filename")

	if opType != "0" {
		c.Status(http.StatusInternalServerError)
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
		c.Status(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

// FileDeleteHandler：删除文件
func FileDelHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	err := meta.RemoveFileMeta(filehash)
	if err != nil {
		fmt.Println("err:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, "删除成功")
}

func TryFastUploadHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))
	//
	fileMeta, err := db.GetFileMeta(filehash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	//
	if fileMeta == nil {
		msg := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败",
			Data: nil,
		}
		c.JSON(http.StatusOK, msg.JSONBytes())
		return
	}
	finished := db.OnUserUploadFinished(username, fileMeta.FileHash, filename, int64(filesize))
	if finished {
		res := util.RespMsg{
			Code: 0,
			Msg:  "妙传成功",
			Data: nil,
		}
		c.JSON(http.StatusOK, res.JSONBytes())
		return
	} else {
		res := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败,稍后重试",
			Data: nil,
		}
		c.JSON(http.StatusOK, res.JSONBytes())
	}
}
