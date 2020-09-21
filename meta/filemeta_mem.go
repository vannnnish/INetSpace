package meta

import (
	"os"
)

// FileMeta:文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

//  UpdateFileMeata:新增/更新文件元信息
func UpdateFileMeta(fMeta FileMeta) {
	fileMetas[fMeta.FileSha1] = fMeta
}

func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetLastFileMetas: 批量获取文件列表
// TODO:必定要优化的
func GetLastFileMeta(count int) []FileMeta {
	metas := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		metas = append(metas, v)
	}
	// 排序
	return metas[0:count]
}

func RemoveFileMeta(fileSha1 string) error {
	fMeta := fileMetas[fileSha1]
	err := os.Remove(fMeta.Location)
	delete(fileMetas, fileSha1)
	return err
}

