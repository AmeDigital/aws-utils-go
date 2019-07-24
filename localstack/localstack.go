package localstack

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
)

var localstackPID string = ""

// startLocalstack - roda o localstack na maquina local ativando os serviços passados como argumento
// utiliza o script ../runLocalstack.sh para efetivamente iniciar o localstack
func StartLocalstack(services []string) error {
	cmd := exec.Command(os.Getenv("GOPATH") + "/src/localstack/runLocalstack.sh")
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
