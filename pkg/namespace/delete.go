package namespace

import (
    "os/exec"
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
    log "github.com/sirupsen/logrus"
)

func Delete(name string) {
    log.Infof("Creating namespace %s\n", name)
    namespaceExistsCmd := exec.Command("kubectl",
        "get", "namespace", name)
    createNamespaceCmd := exec.Command("kubectl",
        "delete", "namespace", name,
    )

    namespaceExists := util.RunCmdNonFatal(namespaceExistsCmd)
    if  namespaceExists {
        log.Info("Delete namespace")
        util.RunCmd(createNamespaceCmd, true)

    } else {
        log.Infof("Namespace %s does not exists\n", name)
    }

}
