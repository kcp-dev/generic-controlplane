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
	cryptorand "crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"

	etcdoptions "github.com/kcp-dev/kcp/pkg/embeddedetcd/options"
	"k8s.io/client-go/informers"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/util/keyutil"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	controlplaneapiserveroptions "k8s.io/kubernetes/pkg/controlplane/apiserver/options"
	"k8s.io/kubernetes/pkg/features"
	kubeoptions "k8s.io/kubernetes/pkg/kubeapiserver/options"
	"k8s.io/kubernetes/pkg/serviceaccount"

	"github.com/kcp-dev/generic-controlplane/server/batteries"
	"github.com/kcp-dev/generic-controlplane/server/tokengetter"
)

// Options holds the configuration for the generic controlplane server.
type Options struct {
	GenericControlPlane controlplaneapiserveroptions.Options
	EmbeddedEtcd        etcdoptions.Options

	AdminAuthentication AdminAuthentication

	Extra ExtraOptions
}

// ExtraOptions holds the extra configuration for the generic controlplane server.
type ExtraOptions struct {
	RootDir   string
	Batteries batteries.Batteries
}

type completedOptions struct {
	GenericControlPlane controlplaneapiserveroptions.CompletedOptions
	EmbeddedEtcd        etcdoptions.CompletedOptions

	AdminAuthentication AdminAuthentication

	Extra ExtraOptions
}

// CompletedOptions holds the completed configuration for the generic controlplane server.
type CompletedOptions struct {
	*completedOptions
}

// NewOptions creates a new Options with default parameters.
func NewOptions(rootDir string) *Options {
	o := &Options{
		GenericControlPlane: *controlplaneapiserveroptions.NewOptions(),
		EmbeddedEtcd:        *etcdoptions.NewOptions(rootDir),
		AdminAuthentication: *NewAdminAuthentication(rootDir),
		Extra: ExtraOptions{
			RootDir:   rootDir,
			Batteries: batteries.New(),
		},
	}

	// Disable node related features to prevent the need for informers.
	utilfeature.DefaultMutableFeatureGate.OverrideDefault(features.ServiceAccountTokenNodeBindingValidation, false)
	utilfeature.DefaultMutableFeatureGate.OverrideDefault(features.ServiceAccountTokenNodeBinding, false)

	factory := func(factory informers.SharedInformerFactory) serviceaccount.ServiceAccountTokenGetter {
		return tokengetter.NewGetterFromClient(factory.Core().V1().Secrets().Lister(), factory.Core().V1().ServiceAccounts().Lister())
	}

	// override for standalone mode
	o.GenericControlPlane.SecureServing.ServerCert.CertDirectory = rootDir
	// We use KCP form of the authentication options as it does not contain nodes and pods
	// informers.
	o.GenericControlPlane.Authentication = kubeoptions.NewBuiltInAuthenticationOptions().
		WithAnonymous().
		WithBootstrapToken().
		WithClientCert().
		WithOIDC().
		WithRequestHeader().
		WithServiceAccounts().
		WithTokenFile().
		WithWebHook()

	o.GenericControlPlane.Authentication.ServiceAccounts.OptionalTokenGetter = factory

	o.GenericControlPlane.Authentication.ServiceAccounts.Issuers = []string{"https://gcp.default.svc"}
	o.GenericControlPlane.Etcd.StorageConfig.Transport.ServerList = []string{"embedded"}
	o.GenericControlPlane.Features.EnablePriorityAndFairness = false
	// turn on the watch cache
	o.GenericControlPlane.Etcd.EnableWatchCache = true

	// Flush out the default admission plugins and set the ones we want bellow.
	// TODO: Generic control plane should come with a default set of plugins working out of the box.
	o.GenericControlPlane.Admission.GenericAdmission.Plugins = admission.NewPlugins()

	return o
}

// AddFlags adds flags for a specific APIServer to the specified FlagSet.
func (o *Options) AddFlags(fss *cliflag.NamedFlagSets) {
	o.GenericControlPlane.AddFlags(fss)

	etcdServers := fss.FlagSet("etcd").Lookup("etcd-servers")
	etcdServers.Usage += " By default an embedded etcd server is started."

	o.EmbeddedEtcd.AddFlags(fss.FlagSet("Embedded etcd"))
	o.AdminAuthentication.AddFlags(fss.FlagSet("GCP Standalone Authentication"))

	o.Extra.Batteries.AddFlags(fss.FlagSet("Batteries"))
}

