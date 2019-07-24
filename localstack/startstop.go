package localstack

import (
	"os"
	"os/exec"
	"strings"
)

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

func StartLocalstack2(services ...Service) {
	var servicesNames []string
	for _, s := range services {
		servicesNames = append(servicesNames, s.Name)
	}
	var SERVICES_ENV_VAR = "SERVICES=" + strings.Join(servicesNames, ",")

	localstackCmd := exec.Command("localstack", "start")
	localstackCmd.Env = os.Environ()
	localstackCmd.Env = append(localstackCmd.Env, SERVICES_ENV_VAR)
}
