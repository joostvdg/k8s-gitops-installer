package weavenet

import (
    "encoding/base64"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "os/exec"
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
    "strings"
)

// TODO: randomly generate password

func Install() {
    log.Info("Installing Weavenet with encryption")

    networkPassword := "vjStsrzC4q7xDnb1wZkYacnk"
    weaveSecretFileName := "weave-passwd"
    weaveSecretName := "weave-passwd"

    checkSecretExistenceCmd := exec.Command("kubectl",
        "get",  "secret",  "-n",  "kube-system",  "weave-passwd",  "-o",  "name",
    )
    secretCmd := exec.Command("kubectl",
        "create", "secret",
        "-n", "kube-system",
        "generic", weaveSecretName,
        "--from-file="+weaveSecretFileName,
    )
    kubernetesVersionCmd := exec.Command("kubectl","version" )

    secretExist := util.RunCmdNonFatal(checkSecretExistenceCmd)

    if secretExist {
        log.Info("Weavenet encryption secret already exists")
    } else {
        d1 := []byte(networkPassword)
        err := ioutil.WriteFile(weaveSecretFileName, d1, 0644)
        if err != nil {
            log.Fatal("Could not write weavenet secret file")
        }
        log.Info("Creating encryption secret")
        util.RunCmd(secretCmd, true)
    }

    kubernetesVersion := util.RunCmd(kubernetesVersionCmd, true)
    encoded := base64.StdEncoding.EncodeToString([]byte(kubernetesVersion))
    encodedVersion := strings.Replace(encoded, "\n", "", -1)

    createCmd := exec.Command("kubectl",
        "apply", "-f",
        "\"https://cloud.weave.works/k8s/net?k8s-version="+encodedVersion+"&password-secret=weave-passwd&env.IPALLOC_RANGE=10.10.0.0/24\"",
    )
    log.Info("Installing Weavenet")
    util.RunCmd(createCmd, true)

}
