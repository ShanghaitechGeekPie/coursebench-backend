// Copyright (C) 2021-2024 ShanghaiTech GeekPie
// This file is part of CourseBench Backend.
//
// CourseBench Backend is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// CourseBench Backend is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with CourseBench Backend.  If not, see <http://www.gnu.org/licenses/>.

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
	Endpoint string `mapstructure:"endpoint"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Bucket   string `mapstructure:"bucket"`
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
func GetEndpoint() string {
	return MinioConf.Endpoint
}
