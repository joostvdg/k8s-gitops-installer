package ingress

import (
    "os/exec"
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
    log "github.com/sirupsen/logrus"
)

func InstallNginx() {
    log.Info("Installing Nginx as Ingress Controller")
    installMandatoryBase := exec.Command("kubectl",
        "apply", "-f", "https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml",
    )
    installCloudGeneric := exec.Command("kubectl",
        "apply", "-f", "https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/provider/cloud-generic.yaml",
    )

    log.Info("Installing Mandatory Nginx resources")
    util.RunCmd(installMandatoryBase, true)
    log.Info("Installing Additional (Cloud) Nginx resources")
    util.RunCmd(installCloudGeneric, true)
}
