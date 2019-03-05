package core

import (
    "fmt"
    "strings"

    "io"
    "net/http"
    "os"
    "os/exec"
    "github.com/joostvdg/k8s-gitops-installer/pkg/util"
    log "github.com/sirupsen/logrus"
)

// https://downloads.cloudbees.com/cloudbees-core/cloud/2.138.2.2/cloudbees-core_2.138.2.2_kubernetes.tgz
var (
    Filename = `cloudbees-core_%s_%s.tgz`
    DownloadURLProto = `https://downloads.cloudbees.com/cloudbees-core/cloud/%s/`
)

func DownloadAndUnpack(config CoreModernConfig) {
    log.Infof("Downloading CloudBees Core version %s for platform %s\n", config.Version, config.Platform)
    baseUrl := fmt.Sprintf(DownloadURLProto, config.Version)
    filename := fmt.Sprintf(Filename, config.Version, config.Platform)
    shaFilename := filename+".sha256"
    url := baseUrl+filename
    shaUrl := baseUrl+shaFilename

    checkDistributionCmd := exec.Command("uname")
    uname := util.RunCmd(checkDistributionCmd, false)
    uname = strings.TrimSpace(uname)
    if config.Verbose {
        log.Infof("Found uname: [%s]", uname)
    }
    var validateDownloadCmd *exec.Cmd
    if uname == "Darwin" {
        // Darwin = macOs
        validateDownloadCmd = exec.Command("shasum",
            "-a", "256", shaFilename,
        )
    } else {
        // sha256sum -c $INSTALLER.sha256
        validateDownloadCmd = exec.Command("sha256sum",
            "-c", shaFilename,
        )
    }

    unpackDownloadCmd := exec.Command("tar",
        "xzvf", filename,
    )

    if config.Verbose {
        log.Infof("\tFilename: %s\n", filename)
        log.Infof("\tFilename SHA: %s\n", shaFilename)
        log.Infof("\tURL: %s\n", url)
        log.Infof("\tURL SHA: %s\n", shaUrl)
    }

    err1 := downloadFile(filename, url)
    if err1 != nil {
        log.Errorf("Could not download %s: %s", url, err1)
    }
    err2 := downloadFile(shaFilename, shaUrl)
    if err2 != nil {
        log.Errorf("Could not download %s: %s", url, err2)
    }


    log.Info("Validating download via sha file")
    util.RunCmd(validateDownloadCmd, true)

    log.Info("Unpacking download")
    util.RunCmd(unpackDownloadCmd, true)

}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
// Originally from: https://golangcode.com/download-a-file-from-a-url/
func downloadFile(filepath string, url string) error {
    log.Infof("Downloading %s from %s\n", filepath, url)
    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}
