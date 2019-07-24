package localstack

import "testing"

func TestInstanciarServicos(t *testing.T) {
	// t.Log(ApiGateway)
	// t.Log(ApiGateway + "oi")
	// t.Log(ApiGateway.String())

	// t.Log(ApiGateway.Endpoint())
	// t.Log(S3.Endpoint())

	t.Log(Services.IAM)
	t.Log(Services.SNS)

}
