package impl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Ankr-network/dccn-tools/ceph-toolbox/pkg/ssh/face"
	"github.com/golang/glog"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	c *ssh.Client
}

func (c *Client) Run(cmd string) (string, error) {
	sess, err := c.c.NewSession()
	if err != nil {
		glog.Error(err)
		return "", err
	}
	defer sess.Close()
	var buf bytes.Buffer
	sess.Stdout = &buf
	if err := sess.Run(cmd); err != nil {
		glog.Error(err)
		return "", nil
	}
	return buf.String(), nil
}

// key file : $HOME/$USER/.ssh/id_rsa
func getKeyFile(keyFile string) (key ssh.Signer, err error) {
	buf, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}
	return
}

func NewClient(host, user, password, priKey string) (face.Issh, error) {
	var auth ssh.AuthMethod
	if priKey != "" {
		key, err := getKeyFile(priKey)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		auth = ssh.PublicKeys(key)
	} else {
		auth = ssh.Password(password)
	}
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if !strings.Contains(host, ":") {
		host = fmt.Sprintf("%s:22", host)
	}
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return &Client{c: client}, nil
}
