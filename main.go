package main

import (
	"fmt"
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
		fmt.Println(err)
	}

	svc := sqs.New(sess)
	queueURL := "https://sqs.us-west-2.amazonaws.com/233704588990/barry-test-sqs.fifo"

	// M2Q(svc, queueURL)
	Q2M(svc, queueURL)

}

// M2Q is a method that message to aws SQS
func M2Q(svc *sqs.SQS, url string) {

	sendMessage := &sqs.SendMessageInput{
		DelaySeconds:           aws.Int64(0),
		MessageGroupId:         aws.String("the_first_group"),
		MessageDeduplicationId: aws.String("the_first_group"),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("The Whistler"),
			},
			"Author": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("John Grisham"),
			},
			"WeeksOn": &sqs.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String("6"),
			},
		},
		MessageBody: aws.String("message body"),
		QueueUrl:    aws.String(url),
	}

	result, err := svc.SendMessage(sendMessage)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Message send Succeeded", *result)
}

// Q2M is a method that receive message from aws SQS
func Q2M(svc *sqs.SQS, url string) {

	receiveMessage := &sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(url),
		MaxNumberOfMessages: aws.Int64(5),
		VisibilityTimeout:   aws.Int64(20),
		WaitTimeSeconds:     aws.Int64(0),
	}

	result, err := svc.ReceiveMessage(receiveMessage)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Message receive amount: ", len(result.Messages))

	for index, message := range result.Messages {
		fmt.Println(message)
		deleteMessage := &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(url),
			ReceiptHandle: message.ReceiptHandle,
		}
		result, err := svc.DeleteMessage(deleteMessage)
		if err != nil {
			fmt.Println(err, "Delete message error")
		} else {
			fmt.Println("Delete message success", index)
			fmt.Println(result.String())
		}
	}

}
