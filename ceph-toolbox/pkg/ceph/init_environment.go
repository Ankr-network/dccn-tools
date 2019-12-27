package ceph

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"github.com/Ankr-network/dccn-tools/ceph-toolbox/pkg/ssh"
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

func InitEnv(cmd *cobra.Command) {
	config, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		glog.Error(err)
		return
	}
	priKey, err := cmd.Flags().GetString("privateKey")
	if err != nil {
		glog.Error(err)
		return
	}
	name, err := cmd.Flags().GetString("user")
	if err != nil {
		glog.Error(err)
		return
	}
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		glog.Error(err)
		return
	}

	// initialize directory name
	DirName = fmt.Sprintf("%x", md5.Sum([]byte(time.Now().Format(time.RFC3339))))

	ns := GetNodeInfo(config)
	for _, v := range ns {
		if priKey != "" {
			name = v.Name
		}
		c, err := ssh.Dial(v.Addr, name, priKey, password)
		if err != nil {
			glog.Error(err)
			return
		}
		if err = makeValidDirectory(c, DirName); err != nil {
			glog.Error(err)
			return
		}
	}

}

func makeValidDirectory(c *ssh.Client, dirName string) error {
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
	cmd := fmt.Sprintf("mkdir -p %s", validPath)

	// create store directory
	_, err = c.Run(cmd)
	if err != nil {
		return err
	}

	// create new soft link
	cmd = fmt.Sprintf("mkdir -p %s", RookStorePath)
	_, err = c.Run(cmd)
	if err != nil {
		return err
	}

	cmd = fmt.Sprintf("ln -s %s %s", validPath, RookStorePath)
	_, err = c.Run(cmd)
	if err != nil {
		return err
	}

	return nil
}
