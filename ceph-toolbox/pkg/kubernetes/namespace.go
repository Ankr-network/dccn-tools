package kubernetes

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func NamespaceHandler(config, body string) error {
	return namespace(config, body)
}

func namespace(config, body string) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return err
	}
	coreV1Client, err := corev1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return err
	}

	ns := apiv1.Namespace{}
	if err := yaml.Unmarshal([]byte(body), &ns); err != nil {
		glog.Error(err)
		return err
	}

	if _, err := coreV1Client.Namespaces().Get(ns.Name, v1.GetOptions{}); err != nil {
		if _, err = coreV1Client.Namespaces().Create(&ns); err != nil {
			glog.Errorf("create namespace %s failed\n", ns.Name)
			return err
		}
	}

	return nil
}
