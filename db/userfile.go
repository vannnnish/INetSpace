package db

import (
	"fmt"
	"netspace/db/mysql"
	"time"
)

// 用户文件表结构体
type UserFile struct {
	UserName   string
	FileHash   string
	FileName   string
	FileSize   int64
	UploadAt   string
	LastUpdate string
}

func OnUserUploadFinished(username, fileHash, filename string, fileSize int64) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into " +
			"tbl_user_file(" +
			"`user_name`," +
			"`file_sha1`," +
			"`file_name`," +
			"`file_size`," +
			"`upload_at`) values (?,?,?,?,?)")
	if err != nil {
		fmt.Println("创建文件记录失败:", err)
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, fileHash, filename, fileSize, time.Now())
	if err != nil {
		fmt.Println("创建文件记录失败", err)
		return false
	}
	return true
}

// 获取用户信息
func QueryUserFileMeta(username string, limit int) ([]UserFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"select " +
			"file_sha1," +
			"file_name," +
			"file_size," +
			"upload_at," +
			"last_update " +
			"from tbl_user_file " +
			"where user_name = ? limit ?")
	if err != nil {
		fmt.Println("err,", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}
	var userFiles []UserFile
	for rows.Next() {
		uFile := UserFile{}
		err = rows.Scan(&uFile.FileHash, &uFile.FileName, &uFile.FileSize, &uFile.UploadAt, &uFile.LastUpdate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		userFiles = append(userFiles, uFile)
	}
	return userFiles, nil
}
