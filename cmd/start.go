package cmd

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	"okp4/s3-auth-proxy/app"
	"okp4/s3-auth-proxy/auth"
)

const (
	FlagNode            = "node"
	FlagCognitariumAddr = "cognitarium-addr"
	FlagServiceId       = "id"
	FlagListenAddr      = "listen-addr"
	FlagJWTSecretKey    = "jwt-secret-key"
	FlagS3Endpoint      = "s3-endpoint"
	FlagS3AccessKey     = "s3-access-key"
	FlagS3SecretKey     = "s3-secret-key"
	FlagS3Insecure      = "s3-insecure"
)

var (
	nodeGrpcAddr    string
	cognitariumAddr string
	serviceID       string
	listenAddr      string
	jwtSecretKey    []byte
	s3Endpoint      string
	s3AccessKey     string
	s3SecretKey     string
	s3Insecure      bool
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch the auth server",
	RunE: func(cmd *cobra.Command, args []string) error {
		s3Client, err := minio.New(s3Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(s3AccessKey, s3SecretKey, ""),
			Secure: false,
		})
		if err != nil {
			return err
		}

		app.New(
			listenAddr,
			s3Client,
			auth.New(jwtSecretKey, nodeGrpcAddr, cognitariumAddr, serviceID),
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
	startCmd.PersistentFlags().BytesHexVar(&jwtSecretKey, FlagJWTSecretKey, []byte{}, "The hex encoded secret key used to issue JWT tokens")
	startCmd.PersistentFlags().StringVar(&s3Endpoint, FlagS3Endpoint, "", "The S3 endpoint to proxy")
	startCmd.PersistentFlags().StringVar(&s3AccessKey, FlagS3AccessKey, "", "The S3 access key")
	startCmd.PersistentFlags().StringVar(&s3SecretKey, FlagS3SecretKey, "", "The S3 secret key")
	startCmd.PersistentFlags().BoolVar(&s3Insecure, FlagS3Insecure, false, "If specified we'll accept non encrypted connection with the S3")
}
