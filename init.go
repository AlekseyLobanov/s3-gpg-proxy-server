package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const ENCRYPT_KEY_PATH = "_key.gpg"

func getMinio(config *S3Config) (*minio.Client, error) {
	creds := credentials.NewStaticV4(config.AccessKey, config.SecretKey, "")
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  creds,
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}

func getInternalMinio(config *AppConfig) (*minio.Client, error) {
	return getMinio(&config.LocalS3)
}

func getTargetMinio(config *AppConfig) (*minio.Client, error) {
	return getMinio(&config.TargetS3)
}

func prepareKey(config *AppConfig, minioClient *minio.Client) error {
	var err error
	if _, err := os.Stat(config.Key.LocalPath); !os.IsNotExist(err) {
		return nil
	}
	err = minioClient.FGetObject(
		context.Background(),
		config.Key.Bucket,
		config.Key.Name,
		ENCRYPT_KEY_PATH,
		minio.GetObjectOptions{},
	)
	if err == nil {
		config.Key.LocalPath = ENCRYPT_KEY_PATH
	}
	return err
}

func init_app() (*minio.Client, *minio.Client, *AppConfig) {
	fmt.Println("initing...")
	var err error

	config := GetConfig()

	var internalMinio *minio.Client
	internalMinio, err = getInternalMinio(config)
	if err != nil {
		log.Fatalln(err)
	}

	var targetMinio *minio.Client
	targetMinio, err = getTargetMinio(config)
	if err != nil {
		log.Fatalln(err)
	}

	err = prepareKey(config, targetMinio)
	if err != nil {
		fmt.Println("Unable to get public key from target S3")
		log.Fatalln(err)
	}
	return internalMinio, targetMinio, config
}
