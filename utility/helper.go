package utility

import (
	"encoding/json"
	"fmt"
	"gosqs/sqshelper"
	"io/ioutil"
	"os"
)

type Env struct {
	Queueurl string `json:"queueurl"`
}

var Envir Env
var Messages []sqshelper.Msg

func init() {
	readEnv()
	readMessage()
}

func readEnv() {
	readenv, err := os.Open("config.json")
	if err != nil {
		checkerr(err)
		return
	}
	defer readenv.Close()

	envjson, _ := ioutil.ReadAll(readenv)
	err = json.Unmarshal(envjson, &Envir)
	if err != nil {
		checkerr(err)
		return
	}
}

func readMessage() {
	readMessage, err := os.Open("message.json")
	if err != nil {
		checkerr(err)
		return
	}
	defer readMessage.Close()

	messagejson, _ := ioutil.ReadAll(readMessage)
	err = json.Unmarshal(messagejson, &Messages)
	if err != nil {
		checkerr(err)
		return
	}
}

func checkerr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
