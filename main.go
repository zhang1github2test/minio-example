package main

import (
	"context"
	"github.com/minio/minio-go/v7"
	"minio-example/version"
)

func main() {
	client := version.InitMinioClient()
	client.RemoveBucket(context.Background(), "version-bucket")
	client.RemoveBucketWithOptions(context.Background(), "version-bucket", minio.RemoveBucketOptions{
		ForceDelete: true,
	})
	//version.CreateBucket(client, "version-bucket")
	//version.EnableVersion(client, "version-bucket")
	//
	//base.PutObject(client, "version-bucket", "test.txt", []byte("hello world,v1"))
	//base.PutObject(client, "version-bucket", "test.txt", []byte("hello world,v2"))
	//base.PutObject(client, "version-bucket", "test.txt", []byte("hello world,v3"))
	//
	//// 下载
	//bytes2, err := base.DownloadAsByte(client, "version-bucket", "test.txt")
	//if err != nil {
	//	log.Fatalln("minio下载文件失败", err)
	//}
	//fmt.Println(string(bytes2))
	//
	//// 下载指定版本
	//object, err := client.GetObject(context.Background(), "version-bucket", "test.txt", minio.GetObjectOptions{
	//	VersionID: "8d7f7447-bb09-4468-a730-dd45ae0dd57b",
	//})
	//// 读取对象为字节数组
	//buf := new(bytes.Buffer)
	//_, err = io.Copy(buf, object)
	//if err != nil {
	//	slog.Error("下载文件失败", err)
	//	return
	//}
	//data := buf.Bytes()
	//fmt.Println(string(data))
	//
	//// 列出对象的版本
	//objectPrefix := "test.txt"
	//
	//// 开始列出版本
	//ctx := context.Background()
	//opts := minio.ListObjectsOptions{
	//	Prefix:       objectPrefix,
	//	Recursive:    true,
	//	WithVersions: true, // 关键：启用版本列出
	//}
	//
	//for object := range client.ListObjects(ctx, "version-bucket", opts) {
	//	if object.Err != nil {
	//		log.Println("Error:", object.Err)
	//		continue
	//	}
	//	fmt.Printf("Object: %s, Size: %d, VersionID: %s, IsLatest: %v\n",
	//		object.Key, object.Size, object.VersionID, object.IsLatest)
	//}

}
