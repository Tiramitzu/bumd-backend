package utils

import (
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOClient() (*minio.Client, error) {
	// Initialize minio client object.
	minioClient, err := MinioConnection()
	if err != nil {
		return minioClient, err
	}

	return minioClient, err
}

// MinioConnection func for opening minio connection.
func MinioConnection() (*minio.Client, error) {
	var err error
	var minioClient *minio.Client

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")

	// Initialize minio client object.
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return minioClient, err
	}

	return minioClient, err
}
