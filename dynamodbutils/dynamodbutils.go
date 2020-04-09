package dynamodbutils

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/AmeDigital/aws-utils-go/sessionutils"

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

// DeleteItem - deletes an item from dynamodb
// Note: this function won't return error if the item was not found on the table.
func DeleteItem(tablename string, key Key) (err error) {
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

	_, err = svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &tablename,
		Key:       keyAttributes,
	})

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
// type Person struct {
//     // it has to start with Uppercase letter otherwise GetItem will silently fail!
//     Id int           `json:"id"`
//	   // it has to start with Uppercase letter otherwise GetItem will silently fail!
//     Name string      `json:"name"`
// }
//
// person := Person{}
//
// key := Key{PKName: "personId", PKValue: "10019911"}
//
// err = GetItem(tablename, key, &person)
//
// The errors returned are:
//     - ItemNotFoundException: no matching item was found on the database for the given Key.
//		 This error indicates that the database was queried successfully but the item does not exist.
//		 Note: use 'err.Error() == "ItemNotFoundException"' to identify this error.
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

// FindOneFromIndex works like GetItem() but runs a query to a secondary index table in order to find the item.
// This method is meant to be used when the given key will match a single item from the index and
// it will throw a 'MultipleItemsFound' error if the query returns more than one item.
//
// The errors returned are:
//     - ItemNotFoundException: no matching item was found on the database for the given Key.
//		 This error indicates that the database was queried successfully but the item does not exist.
//		 Note: use 'err.Error() == "ItemNotFoundException"' to identify this error.
// 	   - MultipleItemsFound: the query retrieved more than one item.
//		 Note: use 'err.Error() == "MultipleItemsFound"' to identify this error.
//     - errors from the aws sdk: see https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#DynamoDB.GetItem
func FindOneFromIndex(tablename string, indexname string, key Key, pointerToOutputObject interface{}) (err error) {
	svc := dynamodb.New(sessionutils.Session)

	keyCondition := expression.Key(key.PKName).Equal(expression.Value(key.PKValue))

	if len(key.SKName) > 0 {
		keyCondition = expression.KeyAnd(keyCondition, expression.Key(key.SKName).Equal(expression.Value(key.SKValue)))
	}

	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()

	if err != nil {
		return err
	}

	queryOutput, err := svc.Query(&dynamodb.QueryInput{
		TableName:                 &tablename,
		IndexName:                 &indexname,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return err
	}

	if len(queryOutput.Items) == 0 {
		return errors.New("ItemNotFoundException")
	} else if len(queryOutput.Items) > 1 {
		return errors.New("MultipleItemsFound")
	}

	err = dynamodbattribute.UnmarshalMap(queryOutput.Items[0], pointerToOutputObject)

	return err
}

// PutItem creates or replaces an Item on a Dynamodb table.
// The given item must be a struct or a map[string]interface{} instance
func PutItem(tablename string, item interface{}) error {
	return PutItemWithConditional(tablename, item, "", nil)
}

// PutItemWithConditional put item with conditional
// example:
// queryConditional := "deleted = :deleted"
// valuesConditional := map[string]interface{}{":deleted": false}
// err := dynamodbutils.PutItemWithConditional(PROMOTION_TABLE_NAME, promotionPersisted, queryConditional, valuesConditional)
func PutItemWithConditional(tablename string, item interface{}, conditionalExpression string, conditionalValues map[string]interface{}) error {
	dynamoItem, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	var putItemInput *dynamodb.PutItemInput

	var condExp *string
	if len(conditionalExpression) > 0 {
		condExp = &conditionalExpression
	}

	var condValues map[string]*dynamodb.AttributeValue
	if len(conditionalValues) > 0 {
		condValues, err = dynamodbattribute.MarshalMap(conditionalValues)
		if err != nil {
			return err
		}
	}

	putItemInput = &dynamodb.PutItemInput{
		TableName:                 aws.String(tablename),
		Item:                      dynamoItem,
		ConditionExpression:       condExp,
		ExpressionAttributeValues: condValues,
	}

	dynamodbClient := dynamodb.New(sessionutils.Session)
	_, err = dynamodbClient.PutItem(putItemInput)

	return err
}

// KeyCondition allows you set the parameters for a query with 'key condition expression'
// ref https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.html#Query.KeyConditionExpressions
//   - IndexName: optional. If this value is set the query will run on the index table.
//   - PKName: primary key name, mandatory.
//   - PKValue: primary key value, mandatory.
//   - SKName: sort key name, optional. You must set this value if you what to use the conditions that operate on the sort keys.
//   - SKValueEqual: selects an item with the sort key value equal to the given value
//   - SKValueLessThan: selects items having the sort key value less than the given value
//   - SKValueLessThanEqual: selects items having the sort key value less than or equal to the the given value
//   - SKValueGreaterThan: selects items having the sort key value greater than the given value
//   - SKValueGreaterThanEqual: selects items having the sort key value greater than or equal to the the given value
//   - SKValueBetweenStart and SKValueBetweenEnd: selects items having the sort key value between the given limits, including the limiting items.
type KeyCondition struct {
	IndexName               string      // optional
	PKName                  string      // mandatory
	PKValue                 interface{} // mandatory
	SKName                  string      // optional
	SKValueEqual            interface{} // optional
	SKValueLessThan         interface{} // optional
	SKValueLessThanEqual    interface{} // optional
	SKValueGreaterThan      interface{} // optional
	SKValueGreaterThanEqual interface{} // optional
	SKValueBetweenStart     interface{} // optional
	SKValueBetweenEnd       interface{} // optional
}

// Runs the query specified by the keyCondition argument on the given table or index and fills the slice
// pointed by 'pointerToOutputSlice' with the items found, if any.
func Query(tablename string, keyCondition KeyCondition, pointerToOuputSlice interface{}) (err error) {
	rv := reflect.ValueOf(pointerToOuputSlice)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dynamodbutils.Query: pointerToOutputSlice must be a slice pointer")
	}

	dynamodbClient := dynamodb.New(sessionutils.Session)

	attributeValues := make(map[string]*dynamodb.AttributeValue)
	attributeNames := make(map[string]*string)

	var keyConditionExpression = "#pkname = :pkval"

	attributeNames["#pkname"] = &keyCondition.PKName

	attributeValues[":pkval"], err = dynamodbattribute.Marshal(keyCondition.PKValue)
	if err != nil {
		return err
	}

	if len(keyCondition.SKName) > 0 {
		var err error

		attributeNames["#skname"] = &keyCondition.SKName

		if keyCondition.SKValueEqual != nil {
			keyConditionExpression = keyConditionExpression + " and #skname = :skval"
			attributeValues[":skval"], err = dynamodbattribute.Marshal(keyCondition.SKValueEqual)

		} else if keyCondition.SKValueBetweenStart != nil && keyCondition.SKValueBetweenEnd != nil {
			keyConditionExpression = keyConditionExpression + " and #skname BETWEEN :skval1 AND :skval2"
			attributeValues[":skval1"], err = dynamodbattribute.Marshal(keyCondition.SKValueBetweenStart)
			if err == nil {
				attributeValues[":skval2"], err = dynamodbattribute.Marshal(keyCondition.SKValueBetweenEnd)
			}

		} else if keyCondition.SKValueGreaterThan != nil {
			keyConditionExpression = keyConditionExpression + " and #skname > :skval"
			attributeValues[":skval"], err = dynamodbattribute.Marshal(keyCondition.SKValueGreaterThan)

		} else if keyCondition.SKValueGreaterThanEqual != nil {
			keyConditionExpression = keyConditionExpression + " and #skname >= :skval"
			attributeValues[":skval"], err = dynamodbattribute.Marshal(keyCondition.SKValueGreaterThanEqual)

		} else if keyCondition.SKValueLessThan != nil {
			keyConditionExpression = keyConditionExpression + " and #skname < :skval"
			attributeValues[":skval"], err = dynamodbattribute.Marshal(keyCondition.SKValueLessThan)

		} else if keyCondition.SKValueLessThanEqual != nil {
			keyConditionExpression = keyConditionExpression + " and #skname <= :skval"
			attributeValues[":skval"], err = dynamodbattribute.Marshal(keyCondition.SKValueLessThanEqual)

		} else {
			return errors.New("keyCondition is invalid")
		}
		if err != nil {
			return err
		}
	}

	queryInput := dynamodb.QueryInput{
		TableName:                 &tablename,
		KeyConditionExpression:    &keyConditionExpression,
		ExpressionAttributeValues: attributeValues,
		ExpressionAttributeNames:  attributeNames,
	}

	if len(keyCondition.IndexName) > 0 {
		queryInput.IndexName = &keyCondition.IndexName
	}

	queryOutput, err := dynamodbClient.Query(&queryInput)

	if err != nil {
		return err
	}

	if len(queryOutput.Items) == 0 {
		return nil
	}

	err = dynamodbattribute.UnmarshalListOfMaps(queryOutput.Items, pointerToOuputSlice)

	return err
}

