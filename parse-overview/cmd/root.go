package cmd

import (
	"fmt"
	"github.com/Ankr-network/dccn-tools/parse-overview/handle"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "parse-overview",
	Short: "parse overview.yaml",
	Long:  `It is used to parse the Overview.yaml in the chart`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handle.ParseOverview(viper.GetString("file"))
	},
}

func init() {
	rootCmd.Flags().StringP("file", "f", "", "file address to be parsed")
	rootCmd.MarkFlagRequired("file")
	viper.BindPFlags(rootCmd.Flags())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
