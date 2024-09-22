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

package server

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kcp-dev/kcp/cli/pkg/help"
	"github.com/kcp-dev/kcp/pkg/embeddedetcd"
	"github.com/spf13/cobra"

	apiextensionapiserver "k8s.io/apiextensions-apiserver/pkg/apiserver"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	_ "k8s.io/apiserver/pkg/admission" // for admission plugins
	genericapifilters "k8s.io/apiserver/pkg/endpoints/filters"
	genericapiserver "k8s.io/apiserver/pkg/server"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/apiserver/pkg/util/notfoundhandler"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/cli/globalflag"
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
	_ "k8s.io/component-base/metrics/prometheus/workqueue" // for workqueue metrics
	"k8s.io/component-base/term"
	"k8s.io/component-base/version"
	"k8s.io/component-base/version/verflag"
	"k8s.io/klog/v2"
	aggregatorapiserver "k8s.io/kube-aggregator/pkg/apiserver"
	controlplaneapiserver "k8s.io/kubernetes/pkg/controlplane/apiserver"
	_ "k8s.io/kubernetes/pkg/features" // add the kubernetes feature gates

	"github.com/kcp-dev/generic-controlplane/server/batteries"
	options "github.com/kcp-dev/generic-controlplane/server/cmd/options"
	"github.com/kcp-dev/generic-controlplane/server/readiness"
)

// Order for settings:
// Options -> CompletedOptions -> Config -> CompletedConfig -> Server -> Prepared -> Run

func init() {
	utilruntime.Must(logsapi.AddFeatureGates(utilfeature.DefaultMutableFeatureGate))
}

// NewCommand creates a *cobra.Command object with default parameters
func NewCommand() *cobra.Command {
	// manually extract root directory from flags first as it influence all other flags
	rootDir := ".gcp"
	for i, f := range os.Args {
		if f == "--root-directory" {
			if i < len(os.Args)-1 {
				rootDir = os.Args[i+1]
			} // else let normal flag processing fail
		} else if strings.HasPrefix(f, "--root-directory=") {
			rootDir = strings.TrimPrefix(f, "--root-directory=")
		}
	}

	s := options.NewOptions(rootDir)

	cmdStart := &cobra.Command{
		Use: "start",
		Long: `The generic controlplane is a generic controlplane server,
a system serving APIs like Kubernetes, but without the container domain specific
APIs.`,

		// stop printing usage when the command errors
		SilenceUsage: true,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			// silence client-go warnings.
			// kube-apiserver loopback clients should not log self-issued warnings.
			rest.SetDefaultWarningHandler(rest.NoWarnings{})
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested()
			fs := cmd.Flags()

			// Activate logging as soon as possible, after that
			// show flags with the final logging configuration.
			if err := logsapi.ValidateAndApply(s.GenericControlPlane.Logs, utilfeature.DefaultFeatureGate); err != nil {
				return err
			}
			cliflag.PrintFlags(fs)

			completedOptions, err := s.Complete()
			if err != nil {
				return err
			}

			if errs := completedOptions.Validate(); len(errs) != 0 {
				return kerrors.NewAggregate(errs)
			}

			// add feature enablement metrics
			utilfeature.DefaultMutableFeatureGate.AddMetrics()
			ctx := genericapiserver.SetupSignalContext()

			return Run(ctx, *completedOptions)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	var namedFlagSets cliflag.NamedFlagSets
	s.AddFlags(&namedFlagSets)
	verflag.AddFlags(namedFlagSets.FlagSet("global"))
	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmdStart.Name(), logs.SkipLoggingConfigurationFlags())

	fs := cmdStart.Flags()
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	cols, _, _ := term.TerminalSize(cmdStart.OutOrStdout())
	cliflag.SetUsageAndHelpFunc(cmdStart, namedFlagSets, cols)

	startOptionsCmd := &cobra.Command{
		Use:   "options",
		Short: "Show all start command options",
		Long: help.Doc(`
			Show all start command options

			"gcp start"" has a large number of options. This command shows all of them.
		`),
		PersistentPreRunE: func(*cobra.Command, []string) error {
			// silence client-go warnings.
			// apiserver loopback clients should not log self-issued warnings.
			rest.SetDefaultWarningHandler(rest.NoWarnings{})
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmdStart.UseLine())
			cliflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)
			return nil
		},
	}
	cmdStart.AddCommand(startOptionsCmd)

	setPartialUsageAndHelpFunc(cmdStart, namedFlagSets, cols, []string{
		"etcd-servers",
	})

	help.FitTerminal(cmdStart.OutOrStdout())

	return cmdStart
}

