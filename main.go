package main

import (
	"encoding/json"
	"fmt"
	"gosqs/msgaction"
	Queue "gosqs/sqshelper"
	"gosqs/utility"
	"io/ioutil"
	"log"
	"os"
)

// Msg is User defined structure
type Msg struct {
	Title   string `json:"title"`
	Action  string `json:"action"`
	Message string `json:"message"`
}

// initialize the message from message.json
func init() {
	readMessage()
}

var messages []Msg

func main() {
	url := utility.Envir.Queueurl
	region := "us-west-2"
	q := Queue.New(region, url)

	for _, msg := range messages {
		if msg.Title == "" || msg.Action == "" || msg.Message == "" {
			continue
		}
		bmsg, _ := json.Marshal(msg)
		err := q.SendMessage(string(bmsg))
		if err != nil {
			fmt.Println(err)
		}
	}

	msgs := q.ReceiveMessage(3, 10)

	for _, msg := range *msgs {
		var m Msg
		json.Unmarshal([]byte(*msg.Msg), &m)
		mission := msgaction.NewMission(m.Action, msg.Msg)
		url, err := mission.Do()
		if err != nil {
			fmt.Println("mission failed:", err)
		}
		if url != nil {
			fmt.Println("s3 file url:", *url)
		}
		q.Delete(msg.ReceiptHandle)
	}
}

func readMessage() {
	readMessage, err := os.Open("message.json")
	checkerr(err)
	defer readMessage.Close()

	messagejson, _ := ioutil.ReadAll(readMessage)
	err = json.Unmarshal(messagejson, &messages)
	checkerr(err)
}

func checkerr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
