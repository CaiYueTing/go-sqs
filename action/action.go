package action

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
	uploadfile, err := s3helper.NewS3File(
		"barry-dlm-test",
		"shadow/upload.json",
		"filefolder/upload.json",
	)
	if err != nil {
		fmt.Println("new file error", err)
	}
	uploadfile.Upload2S3(sess)
}

type upgrade struct{}

func (u upgrade) DoMission() {
	fmt.Println("this is upgrade mission, just print")
}

type reboot struct{}

func (r reboot) DoMission() {
	fmt.Println("this is reboot mission, just print")
}
