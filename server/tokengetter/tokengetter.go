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

package tokengetter

import (
	"context"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/kubernetes/pkg/serviceaccount"
)

// clientGetter implements ServiceAccountTokenGetter using a factory function
type clientGetter struct {
	secretLister         v1listers.SecretLister
	serviceAccountLister v1listers.ServiceAccountLister
}

// NewGetterFromClient returns a ServiceAccountTokenGetter that
// uses the specified client to retrieve service accounts, secrets and
// return errors for nodes, pods, and webhook configurations.
func NewGetterFromClient(secretLister v1listers.SecretLister, serviceAccountLister v1listers.ServiceAccountLister) serviceaccount.ServiceAccountTokenGetter {
	return clientGetter{secretLister, serviceAccountLister}
}

func (c clientGetter) GetServiceAccount(ctx context.Context, namespace, name string) (*v1.ServiceAccount, error) {
	return c.serviceAccountLister.ServiceAccounts(namespace).Get(name)
}

func (c clientGetter) GetPod(ctx context.Context, _, name string) (*v1.Pod, error) {
	return nil, apierrors.NewNotFound(v1.Resource("pods"), name)
}

func (c clientGetter) GetSecret(ctx context.Context, namespace, name string) (*v1.Secret, error) {
	return c.secretLister.Secrets(namespace).Get(name)
}

func (c clientGetter) GetNode(ctx context.Context, name string) (*v1.Node, error) {
	return nil, apierrors.NewNotFound(v1.Resource("nodes"), name)
}

func (c clientGetter) GetValidatingWebhookConfiguration(ctx context.Context, name string) (*admissionregistrationv1.ValidatingWebhookConfiguration, error) {
	return nil, apierrors.NewNotFound(admissionregistrationv1.Resource("validatingwebhookconfigurations"), name)
}

func (c clientGetter) GetMutatingWebhookConfiguration(ctx context.Context, name string) (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	return nil, apierrors.NewNotFound(admissionregistrationv1.Resource("mutatingwebhookconfigurations"), name)
}
