package cmd

import (
	"github.com/spf13/cobra"
	"okp4/minio-auth-plugin/app"
	"okp4/minio-auth-plugin/auth"
)

const (
	FlagNode            = "node"
	FlagCognitariumAddr = "cognitarium-addr"
	FlagServiceId       = "id"
	FlagListenAddr      = "listen-addr"
)

var (
	nodeGrpcAddr    string
	cognitariumAddr string
	serviceID       string
	listenAddr      string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch the auth server",
	RunE: func(cmd *cobra.Command, args []string) error {
		app.New(listenAddr, auth.New(nodeGrpcAddr, cognitariumAddr, serviceID)).Start()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringVar(&nodeGrpcAddr, FlagNode, "127.0.0.1:9090", "The node grpc address")
	startCmd.PersistentFlags().StringVar(&cognitariumAddr, FlagCognitariumAddr, "", "The cognitarium contract address")
	startCmd.PersistentFlags().StringVar(&serviceID, FlagServiceId, "", "The service's identifier served")
	startCmd.PersistentFlags().StringVar(&listenAddr, FlagListenAddr, "127.0.0.1:8080", "The server's listen address")
}
