package sqshelper

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Message from the message.json
type Msg struct {
	Title   string `json:"title"`
	Action  string `json:"action"`
	Message string `json:"message"`
}

func ToStruct(m map[string]*sqs.MessageAttributeValue) Msg {
	msg := Msg{}
	msg.Title = *m["title"].StringValue
	msg.Action = *m["action"].StringValue
	msg.Message = *m["message"].StringValue
	return msg
}

func NewSendMessage(m Msg, url string) *sqs.SendMessageInput {
	message := &sqs.SendMessageInput{
		DelaySeconds:           aws.Int64(0),
		MessageGroupId:         aws.String("the_group_id"),
		MessageDeduplicationId: aws.String("the_first_group"),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"title":   &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: &m.Title},
			"action":  &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: &m.Action},
			"message": &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: &m.Message},
		},
		MessageBody: aws.String("message body"),
		QueueUrl:    &url,
	}
	return message
}

func NewReceiveMessage(url string) *sqs.ReceiveMessageInput {
	message := &sqs.ReceiveMessageInput{
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
	return message
}

func Send(sess *session.Session, input *sqs.SendMessageInput) error {
	svc := sqs.New(sess)
	result, err := svc.SendMessage(input)
	if err != nil {
		return err
	}
	fmt.Println("send message success", result.MessageId)
	return nil
}

func Receive(sess *session.Session, input *sqs.ReceiveMessageInput) ([]*sqs.Message, error) {
	svc := sqs.New(sess)
	result, err := svc.ReceiveMessage(input)
	if err != nil {
		return nil, err
	}
	fmt.Println("Message receive amount: ", len(result.Messages))
	return result.Messages, nil
}
