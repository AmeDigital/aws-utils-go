package dynamodbutils

import (
	"errors"

	"stash.b2w/asp/aws-utils-go.git/sessionutils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Key holds the name and value of a partition key (PK) and an optional sort key (SK)
type Key struct {
	PKName  string
	PKValue interface{}
	SKName  string      // optional
	SKValue interface{} // optional
}

// UpdateItem updates the fields of an item identified by its partitionKey and sortKey(optional).
//
// Arguments:
//
// tablename: the name of the table
//
// key: the item's partition key and optional sort key
//
// fields: a map of field name/value pairs that will be updated
func UpdateItem(tablename string, key Key, fields map[string]interface{}) (err error) {

	svc := dynamodb.New(sessionutils.Session)

	keyAttributes := make(map[string]*dynamodb.AttributeValue)

	keyAttributes[key.PKName], err = dynamodbattribute.Marshal(key.PKValue)
	if err != nil {
		return err
	}

	if len(key.SKName) > 0 {
		keyAttributes[key.SKName], err = dynamodbattribute.Marshal(key.SKValue)
		if err != nil {
			return err
		}
	}

	pkCondition := expression.Key(key.PKName).Equal(expression.Value(key.PKValue))
	if len(key.SKName) > 0 {
		skCondition := expression.Key(key.SKName).Equal(expression.Value(key.SKValue))
		pkCondition = pkCondition.And(skCondition)
	}

	updateBuilder := expression.UpdateBuilder{}
	for fieldName, fieldValue := range fields {
		updateBuilder = updateBuilder.Set(expression.Name(fieldName), expression.Value(fieldValue))
	}

	expr, err := expression.NewBuilder().WithKeyCondition(pkCondition).WithUpdate(updateBuilder).Build()
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 aws.String(tablename),
		Key:                       keyAttributes,
		ConditionExpression:       expr.KeyCondition(),
		UpdateExpression:          expr.Update(),
	}

	_, err = svc.UpdateItem(input)

	return err
}

// GetItem retrieves from the table the item identified by its partition key (and sort key if given)
// then returns the item in the form of an instance of your choice.
//
// Arguments:
//
// table: tablename
//
// key: the partion key field name/value and sort key name/value(optional, can be empty)
//
// pointerToOutputObject: pointer to a struct or map[string]interface{} instance that will be filled with
// the data comming from dynamo
//
// Example:
//
// person := Person{}
//
// key := Key{PKName: "personId", PKValue: "10019911"}
//
// err = GetItem(tablename, key, &person)
//
// The errors returned are:
//     - ItemNotFoundException: no matching item was found on the database for the given Key
//     - errors from the aws sdk: see https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#DynamoDB.GetItem
func GetItem(tablename string, key Key, pointerToOutputObject interface{}) (err error) {
	svc := dynamodb.New(sessionutils.Session)

	keyAttributes := make(map[string]*dynamodb.AttributeValue)

	keyAttributes[key.PKName], err = dynamodbattribute.Marshal(key.PKValue)
	if err != nil {
		return err
	}

	if len(key.SKName) > 0 {
		keyAttributes[key.SKName], err = dynamodbattribute.Marshal(key.SKValue)
		if err != nil {
			return err
		}
	}

	getItemOutput, err := svc.GetItem(&dynamodb.GetItemInput{
		Key:       keyAttributes,
		TableName: aws.String(tablename),
	})
	if err != nil {
		return err
	}
	if len(getItemOutput.Item) == 0 {
		return errors.New("ItemNotFoundException")
	}

	err = dynamodbattribute.UnmarshalMap(getItemOutput.Item, pointerToOutputObject)

	return err
}

// PutItem creates or replaces an Item on a Dynamodb table.
// The given item must be a struct or a map[string]interface{} instance
func PutItem(tablename string, item interface{}) error {
	dynamoItem, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(tablename),
		Item:      dynamoItem,
	}

	dynamodbClient := dynamodb.New(sessionutils.Session)
	_, err = dynamodbClient.PutItem(putItemInput)

	return err
}