// Retrieves a list of items identified by their keys from the given table and fills the slice
// pointed by 'pointerToOutputSlice' with the items found, if any
func BatchGetItem(tablename string, keys []Key, pointerToOuputSlice interface{}) (err error) {
	rv := reflect.ValueOf(pointerToOuputSlice)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dynamodbutils.BatchGetItem: pointerToOutputSlice must be a slice pointer")
	}

	dynamodbClient := dynamodb.New(sessionutils.Session)

	keyAttributesListOfMaps := []map[string]*dynamodb.AttributeValue{}

	for _, key := range keys {
		keyAttributesMap := make(map[string]*dynamodb.AttributeValue)

		keyAttributesMap[key.PKName], err = dynamodbattribute.Marshal(key.PKValue)
		if err != nil {
			return err
		}

		if len(key.SKName) > 0 {
			keyAttributesMap[key.SKName], err = dynamodbattribute.Marshal(key.SKValue)
			if err != nil {
				return err
			}
		}

		keyAttributesListOfMaps = append(keyAttributesListOfMaps, keyAttributesMap)
	}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tablename: {
				Keys: keyAttributesListOfMaps,
			},
		},
	}

	result, err := dynamodbClient.BatchGetItem(input)
	if err != nil {
		return err
	}

	resultsForThisTable := result.Responses[tablename]

	if len(resultsForThisTable) > 0 {
		err = dynamodbattribute.UnmarshalListOfMaps(resultsForThisTable, pointerToOuputSlice)
	}

	return err
}
