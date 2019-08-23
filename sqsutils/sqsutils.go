package sqsutils

import (
	"aws-utils-go/sessionutils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func SendMessage(queueUrl string, message string) error {

	SQSclient := sqs.New(sessionutils.Session)

	_, err := SQSclient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    &queueUrl,
	})

	return err
}
