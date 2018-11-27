package gke

import (
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "os/exec"
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
)

func InstallSsdSC() {
    storageClassFilename := "sc-ssd.yml"
    applyStorageClassResourceCmd := exec.Command("kubectl",
        "apply", "-f", storageClassFilename,
    )

    d1 := []byte(SSD_Storage_Class)
    err := ioutil.WriteFile(storageClassFilename, d1, 0644)
    if err != nil {
        log.Fatal("Could not write cluster issuer file")
    }

    log.Info("Installing SSD Storage Class resource")
    util.RunCmd(applyStorageClassResourceCmd, true)
}
