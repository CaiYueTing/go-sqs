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

type ActionFactory struct {
}

func (a ActionFactory) GenerateMission(m string) missionInterface {
	switch m {
	case "shadow":
		return shadow{}
	case "upgrade":
		return upgrade{}
	case "reboot":
		return reboot{}
	default:
		return nil
	}
}

type missionInterface interface {
	DoMission()
}

type shadow struct{}

func (s shadow) DoMission() {
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

type upgrade struct{}

func (u upgrade) DoMission() {
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

type reboot struct{}

func (r reboot) DoMission() {
	fmt.Println("this is reboot mission, just print")
}
