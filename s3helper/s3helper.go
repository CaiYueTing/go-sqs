package s3helper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3File struct {
	Bucket *string
	Key    *string
	file   *string
}

func NewS3File(bucket string, key string, filename string) *S3File {
	return &S3File{
		Bucket: &bucket,
		Key:    &key,
		file:   &filename,
	}
}

func (sf *S3File) Upload2S3(sess *session.Session) error {
	f, err := os.Open(*sf.file)
	if err != nil {
		fmt.Println("file not found:", err)
		return err
	}
	r, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("read file error:", err)
		return err
	}
	defer f.Close()
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.BufferProvider = s3manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
	})
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: sf.Bucket,
		Key:    sf.Key,
		Body:   bytes.NewReader(r),
	})
	if err != nil {
		fmt.Println("upload file failed:", err)
		return err
	}
	return nil
}

func (sf *S3File) Download(sess *session.Session) error {
	downloader := s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
		d.PartSize = 64 * 1024 * 1024
	})
	file, err := os.Create(*sf.file)
	if err != nil {
		fmt.Println(err, "Create file failed")
		return err
	}
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: sf.Bucket,
			Key:    sf.Key,
		})
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	return nil
}
