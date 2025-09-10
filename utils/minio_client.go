package utils

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConn struct {
	MinioClient *minio.Client
	BucketName  string
}

func NewMinIOConn(endpoint, accessKey, secretKey, bucketName string) (*MinioConn, error) {
	var err error
	conn := new(MinioConn)
	conn.BucketName = bucketName

	// Initialize minio client object.
	conn.MinioClient, err = minioConnection(endpoint, accessKey, secretKey)
	if err != nil {
		return conn, err
	}

	return conn, err
}

// MinioConnection func for opening minio connection.
func minioConnection(endpoint, accessKey, secretKey string) (*minio.Client, error) {
	var err error
	var minioClient *minio.Client

	// Initialize minio client object.
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return minioClient, err
	}

	return minioClient, err
}
