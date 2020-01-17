package main

import (
	"fmt"
	Queue "gosqs/sqshelper"
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
		input := Queue.NewSendMessage(message, url, strconv.Itoa(i))
		err := Queue.Send(sess, input)
		if err != nil {
			fmt.Println(err, i)
		}
	}

	rec := Queue.NewReceiveMessage(url)
	msgs, err := Queue.Receive(sess, rec)

	recipes := []string{}

	for _, msg := range msgs {
		m := Queue.ToStruct(msg.MessageAttributes)
		fmt.Println(m.Title, m.Message, m.Action)
		recipes = append(recipes, *msg.ReceiptHandle)
	}

	for _, recipe := range recipes {
		deleteMsg := Queue.NewDeleteMessage(&recipe, url)
		err := Queue.Delete(sess, deleteMsg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
