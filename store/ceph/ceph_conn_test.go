package ceph

import (
	"fmt"
	"gopkg.in/amz.v1/s3"
	"testing"
)

func TestGetCephBucket(t *testing.T) {
	bucket := GetCephBucket("testbucket1")

	err := bucket.PutBucket(s3.PublicRead)
	if err != nil {
		panic(err)
	}
	list, err := bucket.List("", "", "", 100)
	fmt.Println("object key:", list)
	// 创建一个新的bucket
	// 新上传一个对象
	err = bucket.Put("/testupload/a.txt", []byte("just fot test"), "octet-stream", s3.PublicRead)
	if err != nil {
		panic(err)
	}
	// 查询bucket下面指定条件的object keys
	list, err = bucket.List("", "", "", 100)
	fmt.Println("object key:", list)

}
