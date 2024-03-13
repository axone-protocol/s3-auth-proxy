package cmd

import (
	"context"
	"crypto/tls"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	grpccreds "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"okp4/s3-auth-proxy/app"
	"okp4/s3-auth-proxy/auth"
	"okp4/s3-auth-proxy/dataverse"
	"time"
)

const (
	FlagNodeGrpc          = "node-grpc"
	FlagGrpcNoTLS         = "grpc-no-tls"
	FlagGrpcTLSSkipVerify = "grpc-tls-skip-verify"
	FlagDataverseAddr     = "dataverse-addr"
	FlagServiceId         = "id"
	FlagListenAddr        = "listen-addr"
	FlagJWTSecretKey      = "jwt-secret-key"
	FlagS3Endpoint        = "s3-endpoint"
	FlagS3AccessKey       = "s3-access-key"
	FlagS3SecretKey       = "s3-secret-key"
	FlagS3Insecure        = "s3-insecure"
)

var (
	nodeGrpcAddr      string
	grpcNoTls         bool
	grpcTlsSkipVerify bool
	dataverseAddr     string
	serviceID         string
	listenAddr        string
	jwtSecretKey      []byte
	s3Endpoint        string
	s3AccessKey       string
	s3SecretKey       string
	s3Insecure        bool
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

		ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFn()
		dataverseClient, err := dataverse.NewClient(ctx, nodeGrpcAddr, dataverseAddr, getTransportCredentials())
		if err != nil {
			return err
		}

		app.New(
			listenAddr,
			s3Client,
			auth.New(jwtSecretKey, dataverseClient, serviceID),
		).Start()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringVar(&nodeGrpcAddr, FlagNodeGrpc, "127.0.0.1:9090", "The node grpc address")
	startCmd.PersistentFlags().BoolVar(&grpcNoTls, FlagGrpcNoTLS, false, "No encryption with the GRPC endpoint")
	startCmd.PersistentFlags().BoolVar(&grpcTlsSkipVerify,
		FlagGrpcTLSSkipVerify,
		false,
		"Encryption with the GRPC endpoint but skip certificates verification")
	startCmd.PersistentFlags().StringVar(&dataverseAddr, FlagDataverseAddr, "", "The dataverse contract address")
	startCmd.PersistentFlags().StringVar(&serviceID, FlagServiceId, "", "The service's identifier served")
	startCmd.PersistentFlags().StringVar(&listenAddr, FlagListenAddr, "127.0.0.1:8080", "The server's listen address")
	startCmd.PersistentFlags().BytesHexVar(&jwtSecretKey, FlagJWTSecretKey, []byte{}, "The hex encoded secret key used to issue JWT tokens")
	startCmd.PersistentFlags().StringVar(&s3Endpoint, FlagS3Endpoint, "", "The S3 endpoint to proxy")
	startCmd.PersistentFlags().StringVar(&s3AccessKey, FlagS3AccessKey, "", "The S3 access key")
	startCmd.PersistentFlags().StringVar(&s3SecretKey, FlagS3SecretKey, "", "The S3 secret key")
	startCmd.PersistentFlags().BoolVar(&s3Insecure, FlagS3Insecure, false, "If specified we'll accept non encrypted connection with the S3")
}

func getTransportCredentials() grpccreds.TransportCredentials {
	switch {
	case grpcNoTls:
		return insecure.NewCredentials()
	case grpcTlsSkipVerify:
		return grpccreds.NewTLS(&tls.Config{InsecureSkipVerify: true}) // #nosec G402 : skip lint since it's an optional flag
	default:
		return grpccreds.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
	}
}
