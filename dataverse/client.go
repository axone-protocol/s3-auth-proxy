package dataverse

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	logictypes "github.com/okp4/okp4d/v7/x/logic/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	wasmClient      wasmtypes.QueryClient
	logicClient     logictypes.QueryServiceClient
	cognitariumAddr string
}

func NewClient(ctx context.Context, nodeGrpc, dataverseAddr string, transportCreds credentials.TransportCredentials) (*Client, error) {
	grpcConn, err := grpc.Dial(
		nodeGrpc,
		grpc.WithTransportCredentials(transportCreds),
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't create grpc connection: %w", err)
	}

	wasmClient := wasmtypes.NewQueryClient(grpcConn)
	cognitariumAddr, err := queryCognitariumAddr(ctx, dataverseAddr, wasmClient)
	if err != nil {
		return nil, fmt.Errorf("couldn't query cognitarium address: %w", err)
	}

	return &Client{
		wasmClient:      wasmClient,
		logicClient:     logictypes.NewQueryServiceClient(grpcConn),
		cognitariumAddr: cognitariumAddr,
	}, nil
}

func queryCognitariumAddr(ctx context.Context, dataverseAddr string, wasmClient types.QueryClient) (string, error) {
	query, err := json.Marshal(map[string]interface{}{
		"dataverse": struct{}{},
	})
	if err != nil {
		return "", err
	}

	resp, err := wasmClient.SmartContractState(ctx, &types.QuerySmartContractStateRequest{
		Address:   dataverseAddr,
		QueryData: query,
	})
	if err != nil {
		return "", err
	}

	var data struct {
		TriplestoreAddress string `json:"triplestore_address"`
	}
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}
	return data.TriplestoreAddress, nil
}
