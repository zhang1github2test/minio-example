# 💾 MinIO Go 客户端使用详解：对象存储开发实战指南

随着云原生架构的发展，**对象存储**已成为现代数据存储的主流方式。**MinIO** 作为一款高性能、兼容 S3 协议的对象存储服务，因其部署灵活、开源透明等特性，在私有云和本地部署场景中得到了广泛应用。

本文将详细介绍如何使用 **MinIO 的 Go 语言客户端（minio-go）**，实现对象的上传、下载、浏览与删除操作。内容覆盖实际开发常用操作，适合希望通过 Go 操作对象存储的工程师。

---

## ✅ 一、准备工作

### 1. 环境依赖

- Go 版本：建议 Go 1.16+
- MinIO 已部署并运行（本地或远程皆可）
- 获取 AccessKey 和 SecretKey

### 2. 安装 SDK

```bash
go get github.com/minio/minio-go/v7
go get github.com/minio/minio-go/v7/pkg/credentials
```

---

## 🔧 二、初始化 MinIO 客户端

```go
package main

import (
    "context"
    "log"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
    endpoint := "localhost:9000"
    accessKeyID := "minioadmin"
    secretAccessKey := "minioadmin"
    useSSL := false

    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
        Secure: useSSL,
    })
    if err != nil {
        log.Fatalln("初始化失败:", err)
    }

    log.Println("✅ 成功连接 MinIO")
}
```

---

## 📦 三、创建 Bucket（存储桶）

```go
func createBucket(client *minio.Client, bucketName string) {
    ctx := context.Background()

    exists, err := client.BucketExists(ctx, bucketName)
    if err != nil {
        log.Fatal(err)
    }

    if !exists {
        err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
        if err != nil {
            log.Fatal(err)
        }
        log.Println("✅ Bucket 创建成功:", bucketName)
    } else {
        log.Println("ℹ️ Bucket 已存在:", bucketName)
    }
}
```

---

## ⬆️ 四、上传对象

```go
func uploadObject(client *minio.Client, bucketName, filePath, objectName string) {
    ctx := context.Background()

    info, err := client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("✅ 上传成功: %s (%d bytes)\n", info.Key, info.Size)
}
```

---

## ⬇️ 五、下载对象

```go
func downloadObject(client *minio.Client, bucketName, objectName, downloadPath string) {
    ctx := context.Background()

    err := client.FGetObject(ctx, bucketName, objectName, downloadPath, minio.GetObjectOptions{})
    if err != nil {
        log.Fatal(err)
    }

    log.Println("✅ 下载成功:", downloadPath)
}
```

---

## 📂 六、列出对象列表

```go
func listObjects(client *minio.Client, bucketName string) {
    ctx := context.Background()

    objectCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
        Recursive: true,
    })

    for obj := range objectCh {
        if obj.Err != nil {
            log.Println("❌ 错误:", obj.Err)
            continue
        }
        log.Printf("📄 文件: %s | 大小: %d 字节\n", obj.Key, obj.Size)
    }
}
```

---

## 🗑️ 七、删除对象

```go
func deleteObject(client *minio.Client, bucketName, objectName string) {
    ctx := context.Background()

    err := client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
    if err != nil {
        log.Fatal(err)
    }

    log.Println("🗑️ 删除成功:", objectName)
}
```

---

## 🧩 八、完整 main 函数示例

```go
func main() {
    client := initMinIOClient()

    bucket := "mybucket"
    object := "example.txt"
    file := "./example.txt"
    download := "./downloaded.txt"

    createBucket(client, bucket)
    uploadObject(client, bucket, file, object)
    listObjects(client, bucket)
    downloadObject(client, bucket, object, download)
    deleteObject(client, bucket, object)
}
```

## 🔚 九、总结

MinIO Go SDK 提供了简洁高效的 API，适合各类后台服务接入对象存储系统。通过本文的学习，你应该能够完成：

* MinIO 客户端初始化；
* Bucket 创建与检查；
* 对象上传、下载、列出、删除；
* 常见问题排查。

如果你在开发日志存储、图片归档、数据备份等系统中需要接入 S3 存储接口，MinIO 是非常优秀且易用的选择。

---

### 📌 推荐阅读：

* [MinIO 官方文档](https://docs.min.io/)
* [MinIO Go SDK GitHub](https://github.com/minio/minio-go)

---

如果觉得本文有帮助，欢迎点赞👍、评论💬、收藏⭐，持续更新更多 **对象存储 + 云原生 + Go 编程** 实战教程！

