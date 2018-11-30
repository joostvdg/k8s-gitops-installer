package cmd

import (
	"fmt"
    "github.com/go-ozzo/ozzo-validation/is"
    "github.com/spf13/cobra"
    log "github.com/sirupsen/logrus"
	"os"
	"github.com/joostvdg/k8s-gitops-installer/pkg/weavenet"
    "github.com/joostvdg/k8s-gitops-installer/pkg/helm"
    "github.com/joostvdg/k8s-gitops-installer/pkg/ingress"
    "github.com/joostvdg/k8s-gitops-installer/pkg/validate"
    "github.com/joostvdg/k8s-gitops-installer/pkg/namespace"
    "github.com/joostvdg/k8s-gitops-installer/pkg/letsencrypt"
    "github.com/joostvdg/k8s-gitops-installer/pkg/distributions/gke"
    "github.com/joostvdg/k8s-gitops-installer/pkg/cloudbees/core"
	"github.com/go-ozzo/ozzo-validation"
    "time"
)

var Name string
var Default bool
var Email string
var Production bool
var Platform string
var Version string
var Verbose bool
var Namespace string
var DomainName string

func init() {

    installCmd.AddCommand(installWeavenetCmd)
    installCmd.AddCommand(installHelmCmd)
    installCmd.AddCommand(installNginxCmd)

    validateCmd.AddCommand(validateKubectlCmd)

    installLetsEncryptCmd.Flags().StringVarP(&Email, "email", "e", "", "Email address for Let's Encrypt account")
    installLetsEncryptCmd.Flags().BoolVar(&Production, "production", false, "For using Let's Encrypt Production certificates")
    installCmd.AddCommand(installLetsEncryptCmd)

    createNamespaceCmd.Flags().StringVarP(&Name, "name", "n", "", "The name of the namespace")
    createNamespaceCmd.Flags().BoolVar(&Default, "default", false, "Set this namespace as default in current context")
    createCmd.AddCommand(createNamespaceCmd)

    gkeGetIngressIpCmd.Flags().StringVarP(&Name, "name", "n", "ingress-nginx", "The name of the nginx ingress controller")
    gkeGetIngressIpCmd.Flags().StringVarP(&Namespace, "namespace", "s", "ingress-nginx", "The namespace where the ingress controller resides")

    gkeCmd.AddCommand(gkeGetIngressIpCmd)
    gkeCmd.AddCommand(gkeSsdStorageClassCmd)

    cbcInstallCmd.Flags().StringVarP(&Version, "version", "v", "2.138.3.1", "The version for CloudBees Core, or use 2.138.3.1 by leaving empty")
    cbcInstallCmd.Flags().StringVarP(&Platform, "platform", "p", "kubernetes", "Platform for CloudBees Core [*kubernetes, openshift] (* default)")
    cbcInstallCmd.Flags().StringVarP(&DomainName, "domainName", "d", "cje.example.com", "The domain name to use")
    cbcInstallCmd.Flags().StringVarP(&Namespace, "namespace", "s", "default", "The namespace where the cloudbees core should be installed")
    cbcInstallCmd.Flags().BoolVar(&Verbose, "verbose", false, "Set this if you want more logging")
    cbcInstallCmd.Flags().BoolVar(&Production, "production", false, "For using Let's Encrypt Production certificates")
    cbcCmd.AddCommand(cbcInstallCmd)
    cbcCmd.AddCommand(cbcPasswordCmd)

    rootCmd.AddCommand(validateCmd)
    rootCmd.AddCommand(versionCmd)
    rootCmd.AddCommand(installCmd)
    rootCmd.AddCommand(createCmd)
    rootCmd.AddCommand(gkeCmd)
    rootCmd.AddCommand(cbcCmd)
}

var validateCmd = &cobra.Command{
    Use:   "validate",
    Short: "Will validate sub-resources",
    Long:  `Anything to do with validations`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {},
}

var validateKubectlCmd = &cobra.Command{
    Use:   "kubectl",
    Short: "Will validate kubectl",
    Long:  `Anything to do with validating kubectl`,
    Run: func(cmd *cobra.Command, args []string) {
        validate.Kubectl()
    },
}

var createCmd = &cobra.Command{
    Use:   "create",
    Short: "Will create subresources",
    Long:  `Anything to do with creates`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {},
}

