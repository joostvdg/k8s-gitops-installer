package namespace

import (
    "os/exec"
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
    log "github.com/sirupsen/logrus"
)

func Create(name string) {
    log.Infof("Creating namespace %s\n", name)
    namespaceExistsCmd := exec.Command("kubectl",
        "get", "namespace", name)
    createNamespaceCmd := exec.Command("kubectl",
        "create", "namespace", "cje",
    )
    addNamespaceLabelCmd := exec.Command("kubectl",
        "label", "namespace", "cje", "name="+name,
    )

    namespaceExists := util.RunCmdNonFatal(namespaceExistsCmd)
    if  namespaceExists {
        log.Infof("Namespace %s already exists\n", name)
    } else {
        log.Info("Create namespace")
        util.RunCmd(createNamespaceCmd, true)
    }

    log.Info("Add namespace label(s)")
    util.RunCmd(addNamespaceLabelCmd, false)
}
