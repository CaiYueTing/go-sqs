package main

import (
	"fmt"
	"gosqs/sqshelper"
	"gosqs/utility"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})
	if err != nil {
		fmt.Println(err)
	}

	// for _, message := range utility.Messages {
	// 	if message.Title == "" {
	// 		continue
	// 	}
	// 	input := sqshelper.NewSendMessage(message, utility.Envir.Queueurl)
	// 	err := sqshelper.Send(sess, input)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	rec := sqshelper.NewReceiveMessage(utility.Envir.Queueurl)
	msgs, err := sqshelper.Receive(sess, rec)
	for _, msg := range msgs {
		m := sqshelper.ToStruct(msg.MessageAttributes)
		fmt.Println(m.Title, m.Message, m.Action)
	}
}
