package main

import (
	"fmt"
	"gosqs/action"
	Queue "gosqs/sqshelper"
	"gosqs/utility"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})
	if err != nil {
		log.Panic(err)
	}

	svc := sqs.New(sess)
	url := utility.Envir.Queueurl
	envMessages := utility.Messages

	for _, msg := range envMessages {
		if msg.Title == "" || msg.Action == "" || msg.Message == "" {
			continue
		}
		err = msg.Send2Q(svc, url)
		if err != nil {
			fmt.Println(err)
		}
	}

	msgs := Queue.Receive(svc, url)

	for _, msg := range msgs {
		fmt.Println("msg body", msg)
		missionfactory := new(action.ActionFactory)
		mission := missionfactory.GenerateMission(msg.Action)
		mission.DoMission()
		err = msg.Delete(svc, url)
		if err != nil {
			fmt.Println(err)
		}
	}
}
