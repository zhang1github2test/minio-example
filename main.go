package main

import (
	"log"
	"minio-example/base"
	"minio-example/lifecycle"
	"minio-example/version"
)

func main() {
	client := version.InitMinioClient()

	bucketName := "bucket-lifecycle"
	//base.PutObject(client, bucketName, "tmp/upload-123.tmp", []byte("hello,delete after 1 days"))
	//base.PutObject(client, bucketName, "logs/access-2024-06.log", []byte("hello,delete after 2 days"))
	state, err := base.GetObjectState(client, bucketName, "tmp/upload-123.tmp")
	if err != nil {
		log.Println("GetObjectState failed", err)
		return
	}
	log.Println("GetObjectState success", state)
	lifecycle.GetLifeCycle(client, bucketName)
}
