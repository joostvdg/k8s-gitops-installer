package namespace

import (
    "os/exec"
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
    log "github.com/sirupsen/logrus"
)

func SetDefault(namespace string) {
    log.Infof("Setting namespace %s as default context\n", namespace)
    currentContextCmd := exec.Command("kubectl",
        "config", "current-context",
    )

    currentContext := util.RunCmd(currentContextCmd, true)
    setCurrentContextToNamespaceCmd := exec.Command("kubectl",
        "config", "set-context", currentContext, "--namespace="+namespace,
    )
    util.RunCmd(setCurrentContextToNamespaceCmd, true)
}
