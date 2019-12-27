package kubernetes

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/golang/glog"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type PodClient struct {
	namespace  string
	clientSet  *kubernetes.Clientset
	config     clientcmd.ClientConfig
	restConfig *rest.Config
}

func NewPodClient(kcfg, namespace string) (*PodClient, error) {
	client := &PodClient{
		namespace: namespace,
	}
	var err error
	if kcfg == "" {
		return nil, errors.New("kubeconfig shouldn't be null")
	}
	if err := os.Setenv("KUBECONFIG", kcfg); err != nil {
		glog.Error(err)
		return nil, err
	}

	client.config = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{})

	client.restConfig, err = client.config.ClientConfig()
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	client.clientSet, err = kubernetes.NewForConfig(client.restConfig)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	return client, nil
}

func (p *PodClient) GetPodByFilter(filter string) (*apiv1.Pod, error) {
	podClient := p.clientSet.CoreV1().Pods(p.namespace)
	pods, err := podClient.List(metav1.ListOptions{
		LabelSelector: filter,
	})

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	var pod apiv1.Pod
	for _, p := range pods.Items {
		if p.Status.Phase == apiv1.PodRunning {
			pod = p
			break
		}
	}

	return &pod, nil
}

func (p *PodClient) ExecInPod(pod *apiv1.Pod, container string, commands []string, stdout, stderr io.Writer) error {
	restClient := p.clientSet.CoreV1().RESTClient()

	req := restClient.Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		Param("container", container).
		Param("stdout", "true").
		Param("stderr", "true")

	for _, command := range commands {
		req.Param("command", command)
	}

	executor, err := remotecommand.NewSPDYExecutor(p.restConfig, http.MethodPost, req.URL())
	if err != nil {
		glog.Error(err)
		return err
	}

	return executor.Stream(remotecommand.StreamOptions{
		Stdin:             nil,
		Stdout:            stdout,
		Stderr:            stderr,
		Tty:               false,
		TerminalSizeQueue: nil,
	})
}
