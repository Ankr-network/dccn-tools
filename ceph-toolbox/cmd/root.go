/*
Copyright © 2019 Mobius <sv0220@163.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ceph-toolbox",
	Short: "ceph tool box",
	Long:  `automatic install ceph cluster on the kubernetes cluster`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	goflag.Parse()

	// initialize global flags
	rootCmd.PersistentFlags().StringP("kubeconfig", "c", "/root/.kube/config",
		"kubernetes cluster config")
	rootCmd.PersistentFlags().StringP("privateKey", "k", "/root/.ssh/id_rsa",
		"host private key file")
	rootCmd.PersistentFlags().StringP("user", "u", "root",
		"user name")
	rootCmd.PersistentFlags().StringP("password", "p", "1q2w!Q@W",
		"password")

	if err := viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig")); err != nil {
		glog.Error(err)
		return
	}
	if err := viper.BindPFlag("privateKey", rootCmd.PersistentFlags().Lookup("privateKey")); err != nil {
		glog.Error(err)
		return
	}
	if err := viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user")); err != nil {
		glog.Error(err)
		return
	}
	if err := viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password")); err != nil {
		glog.Error(err)
		return
	}

}
