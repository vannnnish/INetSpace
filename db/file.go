package db

import (
	"database/sql"
	"fmt"
	"netspace/db/mysql"
)

// OnfileUploadFinished 上传成功
func OnfileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_file(`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`)values(?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement ,err", err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	affected, err := ret.RowsAffected()
	if err == nil {
		if affected <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize int64
	FileAddr sql.NullString
}

// GetFileMeta 从mysql 获取原信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stat, err := mysql.DBConn().Prepare(
		" select file_sha1,file_addr,file_name,file_size from tbl_file where file_sha1= ? and " +
			"status = 1 limit 1")

	if err != nil {
		return nil, err
	}
	defer stat.Close()
	tfile := &TableFile{}
	err = stat.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		return nil, err
	}
	return tfile, nil

}
