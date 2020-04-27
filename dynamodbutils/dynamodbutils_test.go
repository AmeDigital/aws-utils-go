package dynamodbutils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/AmeDigital/aws-utils-go/localstack"
	"github.com/AmeDigital/aws-utils-go/sessionutils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dynamodbClient *dynamodb.DynamoDB
var tablename = "cities"
var indexname = "NameToPkSk"

var check func(e error) = func(e error) {
	if e != nil {
		localstack.StopLocalstack()
		fmt.Println(e.Error())
		panic(e)
	}
}

// TestMain roda em volta de cada teste executado. Os testes são executados
// na invocação de 'm.Run()'
func TestMain(m *testing.M) {
	check = func(e error) {
		if e != nil {
			localstack.StopLocalstack()
			panic(e)
		}
	}

	// cria recursos no localstack,
	err := localstack.StartLocalstack2(localstack.Services.DynamoDB)
	check(err)

	// configures dynamodb client to use localstack
	awsConfigForDynamodb := aws.Config{Endpoint: aws.String("http://localhost:4569"), Region: aws.String("us-east-1")}
	dynamodbSessionForLocalstack, err := session.NewSession(&awsConfigForDynamodb)
	sessionutils.Session = dynamodbSessionForLocalstack
	check(err)
	dynamodbClient = dynamodb.New(dynamodbSessionForLocalstack)

	// creates the table for testing
	createTable(tablename)

	// executa os testes
	returnCode := m.Run()

	// desliga o localstack
	dynamodbClient.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: &tablename,
	})
	localstack.StopLocalstack()

	os.Exit(returnCode)
}

func createTable(tablename string) {
	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("State"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Name"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("State"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String("RANGE"),
			},
		},
		TableName: aws.String(tablename),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String(indexname),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("Name"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String(dynamodb.ProjectionTypeAll),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},
	}

	_, err := dynamodbClient.CreateTable(createTableInput)

	check(err)
}

type City struct {
	State      string
	Id         int
	Name       string
	Population int
	Aliases    []string
}

func TestPutItemUpdateItemAndGetItem(t *testing.T) {

	var err error

	city := City{
		State:      "NJ",
		Id:         1,
		Name:       "Wayne",
		Population: 351,
		Aliases:    []string{"TomTown", "JohnTown"},
	}

	err = PutItem(tablename, city)
	check(err)

	newAliases := [2]string{"RRRR", "ZZZZZ"}

	updateFields := make(map[string]interface{})
	updateFields["Population"] = 360
	updateFields["Aliases"] = newAliases

	key := Key{PKName: "State", PKValue: "NJ", SKName: "Id", SKValue: 1}

	err = UpdateItem(tablename, key, updateFields)
	check(err)

	newCity := City{}

	err = GetItem(tablename, key, &newCity)
	check(err)

	if newCity.Population != 360 {
		t.Errorf("Population shoud be %d but was %d", 360, newCity.Population)
	}

	var arr [2]string
	copy(arr[:], newCity.Aliases)

	if arr != newAliases {
		t.Errorf("Aliases should be %s but was %s", newAliases, newCity.Aliases)
	}
}

func TestPutItemWithConditionalAvoidOverrideExistingItem(t *testing.T) {
	city := City{
		State:      "MG",
		Id:         100,
		Name:       "Tiradentes",
		Population: 50000,
		Aliases:    []string{"Tira"},
	}

	err := PutItemWithConditional(tablename, city, "attribute_not_exists(Id)", nil)
	check(err)

	err = PutItemWithConditional(tablename, city, "attribute_not_exists(Id)", nil)
	if err == nil {
		t.Error("should not accept overrite item")
	} else if !strings.Contains(err.Error(), "ConditionalCheckFailedException") {
		t.Errorf("error should be of type 'ConditionalCheckFailedException' but was %s\n", err.Error())
	}
}

func TestFindOneFromIndex(t *testing.T) {

	var err error

	expected := City{
		State:      "NJ",
		Id:         2,
		Name:       "Hollow",
		Population: 1000,
		Aliases:    []string{"hollow town"},
	}

	another := City{
		State:      "NJ",
		Id:         3,
		Name:       "Prank",
		Population: 3000,
		Aliases:    []string{"prank town"},
	}

	err = PutItem(tablename, expected)
	check(err)
	err = PutItem(tablename, another)
	check(err)

	found := City{}

	err = FindOneFromIndex(tablename, indexname, Key{PKName: "Name", PKValue: "Hollow"}, &found)
	check(err)

	if !reflect.DeepEqual(expected, found) {
		t.Error(fmt.Sprintf("Expected: %+v, Result: %+v", expected, found))
	}
}

// func constructor() interface{} {
// 	return &City{}
// }

