package main

import (
	"minio-example/base"
	"minio-example/version"
)

func main() {
	client := version.InitMinioClient()
	bucketName := "bucket-version"
	fileName := "bucket-version-file"
	// 创建桶
	version.CreateBucket(client, bucketName)

	// 启用版本管理
	version.EnableVersion(client, bucketName)

	// 上传文件
	base.PutObject(client, bucketName, fileName, []byte("hello, bucket-version ,v1"))

	// 再次上传
	base.PutObject(client, bucketName, fileName, []byte("hello, bucket-version ,v2"))

	// 获取最新的对象
	version.GetLatestObject(client, bucketName, fileName)

	// 获取一个对象所有的版本
	objectVersions := version.ListObjectVersion(client, bucketName, fileName)

	for _, versionId := range objectVersions {
		// 通过指定版本号来对象
		version.GetObjectByVersion(client, bucketName, fileName, versionId)
	}

	// 使用普通方法删除对象
	base.RemoveObject(client, bucketName, fileName)

	// 获取一个对象所有的版本
	objectVersions = version.ListObjectVersion(client, bucketName, fileName)
	for _, versionId := range objectVersions {
		// 通过指定版本号来对象
		version.GetObjectByVersion(client, bucketName, fileName, versionId)
		version.DeleteObjectWithVersion(client, bucketName, fileName, versionId)
	}

	// 彻底删除单个文件
	version.DeleteObjectForce(client, bucketName, fileName)

	// 删除桶
	version.DeleteBucketForce(client, bucketName)

}
