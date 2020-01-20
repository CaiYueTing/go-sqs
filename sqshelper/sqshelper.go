package sqshelper

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	uuid "github.com/satori/go.uuid"
)

// Msg from the message.json
type Msg struct {
	Title         string `json:"title"`
	Action        string `json:"action"`
	Message       string `json:"message"`
	ReceiptHandle *string
}

func (msg *Msg) newSendMessage(url string) *sqs.SendMessageInput {
	body, err := json.Marshal(msg)
	if err != nil {
		log.Panic(err)
	}
	defer recoverfunc(err, "new send message")
	return &sqs.SendMessageInput{
		DelaySeconds:           aws.Int64(0),
		MessageGroupId:         aws.String("GroupId"),
		MessageDeduplicationId: aws.String(uuid.NewV4().String()),
		MessageBody:            aws.String(string(body)), // Unmarshal
		QueueUrl:               &url,
	}
}

func newReceiveMessage(url string) *sqs.ReceiveMessageInput {
	return &sqs.ReceiveMessageInput{
		QueueUrl:            &url,
		MaxNumberOfMessages: aws.Int64(5),
		VisibilityTimeout:   aws.Int64(20),
		WaitTimeSeconds:     aws.Int64(0),
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
	}
}

func (msg *Msg) newDeleteMessage(url string) *sqs.DeleteMessageInput {
	return &sqs.DeleteMessageInput{
		QueueUrl:      &url,
		ReceiptHandle: msg.ReceiptHandle,
	}
}

func (msg *Msg) Send2Q(svc *sqs.SQS, url string) error {
	// gid is duplicate group id
	result, err := svc.SendMessage(msg.newSendMessage(url))

	if err != nil {
		log.Panic(err)
	}
	defer recoverfunc(err, "send message to queue")

	fmt.Println("send message success", result.MessageId)
	return nil
}

func Receive(svc *sqs.SQS, url string) []Msg {
	result, err := svc.ReceiveMessage(newReceiveMessage(url))

	if err != nil {
		log.Panic(err)
	}
	defer recoverfunc(err, "receive message")

	fmt.Println("Message receive amount: ", len(result.Messages))
	return msg2Struct(result.Messages)
}

func msg2Struct(msgs []*sqs.Message) []Msg {
	messages := []Msg{}
	for _, msg := range msgs {
		var m Msg
		json.Unmarshal([]byte(*msg.Body), &m)
		m.ReceiptHandle = msg.ReceiptHandle
		messages = append(messages, m)
	}
	return messages
}

func (msg *Msg) Delete(svc *sqs.SQS, url string) error {
	_, err := svc.DeleteMessage(msg.newDeleteMessage(url))
	if err != nil {
		return err
	}
	fmt.Println("Message delete")
	return nil
}

func recoverfunc(err error, method string) {
	if err := recover(); err != nil {
		fmt.Println("recover"+method, err)
	}
}