func TestPointers(t *testing.T) {
	var product interface{} = City{}
	productType := reflect.TypeOf(product)       // this type of this variable is reflect.Type
	productPointer := reflect.New(productType)   // this type of this variable is reflect.Value.
	productValue := productPointer.Elem()        // this type of this variable is reflect.Value.
	productInterface := productValue.Interface() // this type of this variable is interface{}
	product2 := productInterface.(City)          // this type of this variable is product
	product2.Name = "BH"
}

func TestQuery(t *testing.T) {
	bh := City{
		Name:  "Belo Horizonte",
		Id:    1,
		State: "MG",
	}

	divinopolis := City{
		Name:  "Divinópolis",
		Id:    2,
		State: "MG",
	}

	PutItem(tablename, bh)
	PutItem(tablename, divinopolis)

	var cities = []City{}

	// testing when sort key is equal to a value
	keyCondition := KeyCondition{
		PKName:       "State",
		PKValue:      "MG",
		SKName:       "Id",
		SKValueEqual: 1,
	}

	err := Query(tablename, keyCondition, &cities)

	if err != nil {
		t.Error("Query() failed with error: " + err.Error())
	} else if len(cities) != 1 {
		t.Error(fmt.Sprintf("cities should have length 1 but has %d", len(cities)))
	} else {
		city := cities[0]
		if city.Name != "Belo Horizonte" {
			t.Error(fmt.Sprintf("City name should be 'Belo Horizonte' but was '%s'", city.Name))
		}
		if city.State != "MG" {
			t.Error(fmt.Sprintf("City state should be 'MG' but was '%s'", city.State))
		}
		if city.Id != 1 {
			t.Error(fmt.Sprintf("City id should be '1' but was '%d'", city.Id))
		}
	}

	// testing when sort key is equal to a value on the INDEX table
	keyCondition = KeyCondition{
		IndexName: indexname,
		PKName:    "Name",
		PKValue:   "Belo Horizonte",
	}

	err = Query(tablename, keyCondition, &cities)

	if err != nil {
		t.Error("Query() failed with error: " + err.Error())
	} else if len(cities) != 1 {
		t.Error(fmt.Sprintf("cities should have length 1 but has %d", len(cities)))
	} else if !reflect.DeepEqual(cities[0], bh) {
		t.Errorf("A busca no indice falhou: %+v", cities[0])
	}

	// testing when sork key is greater than a value
	keyCondition = KeyCondition{
		PKName:             "State",
		PKValue:            "MG",
		SKName:             "Id",
		SKValueGreaterThan: 1,
	}

	err = Query(tablename, keyCondition, &cities)

	if err != nil {
		t.Error("Query() failed with error: " + err.Error())
	} else if len(cities) != 1 {
		t.Error(fmt.Sprintf("cities should have length 1 but has %d", len(cities)))
	} else {
		city := cities[0]
		if city.Name != "Divinópolis" {
			t.Error(fmt.Sprintf("City name should be 'Divinópolis' but was '%s'", city.Name))
		}
	}

	// testing when sork key is greater than equal a value
	keyCondition = KeyCondition{
		PKName:                  "State",
		PKValue:                 "MG",
		SKName:                  "Id",
		SKValueGreaterThanEqual: 1,
	}

	err = Query(tablename, keyCondition, &cities)

	if err != nil {
		t.Error("Query() failed with error: " + err.Error())
	} else if len(cities) != 2 {
		t.Error(fmt.Sprintf("cities should have length 2 but has %d", len(cities)))
	} else {
		city := cities[0]
		if city.Id != 1 {
			t.Error(fmt.Sprintf("City Id should be '1' but was '%d'", city.Id))
		}
		city = cities[1]
		if city.Id != 2 {
			t.Error(fmt.Sprintf("City Id should be '2' but was '%d'", city.Id))
		}
	}

	// testing when sork key is less than a value
	keyCondition = KeyCondition{
		PKName:          "State",
		PKValue:         "MG",
		SKName:          "Id",
		SKValueLessThan: 2,
	}

	err = Query(tablename, keyCondition, &cities)

	if err != nil {
		t.Error("Query() failed with error: " + err.Error())
	} else if len(cities) != 1 {
		t.Error(fmt.Sprintf("cities should have length 1 but has %d", len(cities)))
	} else {
		city := cities[0]
		if city.Id != 1 {
			t.Error(fmt.Sprintf("City id should be '1' but was '%d'", city.Id))
		}
	}

	// testing when sork key is less than equal a value
	keyCondition = KeyCondition{
		PKName:               "State",
		PKValue:              "MG",
		SKName:               "Id",
		SKValueLessThanEqual: 2,
	}

	err = Query(tablename, keyCondition, &cities)

	if err != nil {
		t.Error("Query() failed with error: " + err.Error())
	} else if len(cities) != 2 {
		t.Error(fmt.Sprintf("cities should have length 3 but has %d", len(cities)))
	} else {
		city := cities[0]
		if city.Id != 1 {
			t.Error(fmt.Sprintf("City id should be '1' but was '%d'", city.Id))
		}
		city = cities[1]
		if city.Id != 2 {
			t.Error(fmt.Sprintf("City id should be '2' but was '%d'", city.Id))
		}
	}

	// testing when sort key is between two values
	ouroPreto := City{
		Name:  "Ouro Preto",
		Id:    3,
		State: "MG",
	}

	diamantina := City{
		Name:  "Diamantina",
		Id:    4,
		State: "MG",
	}

	PutItem(tablename, ouroPreto)
	PutItem(tablename, diamantina)

	keyCondition = KeyCondition{
		PKName:              "State",
		PKValue:             "MG",
		SKName:              "Id",
		SKValueBetweenStart: 2,
		SKValueBetweenEnd:   4,
	}

	err = Query(tablename, keyCondition, &cities)

	if err != nil {
		t.Error("Query() failed with error: " + err.Error())
	} else if len(cities) != 3 {
		t.Error(fmt.Sprintf("cities should have length 3 but has %d", len(cities)))
	} else {
		city := cities[0]
		if city.Id != 2 {
			t.Error(fmt.Sprintf("City id should be '2' but was '%d'", city.Id))
		}
		city = cities[1]
		if city.Id != 3 {
			t.Error(fmt.Sprintf("City id should be '3' but was '%d'", city.Id))
		}
		city = cities[2]
		if city.Id != 4 {
			t.Error(fmt.Sprintf("City id should be '4' but was '%d'", city.Id))
		}
	}

	// test a non pointer output slice argument
	err = Query(tablename, KeyCondition{}, []City{})
	if err == nil {
		t.Error("err should be depicting that argument should be a pointer to a slice")
	}

	// test a pointer to something that is not a slice
	var dummy = "dummy"
	err = Query(tablename, KeyCondition{}, &dummy)
	if err == nil {
		t.Error("err should be depicting that argument should be a pointer to a slice")
	}

}

