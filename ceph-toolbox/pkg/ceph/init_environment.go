package ceph

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"github.com/Ankr-network/dccn-tools/ceph-toolbox/pkg/ssh/face"
	"github.com/Ankr-network/dccn-tools/ceph-toolbox/pkg/ssh/impl"
	"github.com/Ankr-network/dccn-tools/ceph-toolbox/pkg/sys"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

const (
	RookStorePath = "/var/ankr/rook"
	KubeConfig    = "KUBECONFIG"
)

const (
	_ = 1 << (10 * iota)
	KB
	MB
	GB
	TB
)

var (
	DirName     string
	DbMinSizeMB int = 1024 * 1024 * 1024 * MB
)

func InitEnv(cmd *cobra.Command) error {
	config, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		glog.Error(err)
		return err
	}
	priKey, err := cmd.Flags().GetString("privateKey")
	if err != nil {
		glog.Error(err)
		return err
	}
	name, err := cmd.Flags().GetString("user")
	if err != nil {
		glog.Error(err)
		return err
	}
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		glog.Error(err)
		return err
	}

	// initialize directory name
	DirName = fmt.Sprintf("%x", md5.Sum([]byte(time.Now().Format(time.RFC3339))))

	ns := GetNodeInfo(config)
	if err := checkEnvCondition(ns); err != nil {
		glog.Error(err)
		return err
	}
	for _, v := range ns {

		glog.Infof("host: %s user: %s password: ****** private: %s \n", v.Addr, name, priKey)
		c, err := impl.NewClient(v.Addr, name, password, priKey)
		if err != nil {
			glog.Error(err)
			return err
		}

		if err = makeValidDirectory(c, DirName); err != nil {
			glog.Error(err)
			return err
		}
	}
	return nil
}

func checkEnvCondition(ns []*NodeInfo) error {
	if len(ns) < NodeMinimum {
		glog.Infof("basic condition, minimum nodes: %d current nodes: %d \n", NodeMinimum, len(ns))
		errors.New("kubernetes cluster environment condition can't satisfy with ceph cluster")
	}
	// output env info
	for _, v := range ns {
		glog.Infof("ip: %s node: %s user: %s cpu: %d mem: %d \n",
			v.Addr, v.Name, v.User, v.CPU, v.Mem)
		if v.CPU < CpuMinimum || v.Mem < MemMinimum {
			glog.Infof("basic condition, minimum cpu cores: %d memory size: %d \n",
				CpuMinimum, MemMinimum)
			return errors.New("kubernetes cluster environment condition can't satisfy with ceph cluster")
		}
	}
	return nil
}

func makeValidDirectory(c face.Issh, dirName string) error {
	// find out max size store volume and the path which mounted
	rsp, err := c.Run("df -m")
	if err != nil {
		return err
	}
	maxPath, maxSize := sys.LookUpValidDisk(rsp)
	if maxPath == "/" {
		maxPath = ""
	}
	if maxSize == 0 {
		return errors.New("volume shouldn't be 0 MB")
	}

	// if the value less than global DbMinSizeMB then assign
	if maxSize < DbMinSizeMB {
		DbMinSizeMB = maxSize
	}

	validPath := fmt.Sprintf("%s/rook/%s", maxPath, dirName)
	cmd := fmt.Sprintf("sudo mkdir -p %s", validPath)
	glog.Info("exec cmd: ", cmd)
	// create store directory
	_, err = c.Run(cmd)
	if err != nil {
		return err
	}

	// create new soft link
	cmd = fmt.Sprintf("sudo mkdir -p %s", RookStorePath)
	glog.Info("exec cmd: ", cmd)
	_, err = c.Run(cmd)
	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("sudo ln -s %s %s", validPath, RookStorePath)
	glog.Info("exec cmd: ", cmd)
	_, err = c.Run(cmd)
	if err != nil {
		return err
	}

	return nil
}
