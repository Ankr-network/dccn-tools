package kubernetes

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	storagev1 "k8s.io/client-go/kubernetes/typed/storage/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func StorageClassHandler(config, body string) error {
	return buildStorageClass(config, body)
}

func buildStorageClass(config, body string) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return err
	}

	storageClient, err := storagev1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return err
	}

	storageClass := v1.StorageClass{}
	if err := yaml.Unmarshal([]byte(body), &storageClass); err != nil {
		glog.Error(err)
		return err
	}

	if _, err = storageClient.StorageClasses().Get(storageClass.Name, metav1.GetOptions{}); err != nil {
		if _, err = storageClient.StorageClasses().Create(&storageClass); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}
