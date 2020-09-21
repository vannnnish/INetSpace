package meta

import (
	"fmt"
	"netspace/db"
)

func UpdateFileMetaDB(fmeta FileMeta) bool {
	return db.OnfileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// GetFileMetaDb 从mysql 获取文件源信息
func GetFileMetaDb(fileSha1 string) (FileMeta, error) {
	tfile, err := db.GetFileMeta(fileSha1)
	if err != nil {
		fmt.Println("err:", err)
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil
}
