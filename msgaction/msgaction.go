package msgaction

import (
	"fmt"
	"gosqs/s3helper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type factoryInterface interface {
	GenerateMission(m string) missionInterface
}

type actionFactory struct {
}

func (a actionFactory) GenerateMission(m string, payload *string) missionInterface {
	switch m {
	case "shadow":
		return shadow{payload: payload}
	case "upgrade":
		return upgrade{payload: payload}
	case "reboot":
		return reboot{payload: payload}
	default:
		return nil
	}
}

type missionInterface interface {
	Do()
	GetPayload() string
}

type shadow struct {
	payload *string
}

func (s shadow) Do() {
	fmt.Println("this is shadow mission, upload file to s3")
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))
	uploadfile := s3helper.NewS3File(
		"barry-dlm-test",
		"shadow/upload.json",
		"filefolder/upload.json",
	)
	uploadfile.Upload2S3(sess)
	fmt.Println("success upload file")
}

func (s shadow) GetPayload() string {
	return *s.payload
}

type upgrade struct {
	payload *string
}

func (u upgrade) Do() {
	fmt.Println("this is upgrade mission, just print")
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))
	downloadfile := s3helper.NewS3File(
		"barry-dlm-test",
		"shadow/upload.json",
		"download/download.json",
	)

	err := downloadfile.Download(sess)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success download file")
}

func (u upgrade) GetPayload() string {
	return *u.payload
}

type reboot struct {
	payload *string
}

func (r reboot) Do() {
	fmt.Println("this is reboot mission, just print")
}

func (r reboot) GetPayload() string {
	return *r.payload
}

func NewMission(action string, payload *string) missionInterface {
	factory := new(actionFactory)
	return factory.GenerateMission(action, payload)
}
