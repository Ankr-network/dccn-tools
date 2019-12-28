//+build !windows

package ceph

import (
	"log"
	"os/exec"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	Master           = "node-role.kubernetes.io/master"
	CpuMinimum int64 = 4
	MemMinimum int64 = 7000 * KB
)

type NodeInfo struct {
	Name string
	Addr string
	User string
	CPU  int64
	Mem  int64 // Ki
}

// GetNodeInfo  get all work node info except master
func GetNodeInfo(kconfig string) []*NodeInfo {
	userName := GetUserName()
	nodes := make([]*NodeInfo, 0)
	config, err := clientcmd.BuildConfigFromFlags("", kconfig)
	if err != nil {
		log.Println(err)
		return nil
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err)
		return nil
	}
	ns, err := clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil
	}
	for _, v := range ns.Items {
		if func() bool {
			for _, t := range v.Spec.Taints {
				if t.Key == Master {
					return true
				}
			}
			return false
		}() {
			continue
		}

		ni := &NodeInfo{}
		ni.Name = v.Name
		ni.Addr = v.Status.Addresses[0].Address
		ni.User = userName
		ni.CPU = v.Status.Capacity.Cpu().Value()
		ni.Mem = v.Status.Capacity.Memory().Value()
		nodes = append(nodes, ni)
	}
	return nodes
}

func GetUserName() string {
	whoami := exec.Command("whoami")
	rsp, err := whoami.Output()
	if err != nil {
		panic(err)
	}
	rs := strings.Trim(string(rsp), " ")
	rs = strings.Trim(rs, "\n")
	return rs
}
