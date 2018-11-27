package util

import (
    "bytes"
    log "github.com/sirupsen/logrus"
    "os/exec"
)

func RunCmd(cmd *exec.Cmd, failOnError bool) string {
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
    if outStr != "" {
        log.Infof("\n"+outStr+"\n")
    }
    if errStr != "" {
        log.Warnf("\n"+errStr+"\n")
    }
    if err != nil  && failOnError {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
    return outStr
}

func RunCmdNonFatal(cmd *exec.Cmd) bool {
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
    log.Infof("out:\n%s\nerr:\n%s\n", outStr, errStr)
    if err != nil {
        return false
    }
    return true
}

