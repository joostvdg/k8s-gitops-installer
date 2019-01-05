package letsencrypt

import (
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
    log "github.com/sirupsen/logrus"
    "os/exec"
)

func CleanCRDS() {
    cleanCRDSCmd := exec.Command("kubectl",
        "delete", "customresourcedefinitions.apiextensions.k8s.io", "certificates.certmanager.k8s.io", "clusterissuers.certmanager.k8s.io", "issuers.certmanager.k8s.io",
    )
    log.Info("Delete Let's Encrypt CRDS's")
    util.RunCmd(cleanCRDSCmd, false)
}

func PurgeHelm() {
    cleanHelmCmd := exec.Command("helm",
        "delete", "cert-manager", "--purge",
    )
    log.Info("Delete and Purge Certmanager's Helm install")
    util.RunCmd(cleanHelmCmd, false)
}
