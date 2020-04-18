/*
Copyright © 2020 Mobius <sv0220@163.com>

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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"

	"github.com/Ankr-network/dccn-tools/opennet-assist/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opennet-assist",
	Short: "opennet assist clear ip tool",
	Long:  ``,
	Run:   Runner,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const (
	keyPath    = "/etc/etcd/ssl/etcd-key.pem"
	certPath   = "/etc/etcd/ssl/etcd.pem"
	caCertPath = "/etc/etcd/ssl/etcd-root-ca.pem"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("data", "d", "missing.json", "clear missing ip")
	rootCmd.Flags().StringP("dev", "", "all", "device number")
	rootCmd.Flags().StringP("key", "", keyPath, "etcd client key")
	rootCmd.Flags().StringP("crt", "", certPath, "etcd client cert")
	rootCmd.Flags().StringP("ca", "", caCertPath, "etcd root ca")
	rootCmd.Flags().StringP("endpoints", "",
		"https://10.28.5.248:2379, https://10.28.5.249:2379,https://10.28.5.250:2379",
		"etcd endpoints")
}

const (
	rootPath       = "/opennet/ip"
	lockPath       = "/opennet/lock"
	requestTimeout = 30 * time.Second
)

func Runner(cmd *cobra.Command, args []string) {
	log.Info().Msg("clear work start ...")
	crt, err := cmd.Flags().GetString("crt")
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	key, err := cmd.Flags().GetString("key")
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	ca, err := cmd.Flags().GetString("ca")
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	endpoints, err := cmd.Flags().GetString("endpoints")
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	client, err := pkg.NewETCDClient(crt, key, ca, strings.Split(endpoints, ","))
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	// get data from data file, parse and do clear work
	fileName, err := cmd.Flags().GetString("data")
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	ips, err := getDataFromFile(fileName)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	if err := recycleIP(client, cmd, ips); err != nil {
		log.Error().Msg(err.Error())
		return
	}
	log.Info().Msg("recycle ip work over.")
}

type IPSet map[string]string

func getDataFromFile(fileName string) (IPSet, error) {
	ips := make(IPSet)
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Error().Msg(err.Error())
		return ips, err
	}
	if err := json.Unmarshal(data, &ips); err != nil {
		return ips, err
	}
	return ips, nil
}

func recycleIP(client *clientv3.Client, cmd *cobra.Command, rips IPSet) error {
	// lock here
	session, err := concurrency.NewSession(client)
	if err != nil {
		return err
	}
	defer session.Close()
	dev, err := cmd.Flags().GetString("dev")
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	mux := concurrency.NewMutex(session, fmt.Sprintf("%s/%s", lockPath, dev))
	if err := mux.Lock(context.TODO()); err != nil {
		return err
	}
	defer func() {
		if err := mux.Unlock(context.TODO()); err != nil {
			log.Error().Msg(err.Error())
		}
	}()

	// https://github.com/Ankr-network/dccn-tools ip
	ipStorePath := fmt.Sprintf("%s/%s", rootPath, dev)
	kvc := clientv3.NewKV(client)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	rsp, err := kvc.Get(ctx, ipStorePath)
	cancel()
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	ipset := make(IPSet)
	ipSetSrv := make(IPSet)
	if err := json.Unmarshal(rsp.Kvs[0].Value, &ipset); err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	// recycle the ip and restore into etcd
	for k, v := range ipset {
		if _, ok := rips[k]; !ok {
			ipSetSrv[k] = v
		} else {
			log.Info().Msgf("id: %s ip: %s", k, v)
		}
	}

	data, err := json.Marshal(&ipSetSrv)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = kvc.Put(ctx, ipStorePath, string(data))
	cancel()
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return nil
}
