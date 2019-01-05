package letsencrypt

import (
    "fmt"
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "os/exec"

    //"os/exec"
    //"github.com/joostvdg/k8s-gitops-installer/pkg/util"
)

func InstallWithCertmanager(email string, prod bool) {
    log.Infof("Installing cert-manager and ClusterIssuer for Let's Encrypt (email: %s)\n", email )
    clusterIssuerFilename := "cert-manager-cluster-issuer.yml"
    helmUpgradeCmd := exec.Command("helm",
       "upgrade", "cert-manager", "--install", "--version", "v0.5.1", "--namespace", "default", "stable/cert-manager",
    )
    helmExistsCmd := exec.Command("helm","ls")
    applyClusterIssuerResourceCmd := exec.Command("kubectl",
       "apply", "-f", clusterIssuerFilename,
    )

    helmExists := util.RunCmdNonFatal(helmExistsCmd)
    if !helmExists {
       log.Fatal("Helm not installed or not usable, cannot continue")
    }

    log.Info("Installing cert-manager via Helm (Upgrade --install)")
    util.RunCmd(helmUpgradeCmd, true)

    var clusterIssuerYaml string
    if prod {
        clusterIssuerYaml = fmt.Sprintf(Cert_manager_issuer_prod, email)
    } else {
        clusterIssuerYaml = fmt.Sprintf(Cert_manager_issuer_stage, email)
    }
    log.Infof("Cluster Issuer Yaml: \n%s\n", clusterIssuerYaml)
    d1 := []byte(clusterIssuerYaml)
    err := ioutil.WriteFile(clusterIssuerFilename, d1, 0644)
    if err != nil {
        log.Fatal("Could not write cluster issuer file")
    }

    log.Info("Installing ClusterIssuer resource")
    util.RunCmd(applyClusterIssuerResourceCmd, true)
}


