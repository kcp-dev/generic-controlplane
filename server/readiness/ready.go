package readiness

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

// WaitForReady waits for the control plane to be ready.
func WaitForReady(ctx context.Context, kubeConfigPath string) error {
	// wait for readiness
	logger := klog.FromContext(ctx)
	logger.Info("Waiting for /readyz to succeed")
	lastSeenUnready := sets.New[string]()

	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
		&clientcmd.ConfigOverrides{CurrentContext: "root"},
	)
	config, err := configLoader.ClientConfig()
	if err != nil {
		return err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	for {
		time.Sleep(500 * time.Millisecond)

		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled")
		default:
		}

		res := client.RESTClient().Get().AbsPath("/readyz").Do(ctx)
		if _, err := res.Raw(); err != nil {
			unreadyComponents := unreadyComponentsFromError(err)
			//logger.Error(err, "control plane not ready", "unreadyComponents", sets.List[string](unreadyComponents), "error", err)
			if !lastSeenUnready.Equal(unreadyComponents) {
				logger.Error(err, "control plane not ready", "unreadyComponents", sets.List[string](unreadyComponents), "error", err)
				lastSeenUnready = unreadyComponents
			}
		}

		// When there is an error for invalid certificate, we should exit immediately
		// as there is no point in retrying.
		if res.Error() != nil {
			if strings.Contains(res.Error().Error(), "failed to verify certificate: x509") {
				logger.Error(res.Error(), "control plane not ready")
				logger.Info("This is likely due to certificates folder containing invalid certificates. Please fix them and restart the control plane.")
				return res.Error()
			}
		}

		var rc int
		res.StatusCode(&rc)
		if rc == http.StatusOK {
			logger.Info("Control plane is ready")
			break
		}
		logger.Info("Control plane not ready", "status", rc, "unreadyComponents", sets.List[string](lastSeenUnready))

	}

	return nil
}

// there doesn't seem to be any simple way to get a metav1.Status from the Go client, so we get
// the content in a string-formatted error, unfortunately.
func unreadyComponentsFromError(err error) sets.Set[string] {
	innerErr := strings.TrimPrefix(strings.TrimSuffix(err.Error(), `") has prevented the request from succeeding`), `an error on the server ("`)
	unreadyComponents := sets.New[string]()
	for _, line := range strings.Split(innerErr, `\n`) {
		if name := strings.TrimPrefix(strings.TrimSuffix(line, ` failed: reason withheld`), `[-]`); name != line {
			// NB: sometimes the error we get is truncated (server-side?) to something like: `\n[-]poststar") has prevented the request from succeeding`
			// In those cases, the `name` here is also truncated, but nothing we can do about that. For that reason, the list of components returned is
			// not durable and should not be parsed.
			unreadyComponents.Insert(name)
		}
	}
	return unreadyComponents
}
