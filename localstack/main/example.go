package main

import (
	"fmt"

	"github.com/AmeDigital/aws-utils-go/localstack"
)

func main() {
	var services = localstack.Services
	fmt.Println(services.Firehose)
	fmt.Println(services.Firehose.EndpointUrl())
	localstack.StartLocalstack2(services.SNS, services.SQS)
}
