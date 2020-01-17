package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func init() {
	readenv, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer readenv.Close()

	envjson, _ := ioutil.ReadAll(readenv)
	err = json.Unmarshal(envjson, &env)
	if err != nil {
		fmt.Println(err)
	}
}

func TestM2Q(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})
	if err != nil {
		t.Error(err)
	}

	svc := sqs.New(sess)
	queueURL := env.Queueurl

	err = M2Q(svc, queueURL)
	if err != nil {
		t.Error(err)
	}
}

func TestQ2M(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})
	if err != nil {
		t.Error(err)
	}

	svc := sqs.New(sess)
	queueURL := env.Queueurl

	err = Q2M(svc, queueURL)
	if err != nil {
		t.Error(err)
	}
}
