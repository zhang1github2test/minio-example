package version

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"log"
	"log/slog"
)

// initMinioClient 初始化minio客户端
func InitMinioClient() *minio.Client {
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

// CreateBucket 如果bucket不存在则创建
func CreateBucket(client *minio.Client, bucketName string) {
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

// EnableVersion 开启版本控制
func EnableVersion(client *minio.Client, bucketName string) {
	err := client.EnableVersioning(context.Background(), bucketName)
	if err != nil {
		slog.Info("开启版本控制失败", err)
		return
	}
	slog.Info("开启版本控制成功")
}

// SetLifeCycle 设置生命周期
func SetLifeCycle(client *minio.Client, bucketName string) {
	config := lifecycle.NewConfiguration()
	config.Rules = []lifecycle.Rule{
		{
			ID:     "expire-30day-lifecycle",
			Status: "Enabled",
			// 当前
			Expiration: lifecycle.Expiration{
				Days: 30,
			},
			NoncurrentVersionExpiration: lifecycle.NoncurrentVersionExpiration{
				NoncurrentDays: 1,
			},
		},
	}
	client.SetBucketLifecycle(context.Background(), bucketName, config)

}
