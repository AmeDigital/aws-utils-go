package s3utils

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"stash.b2w/asp/aws-utils-go.git/sessionutils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// GetObject downloads from an s3 bucket an object identified by its key and
// returns its content in raw format, as an array of bytes
func GetObject(bucketName string, key string) (data []byte, err error) {
	buf, err := getObjectAsBuf(bucketName, key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// GetObject downloads from an s3 bucket an object identified by its key and
// returns its content as a string
func GetObjectAsString(bucketName string, key string) (data string, err error) {
	buf, err := getObjectAsBuf(bucketName, key)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}

// ListObjects will retrieve the list of object keys that begins with the given keyPrefix.
func ListObjects(bucketName string, keyPrefix string) (keysList []*string, err error) {
	svc := s3.New(sessionutils.Session)
	res, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(keyPrefix),
	})

	if err != nil {
		fmt.Printf("Error listing bucket:\n%v\n", err)
		return keysList, err
	}

	if len(res.Contents) > 0 {
		for _, obj := range res.Contents {
			keysList = append(keysList, obj.Key)
		}
	}

	return keysList, nil

}

// PutObject uploads an object to a bucket and returns the url of the object created on s3.
func PutObject(bucketname string, key string, body string) (location string, err error) {
	uploader := s3manager.NewUploader(sessionutils.Session)
	reader := strings.NewReader(body)

	uploadOutput, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(key),
		Body:   reader,
	})

	if err != nil {
		return "", err
	}

	return uploadOutput.Location, nil
}

func getObjectAsBuf(bucketName string, key string) (data *bytes.Buffer, err error) {
	var S3Client *s3.S3 = s3.New(sessionutils.Session)

	getObjectOutput, err := S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, errors.New("ERROR: Could not retrieve object from bucket [" + bucketName + "] Key [" + key + "] Error: " + err.Error())
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(getObjectOutput.Body)

	return buf, err
}
