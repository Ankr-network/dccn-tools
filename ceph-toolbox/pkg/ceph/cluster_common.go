package ceph

import (
	"errors"
	"strings"

	"github.com/Ankr-network/dccn-tools/ceph-toolbox/pkg/kubernetes"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

const Kind = "kind"

func ConfigClusterCommon(cmd *cobra.Command) error {
	config, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		glog.Error(err)
		return err
	}

	cms := strings.Split(ClusterCommon, SplitString)
	if len(cms) == 0 {
		return errors.New("cluster common shouldn't be null")
	}
	var m map[string]interface{}
	for _, v := range cms {
		if err := yaml.Unmarshal([]byte(v), &m); err != nil {
			glog.Error(err)
			return err
		}
		if handleKey, ok := m[Kind].(string); ok {
			if err := worker(config, handleKey, v); err != nil {
				return err
			}
		} else {
			return errors.New("assert kind failed")
		}
	}
	return nil
}

func worker(config, name, body string) error {
	if err := kubernetes.ApiHandler[name](config, body); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
