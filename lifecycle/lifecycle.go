package lifecycle

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"log"
)

// SetLifeCycle 设置生命周期
func SetLifeCycle(client *minio.Client, bucketName string) {
	config := lifecycle.NewConfiguration()
	config.Rules = []lifecycle.Rule{
		{
			ID:     "expire-2day-lifecycle",
			Status: "Enabled",
			// 超过两天进行删除
			Expiration: lifecycle.Expiration{
				Days: 2,
			},
			// 超过一天进行转储
			Transition: lifecycle.Transition{
				Days:         1,
				StorageClass: "WARM-MINIO-TIER",
			},
			RuleFilter: lifecycle.Filter{
				Prefix: "logs/",
			},
		},
		{
			ID:     "expire-1day-lifecycle",
			Status: "Enabled",
			// 超过1天进行删除
			Expiration: lifecycle.Expiration{
				Days: 1,
			},
			RuleFilter: lifecycle.Filter{
				Prefix: "tmp/",
			},
		},
	}
	err := client.SetBucketLifecycle(context.Background(), bucketName, config)
	if err != nil {
		log.Println("SetBucketLifecycle failed", err)
		return
	}
}

// GetLifeCycle 获取生命周期规则
func GetLifeCycle(client *minio.Client, bucketName string) {
	config, err := client.GetBucketLifecycle(context.Background(), bucketName)
	if err != nil {
		fmt.Println(err)
	}
	marshal, err := xml.MarshalIndent(config, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(marshal))
}
