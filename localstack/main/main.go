package main

import (
	"fmt"
	"localstack"
)

func main() {
	var services = localstack.Services
	fmt.Println(services.Firehose)
	fmt.Println(services.Firehose.EndpointUrl())
	localstack.StartLocalstack2(services.SNS, services.SQS)
}
