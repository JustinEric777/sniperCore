package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sniperCore/core/s3Cloud"
	"os"
	"path/filepath"
	"time"
)

type S3CloudService struct{}

func (obs *S3CloudService) GetObject(cloud, key string) (string, error) {
	if key == "" {
		return "", errors.New("key param is not empty")
	}

	s3CloudObj := s3Cloud.GetS3Cloud(cloud)
	inputParams := &s3.GetObjectInput{
		Bucket: aws.String(s3CloudObj.Bucket),
		Key:    aws.String(key),
	}

	resp, err := s3CloudObj.S3.GetObject(inputParams)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (obs *S3CloudService) GetObjectToLocal(cloud, key, fileName string) (string, error) {
	if key == "" {
		return "", errors.New("key param is not empty")
	}

	dir := filepath.Dir(fileName)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return "", err
		}
	}

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	s3CloudObj := s3Cloud.GetS3Cloud(cloud)
	inputParams := &s3.GetObjectInput{
		Bucket: aws.String(s3CloudObj.Bucket),
		Key:    aws.String(key),
	}

	_, err = s3CloudObj.Downloader.Download(file, inputParams)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (obs *S3CloudService) GetObjectUrl(cloud, key string, expireTime int) (string, error) {
	if key == "" {
		return "", errors.New("key param is not empty")
	}
	if expireTime == 0 {
		expireTime = 300
	}

	s3CloudObj := s3Cloud.GetS3Cloud(cloud)
	inputParams := &s3.GetObjectInput{
		Bucket: aws.String(s3CloudObj.Bucket),
		Key:    aws.String(key),
	}

	resp, _ := s3CloudObj.S3.GetObjectRequest(inputParams)
	url, _ := resp.Presign(time.Duration(expireTime) * time.Second)

	return url, nil
}

func (obs *S3CloudService) PutObject(cloud, key, contentType, body string, metadata map[string]string) (string, error) {
	if key == "" {
		return "", errors.New("key param is not empty")
	}
	if contentType == "" {
		return "", errors.New("contentType param is not empty")
	}
	if body == "" {
		return "", errors.New("body param is not empty")
	}

	//meta信息处理
	var metadataInfo map[string]*string
	for key, val := range metadata {
		metadataInfo[key] = aws.String(val)
	}

	s3CloudObj := s3Cloud.GetS3Cloud(cloud)
	inputParams := &s3.PutObjectInput{
		Bucket:      aws.String(s3CloudObj.Bucket),
		Key:         aws.String(key),
		ACL:         aws.String("private"),
		ContentType: aws.String(contentType),
		Body:        bytes.NewReader([]byte(body)),
		Metadata:    metadataInfo,
	}

	_, err := s3CloudObj.S3.PutObject(inputParams)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", s3CloudObj.Url, key), nil
}

func (obs *S3CloudService) PutObjectFromLocal(cloud, key, fileName string) (string, error) {
	if fileName == "" {
		return "", errors.New("fileName param is not empty")
	}
	if key == "" {
		return "", errors.New("key param is not empty")
	}

	fileContent, err := os.Open(fileName)
	if err != nil {
		return "", err
	}

	s3CloudObj := s3Cloud.GetS3Cloud(cloud)
	inputParams := &s3manager.UploadInput{
		Bucket: aws.String(s3CloudObj.Bucket),
		Key:    aws.String(key),
		Body:   fileContent,
	}

	_, err = s3CloudObj.Uploader.Upload(inputParams, func(uploader *s3manager.Uploader) {
		uploader.PartSize = 50 * 1024 * 1024
		uploader.Concurrency = 10
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", s3CloudObj.Url, key), nil
}

func (obs *S3CloudService) DeleteObject(cloud, key string) error {
	if key == "" {
		return errors.New("key param is not empty")
	}

	s3CloudObj := s3Cloud.GetS3Cloud(cloud)
	inputParams := &s3.DeleteObjectInput{
		Bucket: aws.String(s3CloudObj.Bucket),
		Key:    aws.String(key),
	}

	_, err := s3CloudObj.S3.DeleteObject(inputParams)
	if err != nil {
		return err
	}
	return nil
}
