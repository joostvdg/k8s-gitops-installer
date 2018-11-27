package helm

import (
    "os/exec"
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
    log "github.com/sirupsen/logrus"
)

func Install() {
    log.Info("Installing Helm 2 (with Tiller)")
    createServiceAccountCmd := exec.Command("kubectl",
        "create", "serviceaccount", "--namespace", "kube-system", "tiller",
    )
    createClusterRoleBindingCmd := exec.Command("kubectl",
        "create", "clusterrolebinding", "tiller-cluster-rule", "--clusterrole=cluster-admin", "--serviceaccount=kube-system:tiller",
    )
    helmInitCmd := exec.Command("helm",
        "init", "--service-account", "tiller", "--upgrade",
    )

    log.Info("Creating tiller service account (SA)")
    util.RunCmd(createServiceAccountCmd, false)
    log.Info("Creating cluster role binding for tiller SA")
    util.RunCmd(createClusterRoleBindingCmd, false)
    log.Info("Initialize Helm")
    util.RunCmd(helmInitCmd, true)
}
