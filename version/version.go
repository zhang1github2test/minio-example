package version

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"io"
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
		slog.Info("检查bucket是否存在失败", "errMsg", errBucketExists)
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
		slog.Error("开启版本控制失败", "errMsg", err)
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

// DeleteObjectForce 强制删除对象
func DeleteObjectForce(client *minio.Client, bucketName, objectName string) {
	client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
}

// DeleteObjectWithVersion  删除指定版本的对象
func DeleteObjectWithVersion(client *minio.Client, bucketName, objectName, versionID string) {
	err := client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{
		VersionID: versionID,
	})
	if err != nil {
		slog.Error("物理删除对象失败：", "errMsg", err)
		return
	}
}

// GetLatestObject 获取最新版本对象
func GetLatestObject(client *minio.Client, bucketName, objectName string) {
	object, err := client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		slog.Error("获取对象失败:", err)
		return
	}
	defer object.Close()
	b := new(bytes.Buffer)
	io.Copy(b, object)
	slog.Info("对象内容：", "objectInfo", string(b.Bytes()))
}

// ListObjectVersion 获取一个对象的所有版本号
func ListObjectVersion(client *minio.Client, bucketName, objectName string) []string {
	var versionIds []string
	objects := client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:       objectName,
		WithVersions: true,
	})
	for object := range objects {
		slog.Info("objectName对象信息,", "versionId", object.VersionID, "IsDeleteMarker",
			object.IsDeleteMarker, "IsLatest", object.IsLatest)
		versionIds = append(versionIds, object.VersionID)
	}
	return versionIds
}

// GetObjectByVersion 指定版本获取对象
func GetObjectByVersion(client *minio.Client, bucketName, objectName, versionID string) {
	object, err := client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{
		VersionID: versionID,
	})
	if err != nil {
		slog.Error("根据版本获取对象失败:", "msg", err)
		return
	}
	defer object.Close()

	buf := make([]byte, 1024)
	n, _ := object.Read(buf)
	slog.Info("指定版本对象内容：", "version", versionID, "objectInfo", string(buf[:n]))
}

// DeleteBucketForce 强制删除开启版本控制的桶
func DeleteBucketForce(client *minio.Client, bucketName string) {
	client.RemoveBucketWithOptions(context.Background(), bucketName, minio.RemoveBucketOptions{
		ForceDelete: true,
	})
}
