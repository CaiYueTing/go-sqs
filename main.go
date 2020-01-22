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
	q := Queue.New("us-west-2", url)

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

	msgs := q.ReceiveMessage(3, 5)

	for _, msg := range *msgs {
		var m Msg
		json.Unmarshal([]byte(*msg.Msg), &m)
		fmt.Println("message body: ", m)
		actionFactory := new(msgaction.ActionFactory)
		mission := actionFactory.GenerateMission(m.Action)
		mission.DoMission()
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
