/*
Copyright 2022 The KCP Authors.

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
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/pflag"

	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/group"
	"k8s.io/apiserver/pkg/authentication/request/bearertoken"
	authenticatorunion "k8s.io/apiserver/pkg/authentication/request/union"
	"k8s.io/apiserver/pkg/authentication/user"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	// A kcp admin being member of system-admin
	gcpAdminUserName = "system-admin"
	// A non-admin user part of the "user" battery.
	gcpUserUserName = "user"
)

type AdminAuthentication struct {
	KubeConfigPath string

	// TODO: move into Secret in-cluster, maybe by using an "in-cluster" string as value
	ShardAdminTokenHashFilePath string
}

func NewAdminAuthentication(rootDir string) *AdminAuthentication {
	return &AdminAuthentication{
		KubeConfigPath:              filepath.Join(rootDir, "admin.kubeconfig"),
		ShardAdminTokenHashFilePath: filepath.Join(rootDir, ".admin-token-store"),
	}
}

func (s *AdminAuthentication) Validate() []error {
	if s == nil {
		return nil
	}

	errs := []error{}

	if s.ShardAdminTokenHashFilePath == "" && s.KubeConfigPath != "" {
		errs = append(errs, fmt.Errorf("--admin-kubeconfig requires --admin-token-hash-file-path"))
	}

	return errs
}

func (s *AdminAuthentication) AddFlags(fs *pflag.FlagSet) {
	if s == nil {
		return
	}

	fs.StringVar(&s.KubeConfigPath, "kubeconfig-path", s.KubeConfigPath,
		"Path to which the administrative kubeconfig should be written at startup. If this is relative, it is relative to --root-directory.")
	fs.StringVar(&s.ShardAdminTokenHashFilePath, "authentication-admin-token-path", s.ShardAdminTokenHashFilePath,
		"Path to which the administrative token hash should be written at startup. If this is relative, it is relative to --root-directory.")
}

// ApplyTo returns a new volatile gcp admin token.
// It also returns a new shard admin token and its hash if the configured shard admin hash file is not present.
// If the shard admin hash file is present only the shard admin hash is returned and the returned shard admin token is empty.
func (s *AdminAuthentication) ApplyTo(config *genericapiserver.Config) (volatileGcpAdminToken, volatileUserToken string, err error) {
	volatileUserToken = uuid.New().String()
	volatileGcpAdminToken = uuid.New().String()

	gcpAdminUser := &user.DefaultInfo{
		Name: gcpAdminUserName,
		UID:  uuid.New().String(),
		Groups: []string{
			"system:masters",
		},
	}

	nonAdminUser := &user.DefaultInfo{
		Name:   gcpUserUserName,
		UID:    uuid.New().String(),
		Groups: []string{},
	}

	newAuthenticator := group.NewAuthenticatedGroupAdder(bearertoken.New(authenticator.WrapAudienceAgnosticToken(config.Authentication.APIAudiences, authenticator.TokenFunc(func(ctx context.Context, requestToken string) (*authenticator.Response, bool, error) {
		if requestToken == volatileGcpAdminToken {
			return &authenticator.Response{User: gcpAdminUser}, true, nil
		}

		if requestToken == volatileUserToken {
			return &authenticator.Response{User: nonAdminUser}, true, nil
		}

		return nil, false, nil
	}))))

	config.Authentication.Authenticator = authenticatorunion.New(newAuthenticator, config.Authentication.Authenticator)

	return volatileGcpAdminToken, volatileUserToken, nil
}

func (s *AdminAuthentication) WriteKubeConfig(config genericapiserver.CompletedConfig, gcpAdminToken, userToken string) error {
	externalCACert, _ := config.SecureServing.Cert.CurrentCertKeyContent()
	externalKubeConfigHost := fmt.Sprintf("https://%s", config.ExternalAddress)

	externalKubeConfig := createKubeConfig(gcpAdminToken, userToken, externalKubeConfigHost, "", externalCACert)
	return clientcmd.WriteToFile(*externalKubeConfig, s.KubeConfigPath)
}

func createKubeConfig(gcpAdminToken, userToken, baseHost, tlsServerName string, caData []byte) *clientcmdapi.Config {
	var kubeConfig clientcmdapi.Config
	// Create Client and Shared
	kubeConfig.AuthInfos = map[string]*clientcmdapi.AuthInfo{
		gcpAdminUserName: {Token: gcpAdminToken},
	}
	kubeConfig.Clusters = map[string]*clientcmdapi.Cluster{
		"root": {
			Server:                   baseHost,
			CertificateAuthorityData: caData,
			TLSServerName:            tlsServerName,
		},
	}
	kubeConfig.Contexts = map[string]*clientcmdapi.Context{
		"root": {Cluster: "root", AuthInfo: gcpAdminUserName},
	}
	kubeConfig.CurrentContext = "root"

	if len(userToken) > 0 {
		kubeConfig.AuthInfos[gcpUserUserName] = &clientcmdapi.AuthInfo{Token: userToken}
		kubeConfig.Contexts[gcpUserUserName] = &clientcmdapi.Context{Cluster: "root", AuthInfo: gcpUserUserName}
	}

	return &kubeConfig
}
