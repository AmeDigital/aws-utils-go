package localstack

// // API Gateway at http://localhost:4567
// // Kinesis at http://localhost:4568
// // DynamoDB at http://localhost:4569
// // DynamoDB Streams at http://localhost:4570
// // Elasticsearch at http://localhost:4571
// // S3 at http://localhost:4572
// // Firehose at http://localhost:4573
// // Lambda at http://localhost:4574
// // SNS at http://localhost:4575
// // SQS at http://localhost:4576
// // Redshift at http://localhost:4577
// // ES (Elasticsearch Service) at http://localhost:4578
// // SES at http://localhost:4579
// // Route53 at http://localhost:4580
// // CloudFormation at http://localhost:4581
// // CloudWatch at http://localhost:4582
// // SSM at http://localhost:4583
// // SecretsManager at http://localhost:4584
// // StepFunctions at http://localhost:4585
// // CloudWatch Logs at http://localhost:4586
// // STS at http://localhost:4592
// // IAM at http://localhost:4593

// type Service string

// const (
// 	ApiGateway           Service = "apigateway"
// 	Kinesis              Service = "kinesis"
// 	DynamoDB             Service = "dynamodb"
// 	DynamoDBStreams      Service = "dynamodbstreams"
// 	S3                   Service = "s3"
// 	Firehose             Service = "firehose"
// 	Lambda               Service = "lambda"
// 	SNS                  Service = "sns"
// 	SQS                  Service = "sqs"
// 	Redshift             Service = "redshift"
// 	ElasticsearchService Service = "es"
// 	SES                  Service = "ses"
// 	Route53              Service = "route53"
// 	CloudFormation       Service = "cloudformation"
// 	CloudWatch           Service = "cloudwatch"
// 	SSM                  Service = "ssm"
// 	SecretsManager       Service = "secretsmanager"
// 	StepFunctions        Service = "stepfunctions"
// )

// func (service Service) String() string {
// 	return string(service)
// }

// func (service Service) Endpoint() string {
// 	switch service {
// 	case ApiGateway:
// 		return "http://localhost:4567"
// 	default:
// 		return "urlnotfound"
// 	}
// }
