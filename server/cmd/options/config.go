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

package options

import (
	"github.com/kcp-dev/kcp/pkg/embeddedetcd"

	apiextensionsapiserver "k8s.io/apiextensions-apiserver/pkg/apiserver"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/util/webhook"
	aggregatorapiserver "k8s.io/kube-aggregator/pkg/apiserver"
	aggregatorscheme "k8s.io/kube-aggregator/pkg/apiserver/scheme"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/controlplane"
	controlplaneapiserver "k8s.io/kubernetes/pkg/controlplane/apiserver"
	generatedopenapi "k8s.io/kubernetes/pkg/generated/openapi"

	"github.com/kcp-dev/generic-controlplane/server/batteries"
)

// Config holds the configuration for the generic controlplane server.
type Config struct {
	Options CompletedOptions

	EmbeddedEtcd  *embeddedetcd.Config
	Aggregator    *aggregatorapiserver.Config
	ControlPlane  *controlplaneapiserver.Config
	APIExtensions *apiextensionsapiserver.Config

	ExtraConfig
}

// ExtraConfig holds the extra configuration for the generic controlplane server.
type ExtraConfig struct {
	// authentication
	GcpAdminToken, UserToken string
	// Batteries holds the batteries configuration for the generic controlplane server.
	Batteries batteries.Batteries
}

type completedConfig struct {
	Options CompletedOptions

	EmbeddedEtcd embeddedetcd.CompletedConfig

	Aggregator    aggregatorapiserver.CompletedConfig
	ControlPlane  controlplaneapiserver.CompletedConfig
	APIExtensions apiextensionsapiserver.CompletedConfig

	ExtraConfig
}

// CompletedConfig holds the completed configuration for the generic controlplane server.
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data.
func (c *Config) Complete() (CompletedConfig, error) {
	c.Batteries.Complete()

	return CompletedConfig{&completedConfig{
		Options: c.Options,

		EmbeddedEtcd:  c.EmbeddedEtcd.Complete(),
		ControlPlane:  c.ControlPlane.Complete(),
		Aggregator:    c.Aggregator.Complete(),
		APIExtensions: c.APIExtensions.Complete(),

		ExtraConfig: c.ExtraConfig,
	}}, nil
}

// NewConfig creates all the self-contained pieces making up an
// generic controlplane server.
func NewConfig(opts CompletedOptions) (*Config, error) {
	c := &Config{
		Options: opts,
		ExtraConfig: ExtraConfig{
			Batteries: opts.Extra.Batteries,
		},
	}

	if opts.EmbeddedEtcd.Enabled {
		var err error
		c.EmbeddedEtcd, err = embeddedetcd.NewConfig(opts.EmbeddedEtcd, opts.GenericControlPlane.Etcd.EnableWatchCache)
		if err != nil {
			return nil, err
		}
	}

	genericConfig, versionedInformers, storageFactory, err := controlplaneapiserver.BuildGenericConfig(
		opts.GenericControlPlane,
		[]*runtime.Scheme{legacyscheme.Scheme, apiextensionsapiserver.Scheme, aggregatorscheme.Scheme},
		controlplane.DefaultAPIResourceConfigSource(),
		generatedopenapi.GetOpenAPIDefinitions,
	)
	if err != nil {
		return nil, err
	}

	// set standalone config
	c.GcpAdminToken, c.UserToken, err = opts.AdminAuthentication.ApplyTo(genericConfig)
	if err != nil {
		return nil, err
	}

	serviceResolver := webhook.NewDefaultServiceResolver()
	kubeAPIs, pluginInitializer, err := controlplaneapiserver.CreateConfig(opts.GenericControlPlane, genericConfig, versionedInformers, storageFactory, serviceResolver, nil)
	if err != nil {
		return nil, err
	}
	c.ControlPlane = kubeAPIs

	authInfoResolver := webhook.NewDefaultAuthenticationInfoResolverWrapper(kubeAPIs.ProxyTransport, kubeAPIs.Generic.EgressSelector, kubeAPIs.Generic.LoopbackClientConfig, kubeAPIs.Generic.TracerProvider)
	apiExtensions, err := controlplaneapiserver.CreateAPIExtensionsConfig(*kubeAPIs.Generic, kubeAPIs.VersionedInformers, pluginInitializer, opts.GenericControlPlane, 3, serviceResolver, authInfoResolver)
	if err != nil {
		return nil, err
	}
	c.APIExtensions = apiExtensions

	aggregator, err := controlplaneapiserver.CreateAggregatorConfig(*kubeAPIs.Generic, opts.GenericControlPlane, kubeAPIs.VersionedInformers, serviceResolver, kubeAPIs.ProxyTransport, kubeAPIs.Extra.PeerProxy, pluginInitializer)
	if err != nil {
		return nil, err
	}
	// IMPORTANT: disable the available condition controller in the aggregator
	// to prevent it to try use Service and Endpoints resources which are not enabled in the generic controlplane.
	aggregator.ExtraConfig.DisableRemoteAvailableConditionController = true
	c.Aggregator = aggregator

	return c, nil
}
