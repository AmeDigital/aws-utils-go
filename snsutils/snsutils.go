package snsutils

import (
	"encoding/json"
	"errors"
	"fmt"

	"stash.b2w/asp/aws-utils-go.git/sessionutils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

// SendMessage sends an SNS message. The message instance must be a struct or a map[string]interface{},
// it will be marshalled into a json string and sent in the sns's "message" field.
func SendMessage(topicArn string, message interface{}) error {
	return SendMessageWithAttributes(topicArn, message, nil)
}

// SendMessage sends an SNS message acompained by SNS Message Attributes (used for subscription filtering).
// The message instance must be a struct or a map[string]interface{}, it will be marshalled
// into a json string and sent in the sns's "message" field.
func SendMessageWithAttributes(topicArn string, message interface{}, messageAttributes map[string]string) error {
	if len(topicArn) == 0 {
		return errors.New("topic arn cannot be empty")
	}

	snsService := sns.New(sessionutils.Session)

	messageJson, err := json.Marshal(message)

	if err != nil {
		fmt.Println(err)
		return err
	}

	params := &sns.PublishInput{
		Message:  aws.String(string(messageJson)),
		TopicArn: aws.String(topicArn),
	}

	if len(messageAttributes) > 0 {
		params.MessageAttributes = buildMessageAttributes(messageAttributes)
	}

	resp, err := snsService.Publish(params)

	if err != nil {
		fmt.Println(resp)
		fmt.Println(err)
		return err
	}

	paramsStr, _ := json.Marshal(params)
	fmt.Printf("Sent to topic '%s' the message: %s\n", topicArn, string(paramsStr))

	return nil
}

func buildMessageAttributes(metadataMap map[string]string) map[string]*sns.MessageAttributeValue {
	var messageAttributes = make(map[string]*sns.MessageAttributeValue)

	for k, v := range metadataMap {
		messageAttributeValue := sns.MessageAttributeValue{
			StringValue: aws.String(v),
			DataType:    aws.String("String"),
		}

		messageAttributes[k] = &messageAttributeValue
	}

	return messageAttributes
}
