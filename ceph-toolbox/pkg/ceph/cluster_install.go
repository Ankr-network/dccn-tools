package ceph

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/Ankr-network/dccn-tools/ceph-toolbox/pkg/kubernetes"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	rookApi "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	rookv1 "github.com/rook/rook/pkg/client/clientset/versioned/typed/ceph.rook.io/v1"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	StatusOK         = "HEALTH_OK"
	ToolBoxFilter    = "app=rook-ceph-tools"
	ToolBoxContainer = "rook-ceph-tools"
	RookNamespace    = "rook-ceph"
	queryInterval    = 2
)

func InstallCluster(cmd *cobra.Command) error {

	config, cephCluster, cephV1Client, e := clusterComm(cmd)
	if e != nil {
		return e
	}

	if _, err := cephV1Client.CephClusters(cephCluster.Namespace).Get(cephCluster.Name, v1.GetOptions{}); err != nil {
		if _, err := cephV1Client.CephClusters(cephCluster.Namespace).Create(&cephCluster); err != nil {
			glog.Error(err)
			return err
		}
	}

	// install cluster tool box , for check ceph cluster health
	if err := kubernetes.DeploymentHandler(config, ClusterToolBox); err != nil {
		glog.Error(err)
		return err
	}

	pc, err := kubernetes.NewPodClient(config, RookNamespace)
	if err != nil {
		glog.Error(err)
		return err
	}

	var toolBoxPod *apiv1.Pod

	// replace standard output and save it
	r, w, err := os.Pipe()
	if err != nil {
		glog.Error(err)
		return err
	}
	oldOutput := os.Stdout
	os.Stdout = w
	glog.Info("save output pipe ")
	elapseTime := 0
	for {

		time.Sleep(queryInterval * time.Second)
		elapseTime += queryInterval

		toolBoxPod, err = pc.GetPodByFilter(ToolBoxFilter)
		if err != nil {
			glog.Error(err)
			return err
		} else if toolBoxPod == nil {
			glog.Infof("cluster install in progress, elapse %d s\n", elapseTime)
			continue
		} else if toolBoxPod.Status.Phase != apiv1.PodRunning {
			glog.Infof("cluster install in progress, elapse %d s\n", elapseTime)
			continue
		}

		if err = pc.ExecInPod(toolBoxPod, ToolBoxContainer, []string{"ceph", "status"},
			os.Stdout, bytes.NewBuffer([]byte{})); err != nil {
			glog.Infof("cluster install in progress, elapse %d s\n", elapseTime)
			continue
		} else {
			break
		}

	}
	glog.Info("read status from cluster")
	// recover standard output
	glog.Info("recover output pipe  ")
	_ = w.Close()
	os.Stdout = oldOutput
	rs, err := ioutil.ReadAll(r)
	if err != nil {
		glog.Error(err)
	}
	glog.Infof("cluster install status: %s \n ", getClusterStatus(rs))

	return nil
}

func clusterComm(cmd *cobra.Command) (string, rookApi.CephCluster, *rookv1.CephV1Client, error) {
	config, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		glog.Error(err)
		return "", rookApi.CephCluster{}, nil, err
	}
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return "", rookApi.CephCluster{}, nil, err
	}
	ciBytes, err := getClusterInstallTpl(ClusterInstall)
	if err != nil {
		glog.Error(err)
		return "", rookApi.CephCluster{}, nil, err
	}
	cephCluster := rookApi.CephCluster{}
	if err := yaml.Unmarshal(ciBytes, &cephCluster); err != nil {
		glog.Error(err)
		return "", rookApi.CephCluster{}, nil, err
	}
	cephV1Client, err := rookv1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return "", rookApi.CephCluster{}, nil, err
	}
	return config, cephCluster, cephV1Client, nil
}

func getClusterStatus(rs []byte) string {
	var r string
	s := bufio.NewScanner(bytes.NewBuffer(rs))
	for s.Scan() {
		if strings.Contains(s.Text(), "health:") {
			r = s.Text()
		}
	}
	res := strings.Split(r, ":")
	return strings.Trim(res[1], " ")
}

func getClusterInstallTpl(body string) ([]byte, error) {
	// replace valid arguments
	buf := bytes.NewBuffer([]byte{})

	tpl, err := template.New("cluster_install").Parse(body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	arg := struct {
		HostDirPath string
	}{
		HostDirPath: fmt.Sprintf("%s/%s", RookStorePath, DirName),
	}

	if err := tpl.Execute(buf, &arg); err != nil {
		glog.Error(err)
		return nil, err
	}

	return buf.Bytes(), nil

}

func DeleteCluster(cmd *cobra.Command) error {
	config, cephCluster, cephV1Client, e := clusterComm(cmd)
	if e != nil {
		return e
	}

	// install cluster tool box , for check ceph cluster health
	if err := kubernetes.DelDeploymentHandler(config, ClusterToolBox); err != nil {
		glog.Error(err)
		return err
	}

	if _, err := cephV1Client.CephClusters(cephCluster.Namespace).Get(cephCluster.Name, v1.GetOptions{}); err == nil {
		if err := cephV1Client.CephClusters(cephCluster.Namespace).Delete(cephCluster.Name, &v1.DeleteOptions{}); err != nil {
			glog.Error(err)
			return err
		}
	}
	// wait the component end
	for {
		if _, err := cephV1Client.CephClusters(cephCluster.Namespace).Get(cephCluster.Name, v1.GetOptions{}); err == nil {
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}

	return nil
}
