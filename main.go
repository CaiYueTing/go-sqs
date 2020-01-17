package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Env is the environment variable
type Env struct {
	Queueurl string `json:"queueurl"`
}

// Message from the message.json
type Message struct {
	Title   string `json:"title"`
	Action  string `json:"action"`
	Message string `json:"message"`
}

var env Env
var messages []Message

func init() {
	readenv, err := os.Open("config.json")
	checkerr(err)
	defer readenv.Close()

	envjson, _ := ioutil.ReadAll(readenv)
	err = json.Unmarshal(envjson, &env)
	checkerr(err)

	readM, err := os.Open("message.json")
	checkerr(err)
	defer readM.Close()

	messagejson, _ := ioutil.ReadAll(readM)
	err = json.Unmarshal(messagejson, &messages)
	checkerr(err)
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})
	if err != nil {
		fmt.Println(err)
	}

	svc := sqs.New(sess)
	queueURL := env.Queueurl

	M2Q(svc, queueURL)
	Q2M(svc, queueURL)

}

func newMessage() {

}

// M2Q is a method that message to aws SQS
func M2Q(svc *sqs.SQS, url string) error {

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
		return err
	}

	fmt.Println("Message send Succeeded", *result)
	return nil
}

// Q2M is a method that receive message from aws SQS
func Q2M(svc *sqs.SQS, url string) error {

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
		return err
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
	return nil
}

func checkerr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
