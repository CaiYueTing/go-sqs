package utility

import (
	"encoding/json"
	"gosqs/sqshelper"
	"io/ioutil"
	"log"
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
	checkerr(err)
	defer readenv.Close()

	envjson, _ := ioutil.ReadAll(readenv)
	err = json.Unmarshal(envjson, &Envir)
	checkerr(err)
}

func readMessage() {
	readMessage, err := os.Open("message.json")
	checkerr(err)
	defer readMessage.Close()

	messagejson, _ := ioutil.ReadAll(readMessage)
	err = json.Unmarshal(messagejson, &Messages)
	checkerr(err)
}

func checkerr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
