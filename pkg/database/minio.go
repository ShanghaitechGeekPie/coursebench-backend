package database

import (
	"context"
	"coursebench-backend/pkg/errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"io"
	syslog "log"
)

type MinioConfig struct {
	Endpoint  string   `mapstructure:"endpoint"`
	Username  string   `mapstructure:"username"`
	Password  string   `mapstructure:"password"`
	Bucket    string   `mapstructure:"bucket"`
	Endpoint2 string   `mapstructure:"endpoint2"`
	IP        []string `mapstructure:"ip"`
}

var MinioConf MinioConfig
var minioClient *minio.Client

func InitMinio() {
	syslog.Println("Init Minio")
	config := viper.Sub("minio")
	if config == nil {
		syslog.Println("Minio config not found")
		return
	}
	err := config.Unmarshal(&MinioConf)
	if err != nil {
		syslog.Fatalln(err)
	}
	minioClient, err = minio.New(MinioConf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioConf.Username, MinioConf.Password, ""),
		Secure: true,
	})
	if err != nil {
		syslog.Fatalln(err)
	}
}

func UploadFile(ctx context.Context, objectName string, file io.Reader, fileSize int64) error {
	_, err := minioClient.PutObject(ctx, MinioConf.Bucket, objectName, file, fileSize, minio.PutObjectOptions{})
	if err != nil {
		return errors.Wrap(err, errors.MinIOError)
	}
	return nil
}

func DeleteFile(ctx context.Context, objectName string) error {
	err := minioClient.RemoveObject(ctx, MinioConf.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, errors.MinIOError)
	}
	return nil
}

// GetEndpoint returns the endpoint of the MinIO server.
// It gives on or off-campus nodes based on the user ip.
func GetEndpoint(ip []string) string {
	for _, i := range ip {
		for _, j := range MinioConf.IP {
			if i == j {
				return MinioConf.Endpoint
			}
		}
	}
	return MinioConf.Endpoint2
}
