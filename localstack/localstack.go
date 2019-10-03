package localstack

import (
	"bytes"
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
func StartLocalstack(services []string) error {
	cmd := exec.Command(os.Getenv("GOPATH") + "/src/stash.b2w/asp/aws-utils-go.git/localstack/runLocalstack.sh")
	newEnv := append(os.Environ(), "SERVICES="+strings.Join(services, ","))
	cmd.Env = newEnv
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	localstackPID = out.String()
	return nil
}

func StopLocalstack() error {
	if localstackPID == "" {
		return errors.New("localstackPID nao foi definido. rode StartLocalstack() primeiro.")
	}

	process := exec.Command("kill", "-2", localstackPID)
	_, err := process.Output()
	return err
}

func StartLocalstack2(services ...Service) error {
	var servicesNames []string
	for _, s := range services {
		servicesNames = append(servicesNames, s.Name)
	}
	var SERVICES_ENV_VAR = "SERVICES=" + strings.Join(servicesNames, ",")

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
