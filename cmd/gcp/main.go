/*
Copyright 2024 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"os"

	server "github.com/kcp-dev/generic-controlplane/server/cmd"
	"github.com/spf13/cobra"

	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/logs/json/register"
	_ "k8s.io/component-base/metrics/prometheus/clientgo"
	_ "k8s.io/component-base/metrics/prometheus/version"

	"github.com/kcp-dev/kcp/cli/pkg/help"
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
