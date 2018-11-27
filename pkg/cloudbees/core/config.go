package core


// CoreModernConfig configuration holder for installing cloudbees core on modern cloud environments
type CoreModernConfig struct {
    // https://downloads.cloudbees.com/cloudbees-core/cloud/2.138.2.2/
    Platform string // kubernetes, openshift
    Version string
    Namespace string
    StorageClassOC string // Storage Class Operations Center
    StorageClassMM string // Storage Class Managed Master
    SSL bool
    Domain string
    Production bool
    Verbose bool
}