// Run runs the specified APIServer. This should never exit.
func Run(ctx context.Context, opts options.CompletedOptions) error {
	// To help debugging, immediately log version
	klog.Infof("Version: %+v", version.Get())

	klog.InfoS("Golang settings", "GOGC", os.Getenv("GOGC"), "GOMAXPROCS", os.Getenv("GOMAXPROCS"), "GOTRACEBACK", os.Getenv("GOTRACEBACK"))

	config, err := options.NewConfig(opts)
	if err != nil {
		return err
	}
	completed, err := config.Complete()
	if err != nil {
		return err
	}

	// the etcd server must be up before NewServer because storage decorators access it right away
	if completed.EmbeddedEtcd.Config != nil {
		klog.Info("Starting embedded etcd server")
		if err := embeddedetcd.NewServer(completed.EmbeddedEtcd).Run(ctx); err != nil {
			return err
		}
	}

	server, err := createServerChain(completed)
	if err != nil {
		return err
	}

	prepared, err := server.PrepareRun()
	if err != nil {
		return err
	}

	// write the kubeconfig file as close to the start of the server as possible
	err = completed.Options.AdminAuthentication.WriteKubeConfig(completed.ControlPlane.Generic, completed.ExtraConfig.GcpAdminToken, completed.ExtraConfig.UserToken)
	if err != nil {
		return err
	}

	// Run the server and wait for readiness

	go func() {
		if err := prepared.RunWithContext(ctx); err != nil {
			klog.Fatal(err, "Failed to run server")
		}
	}()

	// wait for the server to be ready
	klog.Info("Waiting for control plane to be ready")
	err = readiness.WaitForReady(ctx, completed.Options.AdminAuthentication.KubeConfigPath)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

// createServerChain creates the apiservers connected via delegation.
func createServerChain(config options.CompletedConfig) (aggregatorapiserver.APIAggregator, error) {
	// 1. Basic not found handler
	notFoundHandler := notfoundhandler.New(config.ControlPlane.Generic.Serializer, genericapifilters.NoMuxAndDiscoveryIncompleteKey)

	// TODO: we can use single variable here with cleaner logic bellow.
	var aggregatorServer aggregatorapiserver.APIAggregator
	var miniAggregatorServer aggregatorapiserver.APIAggregator

	var apiExtensionsServer *apiextensionapiserver.CustomResourceDefinitions
	var nativeAPIs *controlplaneapiserver.Server
	var err error

	if config.Batteries.IsEnabled(batteries.BatteryCRDs) {
		// Base of CRDs are extension server
		apiExtensionsServer, err = config.APIExtensions.New(genericapiserver.NewEmptyDelegateWithCustomHandler(notFoundHandler))
		if err != nil {
			return nil, fmt.Errorf("failed to create apiextensions-apiserver: %w", err)
		}

		nativeAPIs, err = config.ControlPlane.New("generic-controlplane", apiExtensionsServer.GenericAPIServer)
		if err != nil {
			return nil, fmt.Errorf("failed to create generic controlplane apiserver: %w", err)
		}
	} else {
		// 2. Natively implemented resources
		var err error
		nativeAPIs, err = config.ControlPlane.New("generic-controlplane", genericapiserver.NewEmptyDelegateWithCustomHandler(notFoundHandler))
		if err != nil {
			return nil, fmt.Errorf("failed to create generic controlplane apiserver: %w", err)
		}
	}

	client, err := kubernetes.NewForConfig(config.ControlPlane.Generic.LoopbackClientConfig)
	if err != nil {
		return nil, err
	}
	storageProviders, err := config.ControlPlane.GenericStorageProviders(client.Discovery())
	if err != nil {
		return nil, fmt.Errorf("failed to create storage providers: %w", err)
	}

	// Filter out the disabled batteries
	storageProviders = config.Batteries.FilterStorageProviders(storageProviders)

	if err := nativeAPIs.InstallAPIs(storageProviders...); err != nil {
		return nil, fmt.Errorf("failed to install APIs: %w", err)
	}
	for _, storageProvider := range storageProviders {
		klog.Infof("Serving %s", storageProvider.GroupName())
	}

	// 3. Aggregator for APIServices, discovery and OpenAPI
	// If CRDs are enabled, we wire in, else - its a no-op.
	if config.Batteries.IsEnabled(batteries.BatteryCRDs) {
		aggregatorServer, err = controlplaneapiserver.CreateAggregatorServer(
			config.Aggregator,
			nativeAPIs.GenericAPIServer,
			apiExtensionsServer.Informers.Apiextensions().V1().CustomResourceDefinitions(),
			false,
			controlplaneapiserver.DefaultGenericAPIServicePriorities())
		if err != nil {
			return nil, fmt.Errorf("failed to create kube-aggregator: %w", err)
		}
		miniAggregatorServer, err = config.MiniAggregator.New(nativeAPIs.GenericAPIServer, nativeAPIs, apiExtensionsServer)
		if err != nil {
			return nil, fmt.Errorf("failed to create mini-aggregator: %w", err)
		}

	} else {
		klog.Info("CRDs are disabled, skipping api-extension server")
		aggregatorServer, err = controlplaneapiserver.CreateAggregatorServer(config.Aggregator, nativeAPIs.GenericAPIServer, nil, false, controlplaneapiserver.DefaultGenericAPIServicePriorities())
		if err != nil {
			return nil, fmt.Errorf("failed to create kube-aggregator: %w", err)
		}
		miniAggregatorServer, err = config.MiniAggregator.New(nativeAPIs.GenericAPIServer, nativeAPIs, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create mini-aggregator: %w", err)
		}
	}
	if err != nil {
		// we don't need special handling for innerStopCh because the aggregator server doesn't create any go routines
		return nil, fmt.Errorf("failed to create kube-aggregator: %w", err)
	}

	// 4. Based on if we need APIServices or not return right server
	if config.Batteries.IsEnabled(batteries.BatteryAPIServices) {
		klog.Info("Using aggregator server")
		return aggregatorServer, nil
	} else {
		klog.Info("Using mini-aggregator server")
		return miniAggregatorServer, nil
	}
}
