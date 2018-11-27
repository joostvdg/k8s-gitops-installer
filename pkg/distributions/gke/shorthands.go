package gke

import (
    "os/exec"
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
    log "github.com/sirupsen/logrus"
)

// TODO: this works for GKE only

func GetNginxIngressIp(name string, namespace string) {
    log.Infof("Get Nginx ingress ip in namespace %s\n", namespace)
    ingressIptCmd := exec.Command("kubectl",
        "get", "svc", "-n", namespace, name, "-o", "jsonpath=\"{.status.loadBalancer.ingress[0].ip}\"",
    )

    util.RunCmd(ingressIptCmd, true)
}
