package main

import (
	"fmt"
	"gosqs/sqshelper"
	"gosqs/utility"
	"strconv"

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

	url := utility.Envir.Queueurl

	for i, message := range utility.Messages {
		if message.Title == "" {
			continue
		}
		input := sqshelper.NewSendMessage(message, url, strconv.Itoa(i))
		err := sqshelper.Send(sess, input)
		if err != nil {
			fmt.Println(err, i)
		}
	}

	rec := sqshelper.NewReceiveMessage(url)
	msgs, err := sqshelper.Receive(sess, rec)

	recipes := []string{}

	for _, msg := range msgs {
		m := sqshelper.ToStruct(msg.MessageAttributes)
		fmt.Println(m.Title, m.Message, m.Action)
		recipes = append(recipes, *msg.ReceiptHandle)
	}

	for _, recipe := range recipes {
		deleteMsg := sqshelper.NewDeleteMessage(&recipe, url)
		err := sqshelper.Delete(sess, deleteMsg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
