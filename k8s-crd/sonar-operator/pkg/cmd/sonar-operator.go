package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	goflag "flag"

	"github.com/golang/glog"
	"github.com/openebs/elves/k8s-crd/sonar-operator/pkg/operator"
	"github.com/spf13/cobra"
)

var (
	cmdName = "sonar-operator"
	usage   = fmt.Sprintf("%s", cmdName)
)

// Define a type for the options of SonarOperator
type SonarOperatorOptions struct {
	KubeConfig     string
	Namespace      string
}

func AddKubeConfigFlag(cmd *cobra.Command, value *string) {
	cmd.Flags().StringVarP(value, "kubeconfig", "", *value, "Path to a kube config. Only required if out-of-cluster.")
}

func AddNamespaceFlag(cmd *cobra.Command, value *string) {
	cmd.Flags().StringVarP(value, "namespace", "n", *value, "Namespace to deploy in. If no namespace is provided, POD_NAMESPACE env. var is used. Lastly, the 'default' namespace will be used as a last option.")
}

// Fatal prints the message (if provided) and then exits. If V(2) or greater,
// glog.Fatal is invoked for extended information.
func fatal(msg string) {
	if glog.V(2) {
		glog.FatalDepth(2, msg)
	}
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(1)
}

// NewCmdOptions creates an options Cobra command to return usage
func NewCmdOptions() *cobra.Command {
	cmd := &cobra.Command{
		Use: "options",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	return cmd
}

// Create a new command for the sonar-operator. This cmd includes logging,
// cmd option parsing from flags, and the customization of the Tectonic assets.
func NewCmdSonarOperator() (*cobra.Command, error) {
	// Define the options for SonarOperator command
	options := SonarOperatorOptions{}

	// Create a new command
	cmd := &cobra.Command{
		Use:   usage,
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			checkErr(Run(cmd, &options), fatal)
		},
	}

	// Bind & parse flags defined by external projects.
	// e.g. This imports the golang/glog pkg flags into the cmd flagset
	cmd.Flags().AddGoFlagSet(goflag.CommandLine)
	goflag.CommandLine.Parse([]string{})

	// Define the flags allowed in this command & store each option provided
	// as a flag, into the SonarOperatorOptions
	AddKubeConfigFlag(cmd, &options.KubeConfig)
	AddNamespaceFlag(cmd, &options.Namespace)

	return cmd, nil
}

// Run the customization of the Tectonic assets
func Run(cmd *cobra.Command, options *SonarOperatorOptions) error {
	cntlr, err := operator.NewSubmarineController(
		options.KubeConfig, options.Namespace)

	if err != nil {
		return err
	}

	// Relay OS signals to the chan
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	stop := make(chan struct{})
	go cntlr.Start(stop)

	// Block until signaled to stop
	<-signals

	// Close the stop chan / shutdown the controller
	close(stop)
	glog.Infof("Shutting down SonarOperator...")
	return nil
}

func checkErr(err error, handleErr func(string)) {
	if err == nil {
		return
	}
	handleErr(err.Error())
}
