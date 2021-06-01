package k8s

import (
	"golang.org/x/xerrors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func BuildClientset() (*kubernetes.Clientset, error) {
	configLoadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		configLoadingRules,
		&clientcmd.ConfigOverrides{},
	)

	clientConfig, err := config.ClientConfig()
	if err != nil {
		return nil, xerrors.Errorf("client config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, xerrors.Errorf("kubernetes.NewForConfig: %w", err)
	}

	return clientset, nil
}
