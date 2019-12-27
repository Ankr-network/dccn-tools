package ssh

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Client is a wrap of ssh client
type Client struct {
	*ssh.Client
	password string
	once     sync.Once
}

// Dial build a wrapped ssh client with given config
func Dial(addr, user, priKeyFile, password string) (*Client, error) {
	var auth ssh.AuthMethod
	if priKeyFile != "" {
		keyData, err := ioutil.ReadFile(priKeyFile)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		privateKey, err := ssh.ParsePrivateKey(keyData)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		auth = ssh.PublicKeys(privateKey)
	} else {
		auth = ssh.Password(password)
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: func(string, net.Addr, ssh.PublicKey) error { return nil },
	}

	if !strings.Contains(addr, ":") {
		addr += ":22"
	}

	c, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return &Client{
		Client:   c,
		password: password,
	}, nil
}

// BatchRun execute a group commands by ssh client
func (c *Client) BatchRun(cmds ...string) (err error) {
	for _, cmd := range cmds {
		if _, err := c.Run(cmd); err != nil {
			return err
		}
	}

	log.Printf("ssh " + c.RemoteAddr().String() + " batch run success")
	return nil
}

// Run execute a commands by ssh client
func (c *Client) Run(format string, a ...interface{}) (out string, err error) {
	log.Printf("ssh "+c.RemoteAddr().String()+" run: "+format, a...)
	defer func() {
		if err != nil {
			err = fmt.Errorf("ssh %s: %w", c.RemoteAddr().String(), err)
		}
	}()

	sess, err := c.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()

	if c.password == "" {
		data, err := sess.CombinedOutput(fmt.Sprintf(
			"bash <<EOF\n%s\nEOF", fmt.Sprintf(format, a...)))
		if err != nil {
			return "", fmt.Errorf("err: %w, output: %s", err, data)
		}
		return string(data), nil
	}

	if err := sess.RequestPty("xterm", 80, 40, ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}); err != nil {
		return "", err
	}

	// if password are given, use the password to ensure sudo permission
	var cmd string
	if err := sess.Setenv("SUDO_PASSWD", c.password); err != nil {
		cmd = fmt.Sprintf("echo %s | sudo -S true && bash <<EOF\n%s\nEOF | cat -",
			c.password, fmt.Sprintf(format, a...))
	} else {
		cmd = fmt.Sprintf("echo $SUDO_PASSWD | sudo -S true && bash <<EOF\n%s\nEOF | cat -",
			fmt.Sprintf(format, a...))
	}
	data, err := sess.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("cmd: %s, err: %w, output: %s", fmt.Sprintf(format, a...), err, data)
	}
	return string(data), nil
}

// SendFile send a file to remote with permission keeped
func (c *Client) SendFile(src, dst string) (err error) {
	if dst == "" {
		dst = src
	}

	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	stat, err := srcF.Stat()
	if err != nil {
		return err
	}

	return c.Send(dst, stat.Mode(), srcF)
}

// Send a file to remote from stream content
func (c *Client) Send(dst string, perm os.FileMode, content io.Reader) (err error) {
	log.Printf("sftp %s sending %s", c.RemoteAddr(), dst)
	defer func() {
		if err != nil {
			err = fmt.Errorf("sftp send: %w", err)
		} else {
			log.Printf("sftp %s sended %s", c.RemoteAddr(), dst)
		}
	}()

	client, err := sftp.NewClient(c.Client)
	if err != nil {
		return fmt.Errorf("init sftp client: %w", err)
	}

	// create temporary directory to ensure access permission
	c.once.Do(func() {
		rand.Seed(time.Now().Unix())
	})

	tmpDir := fmt.Sprintf("/tmp/kubeinit-%d", rand.Int())
	if err := client.Mkdir(tmpDir); err != nil {
		return err
	}
	defer client.RemoveDirectory(tmpDir)

	// write file
	dir := filepath.Dir(dst)
	file := filepath.Base(dst)
	dstF, err := client.Create(tmpDir + "/" + file)
	if err != nil {
		return err
	}
	defer dstF.Close()

	if err := dstF.Chmod(perm); err != nil {
		return err
	}

	if _, err := io.Copy(dstF, content); err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	// move file to the right path
	return c.BatchRun(
		fmt.Sprintf("sudo mkdir -p %s", dir),
		fmt.Sprintf("sudo test -e %s && sudo test ! -e %s~ && sudo mv %s %s~ || true",
			dst, dst, dst, dst),
		fmt.Sprintf("sudo mv %s/%s %s", tmpDir, file, dst),
	)
}
