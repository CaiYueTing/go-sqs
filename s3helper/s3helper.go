package s3helper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3File struct {
	Bucket *string
	Key    *string
	file   *os.File
}

func NewS3File(bucket string, key string, filename string) (*S3File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &S3File{
		Bucket: &bucket,
		Key:    &key,
		file:   f,
	}, nil
}

func (sf *S3File) Upload2S3(sess *session.Session) {
	r, _ := ioutil.ReadAll(sf.file)
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		// Define a strategy that will buffer 25 MiB in memory
		u.BufferProvider = s3manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
	})
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: sf.Bucket,
		Key:    sf.Key,
		Body:   bytes.NewReader(r),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
