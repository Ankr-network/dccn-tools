package kubernetes

import (
	"time"

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
	err, crdClient, crd, e := crdComm(config, body)
	if e != nil {
		return e
	}

	if _, err = crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name, v1.GetOptions{}); err != nil {
		if _, err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&crd); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}

func crdComm(config string, body string) (error, *client.Clientset, apiextension.CustomResourceDefinition, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return nil, nil, apiextension.CustomResourceDefinition{}, err
	}
	crdClient, err := client.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return nil, nil, apiextension.CustomResourceDefinition{}, err
	}
	crd := apiextension.CustomResourceDefinition{}
	if err := yaml.Unmarshal([]byte(body), &crd); err != nil {
		glog.Error(err)
		return nil, nil, apiextension.CustomResourceDefinition{}, err
	}
	return err, crdClient, crd, nil
}

func DelCustomResourceDefinitionHandler(config, body string) error {
	return delCustomResourceDefinition(config, body)
}

func delCustomResourceDefinition(config, body string) error {
	err, crdClient, crd, e := crdComm(config, body)
	if e != nil {
		return e
	}

	if _, err = crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name,
		v1.GetOptions{}); err == nil {
		if err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, &v1.DeleteOptions{}); err != nil {
			glog.Error(err)
			return err
		}
	}

	// wait the component end
	for {
		if _, err = crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name,
			v1.GetOptions{}); err == nil {
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}

	return nil
}
