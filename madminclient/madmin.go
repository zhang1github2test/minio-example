package madminclient

import (
	"fmt"
	"github.com/minio/madmin-go/v4"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitAdminClient() *madmin.AdminClient {
	// Use a secure connection.
	ssl := true

	// Initialize minio client object.
	mdmClnt, err := madmin.NewWithOptions("121.43.141.218:9000", &madmin.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: ssl,
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return mdmClnt
}
