package kubernetes

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultPath = "~/.kube/config"
)

// Kubernetes provides the API operation methods for making requests to Kubernetes
type Kubernetes struct {
	cli *k8s.Clientset
}

// NewKubernetes creates a new Clientset for the given config
func NewKubernetes() (*Kubernetes, error) {

	c := os.Getenv("K8S_CONFIG")
	if c == "" {
		var err error
		c, err = homedir.Expand(defaultPath)
		if err != nil {
			return nil, fmt.Errorf("error getting the kubernetes default path: %w", err)
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", c)
	if err != nil {
		return nil, fmt.Errorf("error building config from a kubeconfig filepath: %w", err)
	}

	cli, err := k8s.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating a new kubernetes config: %w", err)
	}

	return &Kubernetes{
		cli: cli,
	}, nil
}

// Namespaces lists all namespaces for a given cluster
func (k *Kubernetes) Namespaces() ([]core.Namespace, error) {
	ns, err := k.cli.CoreV1().Namespaces().List(meta.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing the namespaces: %w", err)
	}

	return ns.Items, nil
}
