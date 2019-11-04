package localstack

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var localstackPID string = ""

func init() {
	checkIfLocalstackIsInstalled()
}

func checkIfLocalstackIsInstalled() {
	localstackCmd := exec.Command("localstack", "-v")
	_, err := localstackCmd.Output()
	if err != nil {
		panic(err)
	}
}

// startLocalstack - roda o localstack na maquina local ativando os serviços passados como argumento
// utiliza o script ../runLocalstack.sh para efetivamente iniciar o localstack
func StartLocalstack(serviceNames []string) error {
	var SERVICES_ENV_VAR = "SERVICES=" + strings.Join(serviceNames, ",")

	cmd := exec.Command(os.Getenv("GOPATH") + "/src/stash.b2w/asp/aws-utils-go.git/localstack/runLocalstack.sh")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, SERVICES_ENV_VAR)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error starting localstack with the script " + cmd.Path)
		fmt.Println("Script output: " + string(out))
		return err
	}
	localstackPID = string(out)
	return nil
}

func StartLocalstack2(services ...Service) error {
	var serviceNames []string
	for _, s := range services {
		serviceNames = append(serviceNames, s.Name)
	}
	return StartLocalstack(serviceNames)
}

func StopLocalstack() error {
	if localstackPID == "" {
		return errors.New("localstackPID nao foi definido. rode StartLocalstack() primeiro.")
	}

	cmd := exec.Command(os.Getenv("GOPATH") + "/src/stash.b2w/asp/aws-utils-go.git/localstack/stopLocalstack.sh")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
