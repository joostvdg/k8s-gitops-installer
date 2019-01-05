package core

import (
    "fmt"
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "os/exec"
    "strings"
)

func Install(config CoreModernConfig) error {
    folderName := fmt.Sprintf("./cloudbees-core_%s_%s", config.Version, config.Platform)
    coreYaml := fmt.Sprintf("%s/cloudbees-core.yml", folderName)
    log.Info("Install CloudBees Core")
    // install cb core
    installCBCoreCmd := exec.Command("kubectl",
        "apply", "-f", coreYaml,
    )
    util.RunCmd(installCBCoreCmd, true)

    log.Info("Wait for Rollout to succeed")
    watchCBCoreRolloutCmd := exec.Command("kubectl",
        "rollout", "status", "sts", "cjoc",
    )
    util.RunCmd(watchCBCoreRolloutCmd, false)
    return nil
}

func PreInstallConfigure(config CoreModernConfig) {
    folderName := fmt.Sprintf("cloudbees-core_%s_%s", config.Version, config.Platform)
    coreYaml := fmt.Sprintf("%s/cloudbees-core.yml", folderName)
    log.Info("Configuring CloudBees Core (CBCore)")
    if config.Verbose {
        log.Infof("\tCBCore Version: %s\n", config.Version)
        log.Infof("\tKubernetes Resource File: %s\n", coreYaml)
        log.Infof("\tProduction: %s\n", config.Production)
        log.Infof("\tPlatform: %s\n", config.Platform)
        log.Infof("\tSSL: %s\n", config.SSL)
        log.Infof("\tStorage Class Operations Center: %s\n", config.StorageClassOC)
        log.Infof("\tStorage Class Managed Masters: %s\n", config.StorageClassMM)
    }
    validateFileExistsCmd := exec.Command("ls",
        "-lath", coreYaml,
    )
    fileOutput := util.RunCmd(validateFileExistsCmd, false)
    if strings.Contains(fileOutput, "No such file or directory") {
        log.Fatal("Could not find configuration yaml!")
    }

    // ALPINE: sed -i  -e s/test me/test you/g test/a.txt
    // DARWIN: sed -i '' "s/test you/test it/g" build/a.txt
    checkDistributionCmd := exec.Command("uname")
    uname := util.RunCmd(checkDistributionCmd, false)
    darwinHack := "''"
    if uname != "Darwin" {
        darwinHack = "-e"
    }
    // ssl config
    //  - nginx config
    //  - cert-manager annotation
    //  - host names
    if config.SSL {
        replacementString := fmt.Sprintf("s/cje.example.com/%s/g", config.Domain)
        alterTlsSecretReplacementString := fmt.Sprintf("s/#  secretName: %s-tls/  secretName: %s-tls/g", config.Domain, config.Domain)
        alterHostReplacementString := fmt.Sprintf("s/#  - %s/  - %s/g", config.Domain, config.Domain)
        addCertManager1ReplacementString := fmt.Sprintf("s/# Uncomment the next line if NGINX Ingress Controller is configured to do SSL offloading at load balancer level/certmanager.k8s.io\\/cluster-issuer: letsencrypt-prod/g")
        addCertManager2ReplacementString := fmt.Sprintf("s/# \"413 Request Entity Too Large\" uploading plugins, increase client_max_body_size/certmanager.k8s.io\\/acme-challenge-type: http01/g")

        alterDomainNameCmd := exec.Command("sed",
            "-i", darwinHack, replacementString, coreYaml,
        )
        alterTlsCmd := exec.Command("sed",
            "-i", darwinHack, "s/#tls:/tls:/g", coreYaml,
        )
        alterHostsCmd := exec.Command("sed",
            "-i", darwinHack, "s/#- hosts:/- hosts:/g", coreYaml,
        )
        alterHostCmd := exec.Command("sed",
            "-i", darwinHack, alterHostReplacementString, coreYaml,
        )
        alterTlsSecretCmd := exec.Command("sed",
            "-i", darwinHack, alterTlsSecretReplacementString, coreYaml,
        )
        addCertManagerCmd1 := exec.Command("sed",
            "-i", darwinHack, addCertManager1ReplacementString, coreYaml,
        )
        addCertManagerCmd2 := exec.Command("sed",
            "-i", darwinHack, addCertManager2ReplacementString, coreYaml,
        )

        if config.Verbose {
            log.Infof("Updating file %s\n", coreYaml)
            log.Infof("\tDomain Replacement Command\t=%s\n", replacementString)
            log.Infof("\tAlter TLS Secret Command\t=%s\n", alterTlsSecretReplacementString)
            log.Infof("\tAlter Host Command\t\t=%s\n", alterHostReplacementString)
        }

        log.Info("Changing Kubernetes Resource file for SSL configuration")
        if config.Verbose {  log.Info(" - domain name")}
        util.RunCmd(alterDomainNameCmd, true)
        if config.Verbose {  log.Info(" - tls")}
        util.RunCmd(alterTlsCmd, true)
        if config.Verbose {  log.Info(" - hosts")}
        util.RunCmd(alterHostsCmd, true)
        if config.Verbose {  log.Info(" - host")}
        util.RunCmd(alterHostCmd, true)
        if config.Verbose {  log.Info(" - tls secret")}
        util.RunCmd(alterTlsSecretCmd, true)
        if config.Verbose {  log.Info(" - certman 1")}
        util.RunCmd(addCertManagerCmd1, true)
        if config.Verbose {  log.Info(" - certman 2")}
        util.RunCmd(addCertManagerCmd2, true)

        // create cert
        InstallCertificate(config)
    } else {
        // TODO: implement non-ssl
        // sed -e s,https://$DOMAIN_NAME,http://$DOMAIN_NAME,g < cloudbees-core.yml > tmp && mv tmp cloudbees-core.yml
        // sed -e s,ssl-redirect:\ \"true\",ssl-redirect:\ \"false\",g < cloudbees-core.yml > tmp && mv tmp cloudbees-core.yml
    }

    log.Info("Updating Storage Class configuration in Kubernetes resource file")
    // change storage class
    //  - cjoc
    //  - masters
    //sed '/CLIENTSCRIPT="foo"/s/.*/&\
    //CLIENTSCRIPT2="hello"/' file
    //-Dcom.cloudbees.masterprovisioning.kubernetes.KubernetesMasterProvisioning.storageClassName=
    // indent needs to be at position twelve
    storageClassPropertyMM := "-Dcom.cloudbees.masterprovisioning.kubernetes.KubernetesMasterProvisioning.storageClassName"
    storageClassPropertyMMSedBase := `/value: >-/s/.*/&\
            %s=%s/`
    storageClassPropertyMMSed := fmt.Sprintf(storageClassPropertyMMSedBase, storageClassPropertyMM, config.StorageClassMM)
    storageClassMMCmd := exec.Command("sed",
        "-i", darwinHack, storageClassPropertyMMSed, coreYaml,
    )
    util.RunCmd(storageClassMMCmd, true)

    storageClassCjocDefOriginial := "# storageClassName: some-storage-class"
    storageClassCjocDefDesired := fmt.Sprintf("storageClassName: %s", config.StorageClassOC)
    storageClassCjocSed := fmt.Sprintf("s/%s/%s/g", storageClassCjocDefOriginial, storageClassCjocDefDesired)
    storageClassCjocCmd := exec.Command("sed",
        "-i", darwinHack, storageClassCjocSed, coreYaml,
    )
    util.RunCmd(storageClassCjocCmd, true)
}

func InstallCertificate(config CoreModernConfig) {
    var resourceName string
    var resourceNamespace string
    var secretName string
    var domain string
    var issuerName string
    certificateFilename := "cjoc-cert.yml"
    applyCertificateResourceCmd := exec.Command("kubectl",
        "apply", "-f", certificateFilename,
    )
    log.Info("Install Certificate resource for Operations Center")

    resourceNamespace = config.Namespace
    domain = config.Domain
    resourceName = strings.Replace(config.Domain, ".", "-", -1)
    secretName = resourceName+"-tls"
    if config.Production {
        issuerName = "letsencrypt-prod"
    } else {
        resourceName = resourceName + "-stg"
        secretName = secretName + "-stg"
        issuerName = "letsencrypt-staging"
    }

    // name, namespace, secretName, dns, dns, clusterIssueName
    certificate := fmt.Sprintf(Certificate, resourceName, resourceNamespace, secretName, domain, domain, issuerName)
    d1 := []byte(certificate)
    err := ioutil.WriteFile(certificateFilename, d1, 0644)
    if err != nil {
        log.Fatal("Could not write cluster issuer file")
    }

    log.Info("Installing CJOC Certificate resource")
    util.RunCmd(applyCertificateResourceCmd, true)
}
