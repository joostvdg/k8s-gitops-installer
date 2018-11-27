package validate

import (
    log "github.com/sirupsen/logrus"
    "os/exec"
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
)

func Kubectl() {
    testCmd := exec.Command("kubectl", "get", "nodes", "-o", "wide" )
    log.Info("Testing kubectl command")
    util.RunCmd(testCmd, true)
}
