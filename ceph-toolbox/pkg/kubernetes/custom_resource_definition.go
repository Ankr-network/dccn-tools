package kubernetes

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	apiextension "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func CustomResourceDefinitionHandler(config, body string) error {
	return buildCustomResourceDefinition(config, body)
}

func buildCustomResourceDefinition(config, body string) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return err
	}

	crdClient, err := client.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return err
	}

	crd := apiextension.CustomResourceDefinition{}
	if err := yaml.Unmarshal([]byte(body), &crd); err != nil {
		glog.Error(err)
		return err
	}

	if _, err = crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name, v1.GetOptions{}); err != nil {
		if _, err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&crd); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}
