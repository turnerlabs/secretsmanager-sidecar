package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func main() {

	id := os.Getenv("SECRET_ID")
	if id == "" {
		fmt.Println("error: environment variable SECRET_ID is required")
		os.Exit(-1)
	}
	file := os.Getenv("SECRET_FILE")
	if file == "" {
		fmt.Println("error: environment variable SECRET_FILE is required")
		os.Exit(-1)
	}

	fmt.Println("fetching secret from secrets manager")
	secret, err := getSecret(id)
	if err != nil {
		panic(fmt.Errorf("error while fetching secret from secrets manager"))
	}

	//write secret to file
	fmt.Println("writing file")
	err = ioutil.WriteFile(file, []byte(*secret.SecretString), 0644)
	if err != nil {
		panic(fmt.Errorf("error while writing file %s: %v", file, err))
	}
	fmt.Println("done")
}

func getSecret(id string) (*secretsmanager.GetSecretValueOutput, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(id),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		panic(err)
	}

	return result, nil
}
