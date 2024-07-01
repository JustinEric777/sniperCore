package s3Cloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Cloud struct {
	Uploader   *s3manager.Uploader
	Downloader *s3manager.Downloader
	S3         *s3.S3
	Bucket     string
	Url        string
}

func NewS3Cloud(cloudConfig *CloudStorageConfig) (*S3Cloud, error) {
	credential := credentials.NewStaticCredentials(cloudConfig.Key, cloudConfig.Secret, cloudConfig.Token)

	config := aws.NewConfig().WithRegion(cloudConfig.Region).
		WithEndpoint(cloudConfig.EndPoint).
		WithCredentials(credential).WithS3ForcePathStyle(true)

	sess := session.New(config)

	return &S3Cloud{
		Uploader:   s3manager.NewUploader(sess),
		Downloader: s3manager.NewDownloader(sess),
		S3:         s3.New(sess),
		Bucket:     cloudConfig.Bucket,
		Url:        cloudConfig.Url,
	}, nil
}
