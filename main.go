package main

import (
	"context"
	"fmt"

	"os/exec"

	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
)

func processEncryption(
	config *AppConfig,
	internalMinio *minio.Client,
	targetMinio *minio.Client,
	bucket string,
	name string,
) error {
	var err error

	found, err := targetMinio.BucketExists(context.Background(), bucket)
	if err != nil {
		return err
	}
	if !found {
		fmt.Printf("Bucket %v is not exists on remote server, creating...\n", bucket)
		err = targetMinio.MakeBucket(
			context.Background(),
			bucket,
			minio.MakeBucketOptions{Region: config.TargetS3.Region},
		)
		if err != nil {
			return err
		}
		fmt.Printf("Bucket %v has been added on remote server\n", bucket)
	}

	var objectReader *minio.Object
	objectReader, err = internalMinio.GetObject(
		context.Background(),
		bucket,
		name,
		minio.GetObjectOptions{},
	)

	if err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command("gpg", "--encrypt", "--recipient-file", ENCRYPT_KEY_PATH)
	cmd.Stdin = objectReader
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	if err = cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	uploadInfo, err_put := targetMinio.PutObject(
		context.Background(),
		bucket,
		name+".gpg",
		stdout,
		-1,
		minio.PutObjectOptions{
			ContentType:  "application/octet-stream",
			UserMetadata: map[string]string{"gpg": "true"},
		},
	)

	if err = cmd.Wait(); err != nil {
		return err
	}

	if err_put != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Successfully uploaded bytes: ", uploadInfo.Size)

	err = internalMinio.RemoveObject(
		context.Background(),
		bucket,
		name,
		minio.RemoveObjectOptions{},
	)
	return err
}

func processNotification(
	config *AppConfig,
	internalMinio *minio.Client,
	targetMinio *minio.Client,
	event notification.Event,
) error {
	var err error
	fmt.Printf(
		"bucket: %v name: %v size: %v user_meta: %v\n",
		event.S3.Bucket.Name,
		event.S3.Object.Key,
		event.S3.Object.Size,
		event.S3.Object.UserMetadata,
	)
	_, gpg_exists := event.S3.Object.UserMetadata["X-Amz-Meta-Gpg"]
	if gpg_exists {
		fmt.Println("Already encrypted")
		return nil
	}
	err = processEncryption(
		config,
		internalMinio,
		targetMinio,
		event.S3.Bucket.Name,
		event.S3.Object.Key,
	)
	return err
}

func main() {
	var err error
	internalMinio, targetMinio, config := init_app()

	for notificationInfo := range internalMinio.ListenNotification(context.Background(), "", "", []string{
		"s3:ObjectCreated:*",
	}) {
		if notificationInfo.Err != nil {
			log.Fatalln(notificationInfo.Err)
		}
		err = processNotification(config, internalMinio, targetMinio, notificationInfo.Records[0])
		if err != nil {
			log.Fatalln(err)
		}
	}
}
