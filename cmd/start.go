package cmd

import (
	"context"
	"crypto/tls"
	"io"
	"time"

	"github.com/axone-protocol/axone-sdk/http"

	"github.com/piprate/json-gold/ld"

	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/axone-protocol/axone-sdk/keys"
	"github.com/axone-protocol/axone-sdk/provider/storage"
	"google.golang.org/grpc"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	grpccreds "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// nolint: gosec
const (
	FlagNodeGrpc          = "node-grpc"
	FlagGrpcNoTLS         = "grpc-no-tls"
	FlagGrpcTLSSkipVerify = "grpc-tls-skip-verify"
	FlagDataverseAddr     = "dataverse-addr"
	FlagServiceMnemonic   = "svc-mnemonic"
	FlagServiceBaseURL    = "svc-base-url"
	FlagListenAddr        = "listen-addr"
	FlagJWTSecretKey      = "jwt-secret-key"
	FlagJWTDuration       = "jwt-duration"
	FlagS3Endpoint        = "s3-endpoint"
	FlagS3Bucket          = "s3-bucket"
	FlagS3AccessKey       = "s3-access-key"
	FlagS3SecretKey       = "s3-secret-key"
	FlagS3Insecure        = "s3-insecure"
)

var (
	nodeGrpcAddr      string
	grpcNoTLS         bool
	grpcTLSSkipVerify bool
	dataverseAddr     string
	mnemonic          string
	baseURL           string
	listenAddr        string
	jwtSecretKey      []byte
	jwtDuration       time.Duration
	s3Endpoint        string
	s3Bucket          string
	s3AccessKey       string
	s3SecretKey       string
	s3Insecure        bool
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch the auth server",
	RunE: func(_ *cobra.Command, _ []string) error {
		s3Client, err := minio.New(s3Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(s3AccessKey, s3SecretKey, ""),
			Secure: false,
		})
		if err != nil {
			return err
		}

		key, err := keys.NewKeyFromMnemonic(mnemonic)
		if err != nil {
			return err
		}

		ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFn()

		dvClient, err := dataverse.NewQueryClient(ctx, nodeGrpcAddr, dataverseAddr, grpc.WithTransportCredentials(getTransportCredentials()))
		if err != nil {
			return err
		}

		storageProxy, err := storage.NewProxy(
			ctx,
			key,
			baseURL,
			dvClient,
			ld.NewCachingDocumentLoader(ld.NewDefaultDocumentLoader(nil)),
			func(ctx context.Context, id string) (io.Reader, error) {
				return s3Client.GetObject(ctx, s3Bucket, id, minio.GetObjectOptions{})
			},
			func(ctx context.Context, s string, reader io.Reader) error {
				_, err := s3Client.PutObject(ctx, s3Bucket, s, reader, -1, minio.PutObjectOptions{})
				return err
			},
		)
		if err != nil {
			return err
		}

		return http.NewServer(
			listenAddr,
			storageProxy.HTTPConfigurator(jwtSecretKey, jwtDuration),
		).Listen()
	},
}

// nolint: lll
func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringVar(&nodeGrpcAddr, FlagNodeGrpc, "127.0.0.1:9090", "The node grpc address")
	startCmd.PersistentFlags().BoolVar(&grpcNoTLS, FlagGrpcNoTLS, false, "No encryption with the GRPC endpoint")
	startCmd.PersistentFlags().BoolVar(&grpcTLSSkipVerify,
		FlagGrpcTLSSkipVerify,
		false,
		"Encryption with the GRPC endpoint but skip certificates verification")
	startCmd.PersistentFlags().StringVar(&dataverseAddr, FlagDataverseAddr, "", "The dataverse contract address")
	startCmd.PersistentFlags().StringVar(&mnemonic, FlagServiceMnemonic, "", "The service's mnemonic")
	startCmd.PersistentFlags().StringVar(&baseURL, FlagServiceBaseURL, "", "The service's base URL")
	startCmd.PersistentFlags().StringVar(&listenAddr, FlagListenAddr, "127.0.0.1:8080", "The server's listen address")
	startCmd.PersistentFlags().BytesHexVar(&jwtSecretKey, FlagJWTSecretKey, []byte{}, "The hex encoded secret key used to issue JWT tokens")
	startCmd.PersistentFlags().DurationVar(&jwtDuration, FlagJWTDuration, time.Hour, "The JWT token duration")
	startCmd.PersistentFlags().StringVar(&s3Endpoint, FlagS3Endpoint, "", "The S3 endpoint to proxy")
	startCmd.PersistentFlags().StringVar(&s3Bucket, FlagS3Bucket, "data", "The S3 bucket to proxy")
	startCmd.PersistentFlags().StringVar(&s3AccessKey, FlagS3AccessKey, "", "The S3 access key")
	startCmd.PersistentFlags().StringVar(&s3SecretKey, FlagS3SecretKey, "", "The S3 secret key")
	startCmd.PersistentFlags().BoolVar(&s3Insecure, FlagS3Insecure, false, "If specified we'll accept non encrypted connection with the S3")
}

// nolint: gosec
func getTransportCredentials() grpccreds.TransportCredentials {
	switch {
	case grpcNoTLS:
		return insecure.NewCredentials()
	case grpcTLSSkipVerify:
		return grpccreds.NewTLS(&tls.Config{InsecureSkipVerify: true}) // #nosec G402 : skip lint since it's an optional flag
	default:
		return grpccreds.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
	}
}
