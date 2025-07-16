package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"log/slog"
	"time"
)

func main() {
	client := initMinioClient()
	createBucket(client, "test")
	PutObject(client, "test", "test2.txt", []byte("hello world 2"))
	FPutObject(client, "test", "uploadfile.txt")
	DownloadAsFile(client, "test", "test2.txt", "test.txt")
	GetObjectState(client, "test", "test2.txt")
	PresignedGetObject(client, "test", "uploadfile.txt")
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

// PutObject 上传文件到minio,这里直接使用byte[]
func PutObject(client *minio.Client, bucketName string, fileName string, fileBytes []byte) {
	ctx := context.Background()
	info, err := client.PutObject(ctx, bucketName, fileName, bytes.NewReader(fileBytes), int64(len(fileBytes)), minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"upload-type": "byte",
		},
		ContentType:        "text/plain",
		ContentDisposition: "attachment; filename=\"" + fileName + "_new\"",
	})
	if err != nil {
		slog.Info("上传文件失败", err)
	} else {
		slog.Info("上传文件成功", info)
	}
}

// FPutObject 上传本地文件到minio
func FPutObject(client *minio.Client, bucketName string, fileName string) {
	ctx := context.Background()
	info, err := client.FPutObject(ctx, bucketName, fileName, "uploadfile.txt", minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"upload-type": "file"},
	})
	if err != nil {
		slog.Info("上传文件失败", err)
	} else {
		slog.Info("上传文件成功", info)
	}
}

// DownloadAsFile 下载文件到本地
func DownloadAsFile(client *minio.Client, bucketName string, fileName string, filePath string) error {
	ctx := context.Background()
	err := client.FGetObject(ctx, bucketName, fileName, filePath, minio.GetObjectOptions{})
	if err != nil {
		slog.Info("下载文件失败", err)
		return err
	}
	return nil
}

// DownloadAsByte 下载文件,直接以[]byte返回
func DownloadAsByte(client *minio.Client, bucketName string, fileName string, filePath string) ([]byte, error) {
	ctx := context.Background()
	object, err := client.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		slog.Info("下载文件失败", err)
		return nil, err
	}
	// 读取对象为字节数组
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, object)
	if err != nil {
		slog.Error("下载文件失败", err)
		return nil, err
	}
	data := buf.Bytes()
	return data, nil
}

// GetObjectState 获取对象状态信息
func GetObjectState(client *minio.Client, bucketName string, fileName string) (*minio.ObjectInfo, error) {
	ctx := context.Background()
	statobject, err := client.StatObject(ctx, bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		slog.Error("获取对象状态信息失败", err)
	}
	// 打印常用系统元数据
	fmt.Println("ETag:", statobject.ETag)
	fmt.Println("Size:", statobject.Size)
	fmt.Println("ContentType:", statobject.ContentType)
	fmt.Println("LastModified:", statobject.LastModified)

	// 打印用户自定义元数据（X-Amz-Meta-*）
	fmt.Println("User Metadata:")
	for k, v := range statobject.UserMetadata {
		fmt.Printf("  %s: %s\n", k, v)
	}

	// 打印完整响应头信息（包含 Content-Disposition、Cache-Control 等）
	fmt.Println("All Headers:")
	for k, v := range statobject.Metadata {
		fmt.Printf("  %s: %v\n", k, v)
	}
	return &statobject, nil
}

// PresignedGetObject 获取文件外链
func PresignedGetObject(client *minio.Client, bucketName string, fileName string) {
	ctx := context.Background()
	u, err := client.PresignedGetObject(ctx, bucketName, fileName, time.Second*1000, nil)
	if err != nil {
		slog.Info("获取外链信息失败", "桶名", bucketName, "文件名", fileName, "错误信息", err)
		return
	}
	slog.Info("获取到的外链信息为", "url", u.String())
	presignedPutObject, err := client.PresignedPutObject(ctx, bucketName, "text_upload_by_presigned", time.Second*1000)
	if err != nil {
		slog.Info("获取上传外链信息失败", "桶名", bucketName, "文件名", fileName, "错误信息", err)
		return
	}
	slog.Info("获取到的上传外链信息为", "url", presignedPutObject.String())
}
