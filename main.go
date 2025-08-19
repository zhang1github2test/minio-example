package main

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func main() {
	endpoint := "47.98.132.219:9000"
	accessKeyID := "8OTHJPNJOAYLX6S9RQB9"
	secretAccessKey := "Wc5dMy1UKIEcHpHMWpOiFL4FfZ6Az9Tx7njcqPrl"
	useSSL := false
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("初始minio客户端失败", err)
	}
	log.Println("初始minio客户端成功")
	// create bucket
	bucketName := "test-bucket-go"

	// check bucket exists
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalln("检查bucket是否存在失败", err)
	}
	if exists {
		log.Println("bucket已存在")
	} else {
		// create bucket
		err := client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln("创建bucket失败", err)
		}
		log.Println("创建bucket成功")
	}
	// upload object
	fileName := "1.txt"
	objectName := "go/1.txt"
	object, err := client.FPutObject(context.Background(), bucketName, objectName, fileName, minio.PutObjectOptions{})
	if err != nil {
		return
	}
	if err != nil {
		log.Fatalln("上传文件失败", err)
	}
	log.Println("上传文件成功", object)

	// download object
	err = client.FGetObject(context.Background(), bucketName, objectName, "down.txt", minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln("下载文件失败", err)
	}

	// list object
	objects := client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    "go/",
		Recursive: true,
	})
	for obj := range objects {
		log.Println("文件名", obj.Key, "大小", obj.Size, "修改时间", obj.LastModified, "元数据", obj.UserMetadata, "版本ID", obj.VersionID)
	}
	// delete object
	err = client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Fatalln("删除文件失败", err)
	}
	log.Println("删除文件成功")

	// delete bucket
	err = client.RemoveBucket(context.Background(), bucketName)
	if err != nil {
		log.Fatalln("删除bucket失败", err)
	}
	log.Println("删除bucket成功")
}
