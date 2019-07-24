package localstack

type Service struct {
	Name string
	Port string
}

func (s Service) EndpointUrl() string {
	return "http://localhost:" + s.Port
}

type suportedServices struct {
	ApiGateway           Service
	CloudFormation       Service
	CloudWatchLogs       Service
	CloudWatch           Service
	DynamoDB             Service
	DynamoDBStreams      Service
	Elasticsearch        Service
	ElasticsearchService Service
	Firehose             Service
	IAM                  Service
	Kinesis              Service
	Lambda               Service
	Redshift             Service
	Route53              Service
	S3                   Service
	SecretsManager       Service
	SecurityTokenService Service
	SimpleEmailService   Service
	SNS                  Service
	SQS                  Service
	StepFunctions        Service
	SystemsManager       Service
}

var Services = suportedServices{
	ApiGateway: Service{
		Name: "apigateway",
		Port: "4567",
	},
	CloudFormation: Service{
		Name: "cloudformation",
		Port: "4581",
	},
	CloudWatch: Service{
		Name: "cloudwatch",
		Port: "4582",
	},
	CloudWatchLogs: Service{
		Name: "logs",
		Port: "4586",
	},
	DynamoDB: Service{
		Name: "dynamodb",
		Port: "4569",
	},
	DynamoDBStreams: Service{
		Name: "dynamodbstreams",
		Port: "4570",
	},
	Elasticsearch: Service{
		Name: "I_DONT_KNOW",
		Port: "4571",
	},
	ElasticsearchService: Service{
		Name: "es",
		Port: "4578",
	},
	Firehose: Service{
		Name: "firehose",
		Port: "4573",
	},
	IAM: Service{
		Name: "iam",
		Port: "4593",
	},
	Kinesis: Service{
		Name: "kinesis",
		Port: "4568",
	},
	Lambda: Service{
		Name: "lambda",
		Port: "4574",
	},
	Redshift: Service{
		Name: "redshift",
		Port: "4577",
	},
	Route53: Service{
		Name: "route53",
		Port: "4580",
	},
	S3: Service{
		Name: "s3",
		Port: "4572",
	},
	SNS: Service{
		Name: "sns",
		Port: "4575",
	},
	SQS: Service{
		Name: "sqs",
		Port: "4576",
	},
	SystemsManager: Service{
		Name: "ssm",
		Port: "4583",
	},
	SecurityTokenService: Service{
		Name: "sts",
		Port: "4592",
	},
	SecretsManager: Service{
		Name: "secretsmanager",
		Port: "4584",
	},
	SimpleEmailService: Service{
		Name: "ses",
		Port: "4579",
	},
	StepFunctions: Service{
		Name: "stepfunctions",
		Port: "4585",
	},
}
