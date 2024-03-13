package cmd

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	"okp4/minio-auth-plugin/app"
	"okp4/minio-auth-plugin/auth"
)

const (
	FlagNode            = "node"
	FlagCognitariumAddr = "cognitarium-addr"
	FlagServiceId       = "id"
	FlagListenAddr      = "listen-addr"
	FlagMinioURL        = "minio-url"
	FlagMinioAccessKey  = "minio-access-key"
	FlagMinioSecretKey  = "minio-secret-key"
)

var (
	nodeGrpcAddr    string
	cognitariumAddr string
	serviceID       string
	listenAddr      string
	minioURL        string
	minioAccessKey  string
	minioSecretKey  string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch the auth server",
	RunE: func(cmd *cobra.Command, args []string) error {
		minioClient, err := minio.New(minioURL, &minio.Options{
			Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
			Secure: false,
		})
		if err != nil {
			return err
		}

		app.New(
			listenAddr,
			minioClient,
			auth.New(nodeGrpcAddr, cognitariumAddr, serviceID),
		).Start()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringVar(&nodeGrpcAddr, FlagNode, "127.0.0.1:9090", "The node grpc address")
	startCmd.PersistentFlags().StringVar(&cognitariumAddr, FlagCognitariumAddr, "", "The cognitarium contract address")
	startCmd.PersistentFlags().StringVar(&serviceID, FlagServiceId, "", "The service's identifier served")
	startCmd.PersistentFlags().StringVar(&listenAddr, FlagListenAddr, "127.0.0.1:8080", "The server's listen address")
	startCmd.PersistentFlags().StringVar(&minioURL, FlagMinioURL, "", "The proxied minio URL")
	startCmd.PersistentFlags().StringVar(&minioAccessKey, FlagMinioAccessKey, "", "The minio access key")
	startCmd.PersistentFlags().StringVar(&minioSecretKey, FlagMinioSecretKey, "", "The minio secret key")
}
