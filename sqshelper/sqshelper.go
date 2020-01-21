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

type QueueMsg struct {
	url *string
	Msg *Msg
}

func NewQMsg(msg Msg, url string) *QueueMsg {
	return &QueueMsg{
		url: &url,
		Msg: &msg,
	}
}

func (qmsg *QueueMsg) newSendMessage() *sqs.SendMessageInput {
	body, err := json.Marshal(qmsg.Msg)
	if err != nil {
		log.Panic(err)
	}
	defer recoverfunc(err, "new send message")
	return &sqs.SendMessageInput{
		DelaySeconds:           aws.Int64(0),
		MessageGroupId:         aws.String("GroupId"),
		MessageDeduplicationId: aws.String(uuid.NewV4().String()),
		MessageBody:            aws.String(string(body)), // Unmarshal
		QueueUrl:               qmsg.url,
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

func (qmsg *QueueMsg) newDeleteMessage() *sqs.DeleteMessageInput {
	return &sqs.DeleteMessageInput{
		QueueUrl:      qmsg.url,
		ReceiptHandle: qmsg.Msg.ReceiptHandle,
	}
}

func (qmsg *QueueMsg) Send2Q(svc *sqs.SQS) error {
	input := qmsg.newSendMessage()
	result, err := svc.SendMessage(input)

	if err != nil {
		log.Panic(err)
	}
	defer recoverfunc(err, "send message to queue")

	fmt.Println("send message success", result.MessageId)
	return nil
}

func Receive(svc *sqs.SQS, url string) *[]QueueMsg {
	result, err := svc.ReceiveMessage(newReceiveMessage(url))

	if err != nil {
		log.Panic(err)
	}
	defer recoverfunc(err, "receive message")

	fmt.Println("Message receive amount: ", len(result.Messages))
	return msg2Struct(result.Messages, url)
}

func msg2Struct(msgs []*sqs.Message, url string) *[]QueueMsg {
	messages := []QueueMsg{}
	for _, msg := range msgs {
		var qm QueueMsg
		json.Unmarshal([]byte(*msg.Body), &qm.Msg)
		qm.Msg.ReceiptHandle = msg.ReceiptHandle
		qm.url = &url
		messages = append(messages, qm)
	}
	return &messages
}

func (qmsg *QueueMsg) Delete(svc *sqs.SQS) {
	input := qmsg.newDeleteMessage()
	_, err := svc.DeleteMessage(input)

	if err != nil {
		log.Panic(err)
	}
	defer recoverfunc(err, "delete message")

	fmt.Println("Message delete")
}

func recoverfunc(err error, method string) {
	if err := recover(); err != nil {
		fmt.Println("recover"+method, err)
	}
}
