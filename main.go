package main

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"log/slog"
)

func main() {
	client := initMinioClient()
	createBucket(client, "test")
	uploadFile(client, "test", "test.pdf", []byte("hello world"))
}

// initMinioClient 初始化minio客户端
func initMinioClient() *minio.Client {
	endpoint := "121.43.141.218:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("初始minio客户端失败", err)
	}
	return client
}

// createBucket 如果bucket不存在则创建
func createBucket(client *minio.Client, bucketName string) {
	exists, errBucketExists := client.BucketExists(context.Background(), bucketName)
	if errBucketExists != nil {
		slog.Info("检查bucket是否存在失败", errBucketExists)
		return
	}
	if !exists {
		err := client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			slog.Info("创建桶失败！", err)
		} else {
			slog.Info("创建桶成功！")
		}
	}
}

func uploadFile(client *minio.Client, bucketName string, fileName string, fileBytes []byte) {
	ctx := context.Background()
	info, err := client.PutObject(ctx, bucketName, fileName, bytes.NewReader(fileBytes), int64(len(fileBytes)), minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"filename": "public-read.pdf"},
	})
	if err != nil {
		slog.Info("上传文件失败", err)
	} else {
		slog.Info("上传文件成功", info)
	}
}
