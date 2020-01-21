package main

import (
	"fmt"
	"gosqs/msgaction"
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
		qmsg := Queue.NewQMsg(msg, url)
		err = qmsg.Send2Q(svc)
		if err != nil {
			fmt.Println(err)
		}
	}

	qmsgs := Queue.Receive(svc, url)

	for _, qmsg := range *qmsgs {
		fmt.Println("msg body", qmsg.Msg)
		actionFactory := new(msgaction.ActionFactory)
		mission := actionFactory.GenerateMission(qmsg.Msg.Action)
		mission.DoMission()
		qmsg.Delete(svc)
	}
}
