package sessionutils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var Session *session.Session

func init() {
	// for _, pair := range os.Environ() {
	// 	fmt.Println(pair)
	// }

	creds := credentials.NewEnvCredentials()
	//fmt.Println("AAAAAAAAAAAAAAAAAAa" + os.Getenv("AWS_SDK_LOAD_CONFIG"))
	//value, _ := creds.Get()
	//fmt.Printf("credentials da aws: %+v\n", value)
	Session = session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
	}))
}