var gkeCmd = &cobra.Command{
    Use:   "gke",
    Short: "Will create gke specific subresources",
    Long:  `Anything to do with gke specific resources`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {},
}

var gkeSsdStorageClassCmd = &cobra.Command{
    Use:   "sc-ssd",
    Short: "Will create gke ssd storage class",
    Long:  `Will create the GKE SSD based Storage Class`,
    Run: func(cmd *cobra.Command, args []string) {
        gke.GetNginxIngressIp(Name, Namespace)
    },
}

var gkeGetIngressIpCmd = &cobra.Command{
    Use:   "ing-ip",
    Short: "Will print ingress ip",
    Long:  `Will print ip of nginx ingress controller service`,
    Run: func(cmd *cobra.Command, args []string) {
        gke.GetNginxIngressIp(Name, Namespace)
    },
}

var createNamespaceCmd = &cobra.Command{
    Use:   "namespace",
    Short: "Will create a namespace",
    Long:  `Anything to do with creating namespaces`,
    Run: func(cmd *cobra.Command, args []string) {
        if Name == "" || len(Name) < 2 {
            log.Fatal("Must have a name with at least two characters")
        }
        namespace.Create(Name)
        if Default {
            namespace.SetDefault(Name)
        }
    },
}

var cbcCmd = &cobra.Command{
    Use:   "cbc",
    Short: "For all CloudBees Core related things",
    Long:  `Anything to do with CloudBees Core`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {},
}

var cbcPasswordCmd = &cobra.Command{
    Use:   "pass",
    Short: "Prints initial cjoc administrator password",
    Long:  `Prints initial cjoc administrator password, doesn't work if there are other users already'`,
    Run: func(cmd *cobra.Command, args []string) {
        core.GetCjocPassword()
    },
}

var cbcInstallCmd = &cobra.Command{
    Use:   "install",
    Short: "Will install CloudBees Core",
    Long:  `Anything to do with installing cloudbees core, download, update yaml's etc'`,
    Run: func(cmd *cobra.Command, args []string) {
        config := core.CoreModernConfig{
            Platform: Platform,
            Version:  Version,
            Verbose:  Verbose,
            StorageClassOC: "ssd",
            StorageClassMM: "ssd",
            SSL: true,
            Domain: DomainName,
            Production: Production,
            Namespace: Namespace,
        }

        // download
        core.DownloadAndUnpack(config)

        // configure
        core.PreInstallConfigure(config)

        // install
        core.Install(config)

        // wait for startup and get print password
        log.Info("Wait 3 minutes for Operations Center to startup")
        time.Sleep(3 * time.Minute)
        core.GetCjocPassword()
    },
}

var installCmd = &cobra.Command{
    Use:   "install",
    Short: "Will install subresources",
    Long:  `Anything to do with installs`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {},
}

var installWeavenetCmd = &cobra.Command{
    Use:   "weavenet",
    Short: "Will install weavenet",
    Long:  `Anything to do with installing weavenet`,
    Run: func(cmd *cobra.Command, args []string) {
        weavenet.Install()
    },
}

var installHelmCmd = &cobra.Command{
    Use:   "helm",
    Short: "Will install helm",
    Long:  `Anything to do with installing helm`,
    Run: func(cmd *cobra.Command, args []string) {
        helm.Install()
    },
}

var installNginxCmd = &cobra.Command{
    Use:   "nginx",
    Short: "Will install nginx controller",
    Long:  `Anything to do with installing nginx as ingress controller`,
    Run: func(cmd *cobra.Command, args []string) {
        ingress.InstallNginx()
    },
}

var installLetsEncryptCmd = &cobra.Command{
    Use:   "letsencrypt",
    Short: "Will install letsencrypt",
    Long:  `Will install letsencrypt through cert-manager`,
    Run: func(cmd *cobra.Command, args []string) {
        err := validation.Errors{ "email": validation.Validate(Email, validation.Required, is.Email),
        }.Filter()
        if err != nil {
            log.Fatalf("There are input validation errors: %s\n", err)
        }

        letsencrypt.InstallWithCertmanager(Email, Production)
    },
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of JPC Go",
	Long:  `All software has versions. This is JPC Go's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("JPC Go 0.1.0")
	},
}

var rootCmd = &cobra.Command{
	Use:   "kgi",
	Short: "kgi is a small cli",
	Long:  `Yada yada yada`,
	Run: func(cmd *cobra.Command, args []string) {
		// return "0.1.0"
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
