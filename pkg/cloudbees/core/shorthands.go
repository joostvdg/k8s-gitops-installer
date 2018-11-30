package core

import (
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
    log "github.com/sirupsen/logrus"
    "os/exec"
)

func GetCjocPassword() {
    printCBCorePasswordCmd := exec.Command("kubectl",
        "exec", "cjoc-0", "--", "cat", "/var/jenkins_home/secrets/initialAdminPassword",
    )
    log.Info("Retrieve Operations Center password")
    util.RunCmd(printCBCorePasswordCmd, false)
}
