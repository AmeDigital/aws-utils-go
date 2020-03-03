package sqsutils

import (
	"fmt"
	"reflect"
	"errors"
	"github.com/AmeDigital/aws-utils-go/sessionutils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func convertAWSDataType(value interface{}) (string, error) {
	typeof := reflect.TypeOf(value)

	switch typeof.Kind() {
		case reflect.String:
			return "String", nil
		case reflect.Array, reflect.Slice:
			return "String.Array", nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			return "Number", nil
		default:
			return "", errors.New("Type unknown")
	}
}

func GetMessageAttribute(queueUrl string, attributeName string) (string, error) {
	SQSclient := sqs.New(sessionutils.Session)
	
	var attributesNamesList []*string
	
	attributesNamesList = append(attributesNamesList, aws.String(attributeName))

	response, err := SQSclient.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		AttributeNames: attributesNamesList,
		QueueUrl:    &queueUrl,
	})

	if err != nil {
		return "", err
	}

	return *response.Attributes[attributeName], err
}

func SendMessage(queueUrl string, message string, messageAttributes map[string]interface{}) error {
	sendMessageInput := sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    &queueUrl,
	}

	if messageAttributes != nil {
		msgAttributeValueMap := make(map[string]*sqs.MessageAttributeValue, len(messageAttributes))

		for key, value := range messageAttributes {
			dataType, err := convertAWSDataType(value)

			if err != nil {
				return err
			}

			msgAttributeValueMap[key] = &sqs.MessageAttributeValue{
				DataType:    aws.String(dataType),
				StringValue: aws.String(value.(string)),
			}
		}

		sendMessageInput.MessageAttributes = msgAttributeValueMap
	}

	SQSclient := sqs.New(sessionutils.Session)

	_, err := SQSclient.SendMessage(&sendMessageInput)

	return err
}

func ReadMessage(queueUrl string, maxNumberOfMessages int64) ([]*sqs.Message, error) {

	SQSclient := sqs.New(sessionutils.Session)

	result, err := SQSclient.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: &maxNumberOfMessages,
		WaitTimeSeconds:     aws.Int64(0), //(Optional) Timeout in seconds for long polling
	})

	return result.Messages, err
}

func DeleteMessage(queueUrl string, receiptHandle string) error {
	SQSclient := sqs.New(sessionutils.Session)

	_, err := SQSclient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: &receiptHandle,
	})

	if err != nil {
		fmt.Println("Delete Error", err)
		return err
	}

	return nil
}