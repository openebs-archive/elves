package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	goflag "flag"

	"github.com/golang/glog"
	"github.com/metral/memhog-operator/pkg/operator"
	"github.com/spf13/cobra"
)

var (
	cmdName = "memhog-operator"
	usage   = fmt.Sprintf("%s", cmdName)
)

// Define a type for the options of MemHogOperator
type MemHogOperatorOptions struct {
	KubeConfig     string
	Namespace      string
	PrometheusAddr string
}

func AddKubeConfigFlag(cmd *cobra.Command, value *string) {
	cmd.Flags().StringVarP(value, "kubeconfig", "", *value, "Path to a kube config. Only required if out-of-cluster.")
}

func AddNamespaceFlag(cmd *cobra.Command, value *string) {
	cmd.Flags().StringVarP(value, "namespace", "n", *value, "Namespace to deploy in. If no namespace is provided, POD_NAMESPACE env. var is used. Lastly, the 'default' namespace will be used as a last option.")
}

func AddPrometheusFlag(cmd *cobra.Command, value *string) {
	cmd.Flags().StringVarP(value, "prometheus-addr", "", *value, "The address & port of the Prometheus service. e.g. http://prometheus.tectonic-system:9090")
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

// Create a new command for the memhog-operator. This cmd includes logging,
// cmd option parsing from flags, and the customization of the Tectonic assets.
func NewCmdMemHogOperator() (*cobra.Command, error) {
	// Define the options for MemHogOperator command
	options := MemHogOperatorOptions{}

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
	// as a flag, into the MemHogOperatorOptions
	AddKubeConfigFlag(cmd, &options.KubeConfig)
	AddNamespaceFlag(cmd, &options.Namespace)
	AddPrometheusFlag(cmd, &options.PrometheusAddr)

	return cmd, nil
}

// Run the customization of the Tectonic assets
func Run(cmd *cobra.Command, options *MemHogOperatorOptions) error {
	cntlr, err := operator.NewAppMonitorController(
		options.KubeConfig, options.Namespace, options.PrometheusAddr)

	if err != nil {
		return err
	}

	// Relay OS signals to the chan
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Example: Create a new AppMontior & instantiate it in the cluster
	// am := operator.NewAppMonitor("my-app-monitor", 80, 2)
	// am.Instantiate(options.KubeConfig, options.Namespace)

	stop := make(chan struct{})
	go cntlr.Start(stop)

	// Block until signaled to stop
	<-signals

	// Close the stop chan / shutdown the controller
	close(stop)
	glog.Infof("Shutting down MemHogOperator...")
	return nil
}

func checkErr(err error, handleErr func(string)) {
	if err == nil {
		return
	}
	handleErr(err.Error())
}
