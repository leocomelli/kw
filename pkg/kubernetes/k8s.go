package kubernetes

import (
	"bufio"
	"context"
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

// LogStream provides the data from a log entry
type LogStream struct {
	Namespace string
	Pod       string
	Container string
	Message   string
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

// Pods lists all pods for a given namespace
func (k *Kubernetes) Pods(ns string) ([]core.Pod, error) {
	pods, err := k.cli.CoreV1().Pods(ns).List(meta.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing the namespaces: %w", err)
	}

	return pods.Items, nil
}

// Pod gets a pod resource for a given namespace and pod
func (k *Kubernetes) Pod(ns, p string) (*core.Pod, error) {
	opts := meta.GetOptions{}

	pod, err := k.cli.CoreV1().Pods(ns).Get(p, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting the pod %s/%s: %w", ns, p, err)
	}

	return pod, nil
}

// ListPods lists all pods for a given namespace or by pod name
func (k *Kubernetes) ListPods(ns, p string) ([]core.Pod, error) {
	if p != "" {
		pod, err := k.Pod(ns, p)
		if err != nil {
			return nil, err
		}

		return []core.Pod{*pod}, nil
	}

	ps, err := k.Pods(ns)
	if err != nil {
		return nil, err
	}

	return ps, nil
}

// Logs streams the logs for a given namespace, pod or container
func (k *Kubernetes) Logs(ns, p, c string, stream chan *LogStream) error {
	pods, err := k.ListPods(ns, p)
	if err != nil {
		return err
	}

	ctx, _ := context.WithCancel(context.Background())

	for _, pod := range pods {
		// use the container name specified or list all containers in pod
		var cnames []string
		if c != "" {
			cnames = append(cnames, c)
		} else {
			for _, ctn := range pod.Spec.Containers {
				cnames = append(cnames, ctn.Name)
			}
		}

		for _, ctn := range cnames {
			opts := &core.PodLogOptions{Follow: true, Container: ctn}

			req := k.cli.CoreV1().Pods(ns).GetLogs(pod.GetName(), opts)
			readCloser, err := req.Context(ctx).Stream()
			if err != nil {
				return err
			}

			go func(pod, container string) {
				scanner := bufio.NewScanner(readCloser)
				for scanner.Scan() {
					line := scanner.Text()
					logStream := &LogStream{
						ns,
						pod,
						container,
						line,
					}

					stream <- logStream
				}
			}(pod.GetName(), ctn)
		}
	}

	done := make(chan int)
	<-done

	return nil
}