func TestGetInexistentItem(t *testing.T) {
	newCity := City{}
	key := Key{PKName: "State", PKValue: "bleh bleh", SKName: "Id", SKValue: 888}
	err := GetItem(tablename, key, &newCity)

	if err == nil {
		t.Error("err object should not be nil")
	} else if err.Error() != "ItemNotFoundException" {
		t.Errorf("err should be 'ItemNotFoundException' but was '%s'.\n", err.Error())
	}
}

func objectToJsonString(obj interface{}) string {
	b, err := json.Marshal(obj)
	check(err)
	return string(b)
}

func TestType(t *testing.T) {
	city := City{
		Id: 1,
	}

	get(&city)
}

func get(items interface{}) {
	itemType := reflect.TypeOf(items)
	fmt.Println(itemType)
}

func TestBatchGetItem(t *testing.T) {
	bh := City{
		Name:  "Belo Horizonte",
		Id:    1,
		State: "MG",
	}

	divinopolis := City{
		Name:  "Divinópolis",
		Id:    2,
		State: "MG",
	}

	ouroPreto := City{
		Name:  "Ouro Preto",
		Id:    3,
		State: "MG",
	}

	PutItem(tablename, bh)
	PutItem(tablename, divinopolis)
	PutItem(tablename, ouroPreto)

	cities := []City{}

	keys := []Key{
		Key{
			PKName:  "State",
			PKValue: "MG",
			SKName:  "Id",
			SKValue: 1,
		},
		Key{
			PKName:  "State",
			PKValue: "MG",
			SKName:  "Id",
			SKValue: 2,
		},
	}

	err := BatchGetItem(tablename, keys, &cities)

	if err != nil {
		t.Error("BatchGetItem() failed with error: " + err.Error())
	} else if len(cities) != 2 {
		t.Error(fmt.Sprintf("cities should have length 2 but has %d", len(cities)))
	} else {
		if !reflect.DeepEqual(cities[0], bh) {
			t.Error("city 0 is not bh")
		} else if !reflect.DeepEqual(cities[1], divinopolis) {
			t.Error("city 1 is not divinopolis")
		}
	}

	// test query for inexistent item
	keys = []Key{
		Key{
			PKName:  "State",
			PKValue: "MG",
			SKName:  "Id",
			SKValue: 100,
		},
	}

	cities = []City{}

	err = BatchGetItem(tablename, keys, &cities)
	if err != nil {
		t.Error("BatchGetItem() failed with error: " + err.Error())
	} else if len(cities) != 0 {
		t.Error(fmt.Sprintf("cities should have length 0 but had %d", len(cities)))
	}

}
