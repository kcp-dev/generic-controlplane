package main

import (
	"os"

	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/logs/json/register"
	_ "k8s.io/component-base/metrics/prometheus/clientgo"
	_ "k8s.io/component-base/metrics/prometheus/version"

	server "github.com/kcp-dev/generic-controlplane/server/cmd"
	"github.com/kcp-dev/kcp/cli/pkg/help"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "gcp",
		Short: "Generic Control Plane (GCP)",
		Long: help.Doc(`
			GCP is a generic control plane server, a system serving APIs like Kubernetes, but without the container domain specific APIs.
		`),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	command := server.NewCommand()
	cmd.AddCommand(command)

	code := cli.Run(cmd)
	os.Exit(code)
}
