package main

import "minio-example/base"

func main() {
	client := base.InitMinioClient()
	base.DownloadAsFile(client, "bucket-lifecycle", "logs/bucket-lifecycle-file", "bucket-lifecycle-file.txt")

}
