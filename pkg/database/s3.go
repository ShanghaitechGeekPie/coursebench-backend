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
	"io"
	syslog "log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type S3Config struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	Secure    bool   `mapstructure:"secure"`
	Prefix    string `mapstructure:"prefix"`
}

var S3Conf S3Config
var s3Client *minio.Client
var s3Enabled bool

func InitS3() {
	syslog.Println("Init S3 (Knowledge Base)")
	config := viper.Sub("s3")
	if config == nil {
		syslog.Println("S3 config not found, knowledge base export disabled")
		s3Enabled = false
		return
	}
	config.SetDefault("secure", true)
	config.SetDefault("region", "")
	config.SetDefault("prefix", "kb")
	err := config.Unmarshal(&S3Conf)
	if err != nil {
		syslog.Fatalln(err)
	}
	s3Client, err = minio.New(S3Conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(S3Conf.AccessKey, S3Conf.SecretKey, ""),
		Secure: S3Conf.Secure,
		Region: S3Conf.Region,
	})
	if err != nil {
		syslog.Fatalln(err)
	}
	s3Enabled = true
	syslog.Println("S3 (Knowledge Base) initialized successfully")
}

func IsS3Enabled() bool {
	return s3Enabled
}

func S3Upload(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) error {
	if !s3Enabled {
		return errors.New(errors.S3Error)
	}
	_, err := s3Client.PutObject(ctx, S3Conf.Bucket, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return errors.Wrap(err, errors.S3Error)
	}
	return nil
}

func S3Delete(ctx context.Context, objectName string) error {
	if !s3Enabled {
		return errors.New(errors.S3Error)
	}
	err := s3Client.RemoveObject(ctx, S3Conf.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, errors.S3Error)
	}
	return nil
}

func S3List(ctx context.Context, prefix string) <-chan minio.ObjectInfo {
	return s3Client.ListObjects(ctx, S3Conf.Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
}
