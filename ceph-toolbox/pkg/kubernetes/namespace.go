package kubernetes

import (
	"time"

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
	coreV1Client, ns, e := nsComm(config, body)
	if e != nil {
		return e
	}

	if _, err := coreV1Client.Namespaces().Get(ns.Name, v1.GetOptions{}); err != nil {
		if _, err = coreV1Client.Namespaces().Create(&ns); err != nil {
			glog.Errorf("create namespace %s failed\n", ns.Name)
			return err
		}
	}

	return nil
}

func nsComm(config string, body string) (*corev1.CoreV1Client, apiv1.Namespace, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return nil, apiv1.Namespace{}, err
	}
	coreV1Client, err := corev1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return nil, apiv1.Namespace{}, err
	}
	ns := apiv1.Namespace{}
	if err := yaml.Unmarshal([]byte(body), &ns); err != nil {
		glog.Error(err)
		return nil, apiv1.Namespace{}, err
	}
	return coreV1Client, ns, nil
}

func DelNamespaceHandler(config, body string) error {
	return delNamespace(config, body)
}

func delNamespace(config, body string) error {
	coreV1Client, ns, e := nsComm(config, body)
	if e != nil {
		return e
	}

	if _, err := coreV1Client.Namespaces().Get(ns.Name, v1.GetOptions{}); err == nil {
		if err = coreV1Client.Namespaces().Delete(ns.Name, &v1.DeleteOptions{}); err != nil {
			glog.Errorf("delete namespace %s failed\n", ns.Name)
			return err
		}
	}

	// wait the component end
	for {
		if _, err := coreV1Client.Namespaces().Get(ns.Name, v1.GetOptions{}); err == nil {
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}

	return nil
}
