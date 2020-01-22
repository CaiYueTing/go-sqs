package utility

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Env struct {
	Queueurl string `json:"queueurl"`
}

var Envir Env

func init() {
	readEnv()
}

func readEnv() {
	readenv, err := os.Open("config.json")
	checkerr(err)
	defer readenv.Close()

	envjson, _ := ioutil.ReadAll(readenv)
	err = json.Unmarshal(envjson, &Envir)
	checkerr(err)
}

func checkerr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
