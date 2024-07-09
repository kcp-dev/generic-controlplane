package tokengetter

import (
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
// uses the specified client to retrieve service accounts, pods, secrets and nodes.
// The client should NOT authenticate using a service account token
// the returned getter will be used to retrieve, or recursion will result.
func NewGetterFromClient(secretLister v1listers.SecretLister, serviceAccountLister v1listers.ServiceAccountLister) serviceaccount.ServiceAccountTokenGetter {
	return clientGetter{secretLister, serviceAccountLister}
}

func (c clientGetter) GetServiceAccount(namespace, name string) (*v1.ServiceAccount, error) {
	return c.serviceAccountLister.ServiceAccounts(namespace).Get(name)
}

func (c clientGetter) GetPod(namespace, name string) (*v1.Pod, error) {
	return nil, apierrors.NewNotFound(v1.Resource("pods"), name)
}

func (c clientGetter) GetSecret(namespace, name string) (*v1.Secret, error) {
	return c.secretLister.Secrets(namespace).Get(name)
}

func (c clientGetter) GetNode(name string) (*v1.Node, error) {
	return nil, apierrors.NewNotFound(v1.Resource("nodes"), name)
}
