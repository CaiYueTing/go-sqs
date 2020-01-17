package main

import (
	"encoding/json"
	"fmt"
	"gosqs/sqshelper"
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

var env Env
var messages []sqshelper.Msg

func init() {
	readenv, err := os.Open("config.json")
	if err != nil {
		checkerr(err)
		return
	}
	defer readenv.Close()

	envjson, _ := ioutil.ReadAll(readenv)
	err = json.Unmarshal(envjson, &env)
	if err != nil {
		checkerr(err)
		return
	}

	readM, err := os.Open("message.json")
	if err != nil {
		checkerr(err)
		return
	}
	defer readM.Close()

	messagejson, _ := ioutil.ReadAll(readM)
	err = json.Unmarshal(messagejson, &messages)
	if err != nil {
		checkerr(err)
		return
	}
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})
	if err != nil {
		checkerr(err)
	}

	// for _, message := range messages {
	// 	if message.Title == "" {
	// 		continue
	// 	}
	// 	input := sqshelper.NewSendMessage(message, env.Queueurl)
	// 	err := sqshelper.Send(sess, input)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	rec := sqshelper.NewReceiveMessage(env.Queueurl)
	msgs, err := sqshelper.Receive(sess, rec)
	for _, msg := range msgs {
		m := toStruct(msg.MessageAttributes)
		fmt.Println(m.Title, m.Message, m.Action)
	}
}

func toStruct(m map[string]*sqs.MessageAttributeValue) sqshelper.Msg {
	msg := sqshelper.Msg{}
	msg.Title = *m["title"].StringValue
	msg.Action = *m["action"].StringValue
	msg.Message = *m["message"].StringValue
	return msg
}

func checkerr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
