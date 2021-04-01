package k8s

import (
	"io/ioutil"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClientInCluster() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

func NewClientFromKubeConfig(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

// Namespace returns the current namespace from running pod
func Namespace() (string, error) {
	const (
		namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	)
	namespace, err := ioutil.ReadFile(namespaceFile)
	if err != nil {
		return "", err
	}
	return string(namespace), nil
}
