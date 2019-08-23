package dynamodbutils

import (
	"encoding/json"
	"fmt"
	"localstack"
	"os"
	"reflect"
	"testing"

	"stash.b2w/asp/aws-utils-go.git/sessionutils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dynamodbClient *dynamodb.DynamoDB
var tablename = "cities"

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
	err := localstack.StartLocalstack([]string{"dynamodb"})
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