// Complete fills in any fields not set that are required to have valid data.
func (o *Options) Complete() (*CompletedOptions, error) {
	if servers := o.GenericControlPlane.Etcd.StorageConfig.Transport.ServerList; len(servers) == 1 && servers[0] == "embedded" {
		klog.Background().Info("enabling embedded etcd server")
		o.EmbeddedEtcd.Enabled = true
	}

	o.Extra.Batteries.Complete()

	var serviceAccountFile string
	if len(o.GenericControlPlane.Authentication.ServiceAccounts.KeyFiles) == 0 {
		// use sa.key and auto-generate if not existing
		serviceAccountFile = filepath.Join(o.Extra.RootDir, "sa.key")
		if _, err := os.Stat(serviceAccountFile); os.IsNotExist(err) {
			klog.Background().WithValues("file", serviceAccountFile).Info("generating service account key file")
			key, err := rsa.GenerateKey(cryptorand.Reader, 4096)
			if err != nil {
				return nil, fmt.Errorf("error generating service account private key: %w", err)
			}

			encoded, err := keyutil.MarshalPrivateKeyToPEM(key)
			if err != nil {
				return nil, fmt.Errorf("error converting service account private key to PEM format: %w", err)
			}
			if err := keyutil.WriteKey(serviceAccountFile, encoded); err != nil {
				return nil, fmt.Errorf("error writing service account private key file %q: %w", serviceAccountFile, err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("error checking service account key file %q: %w", serviceAccountFile, err)
		}

		// set the key file to generic-controlplane server
		o.GenericControlPlane.Authentication.ServiceAccounts.KeyFiles = []string{serviceAccountFile}

		if o.GenericControlPlane.ServiceAccountSigningKeyFile == "" {
			o.GenericControlPlane.ServiceAccountSigningKeyFile = serviceAccountFile
		}
	}

	// override set of admission plugins
	o.Extra.Batteries.RegisterAllAdmissionPlugins(o.GenericControlPlane.Admission.GenericAdmission.Plugins)
	o.GenericControlPlane.Admission.GenericAdmission.DisablePlugins = sets.List[string](o.Extra.Batteries.DefaultOffAdmissionPlugins())
	o.GenericControlPlane.Admission.GenericAdmission.RecommendedPluginOrder = batteries.AllOrderedPlugins

	var err error
	if !filepath.IsAbs(o.EmbeddedEtcd.Directory) {
		o.EmbeddedEtcd.Directory, err = filepath.Abs(o.EmbeddedEtcd.Directory)
		if err != nil {
			return nil, err
		}
	}
	if !filepath.IsAbs(o.GenericControlPlane.SecureServing.ServerCert.CertDirectory) {
		o.GenericControlPlane.SecureServing.ServerCert.CertDirectory, err = filepath.Abs(o.GenericControlPlane.SecureServing.ServerCert.CertDirectory)
		if err != nil {
			return nil, err
		}
	}
	if !filepath.IsAbs(o.AdminAuthentication.ShardAdminTokenHashFilePath) {
		o.AdminAuthentication.ShardAdminTokenHashFilePath, err = filepath.Abs(o.AdminAuthentication.ShardAdminTokenHashFilePath)
		if err != nil {
			return nil, err
		}
	}
	if !filepath.IsAbs(o.AdminAuthentication.KubeConfigPath) {
		o.AdminAuthentication.KubeConfigPath, err = filepath.Abs(o.AdminAuthentication.KubeConfigPath)
		if err != nil {
			return nil, err
		}
	}

	completedGenericServerRunOptions, err := o.GenericControlPlane.Complete(nil, nil)
	if err != nil {
		return nil, err
	}

	completedEmbeddedEtcd := o.EmbeddedEtcd.Complete(o.GenericControlPlane.Etcd)

	return &CompletedOptions{
		completedOptions: &completedOptions{
			GenericControlPlane: completedGenericServerRunOptions,
			EmbeddedEtcd:        completedEmbeddedEtcd,
			AdminAuthentication: o.AdminAuthentication,
			Extra:               o.Extra,
		},
	}, nil
}

// Validate validates the generic controlplane server options.
func (o *CompletedOptions) Validate() []error {
	var errs []error

	errs = append(errs, o.GenericControlPlane.Validate()...)
	errs = append(errs, o.EmbeddedEtcd.Validate()...)
	errs = append(errs, o.AdminAuthentication.Validate()...)
	errs = append(errs, o.Extra.Batteries.Validate()...)

	return errs
}
