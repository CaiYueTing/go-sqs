package msgaction

import (
	"fmt"
	"gosqs/redishelper"
	"gosqs/s3helper"
)

type factoryInterface interface {
	GenerateMission(m string) missionInterface
}

type actionFactory struct {
}

func (a actionFactory) GenerateMission(m string, payload *string) missionInterface {
	switch m {
	case "shadow":
		return shadow{payload: payload}
	case "upgrade":
		return upgrade{payload: payload}
	case "reboot":
		return reboot{payload: payload}
	default:
		return nil
	}
}

type missionInterface interface {
	Do() (*string, error)
	GetPayload() string
}

type shadow struct {
	payload *string
}

func (s shadow) Do() (*string, error) {
	fmt.Println("this is shadow mission, upload file to s3")

	s3 := s3helper.NewS3("barry-dlm-test", "us-west-2")
	url, err := s3.Upload("filefolder/shadowupload.json", "shadow/upload.json")
	if err != nil {
		fmt.Println("upload failed", err)
		return nil, err
	}
	return url, nil
}

func (s shadow) GetPayload() string {
	return *s.payload
}

type upgrade struct {
	payload *string
}

func (u upgrade) Do() (*string, error) {
	fmt.Println("this is upgrade mission, download file from s3")
	s3 := s3helper.NewS3("barry-dlm-test", "us-west-2")

	err := s3.Download("download/download.json", "upload.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("success download file")
	return nil, nil
}

func (u upgrade) GetPayload() string {
	return *u.payload
}

type reboot struct {
	payload *string
}

func (r reboot) Do() (*string, error) {
	fmt.Println("this is reboot mission, use redis")
	redisdb, err := redishelper.NewRedisPool("127.0.0.1:6379")

	err = redisdb.SetString("key", "value")
	if err != nil {
		fmt.Println(err)
	}
	result, err := redisdb.ReadString("key")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*result)

	err = redisdb.PushList("lpush", "list", "aaa")
	if err != nil {
		fmt.Println(err)
	}

	results, err := redisdb.ReadList("list", 0, 10)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*results)

	ms := map[string]string{"a": "b", "c": "d"}
	err = redisdb.SetMap("hash", ms)
	if err != nil {
		fmt.Println(err)
	}
	resultmap, err := redisdb.ReadMap("hash")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*resultmap)

	reply := make(chan []byte)
	redisdb.Subscribe("topic", reply)
	msg := <-reply
	fmt.Println(string(msg))

	return nil, nil
}

func (r reboot) GetPayload() string {
	return *r.payload
}

func NewMission(action string, payload *string) missionInterface {
	factory := new(actionFactory)
	return factory.GenerateMission(action, payload)
}
