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

func (c *Client) GetExecutionOrderContext(ctx context.Context, order, execSvc string) (string, error) {
	resp, err := c.queryCognitariumSelect(ctx, Select{
		Query: SelectQuery{
			Prefixes: []Prefix{{
				Prefix:    "exec",
				Namespace: "https://w3id.org/okp4/ontology/v3/schema/credential/digital-service/execution-order/",
			}},
			Select: []SelectItem{{Variable: "zone"}},
			Where: []WhereCondition{
				{
					Simple: SimpleWhereCondition{TriplePattern: TriplePattern{
						Subject:   VarOrNode{Variable: "credId"},
						Predicate: VarOrNamedNode{NamedNode: &IRI{Full: "dataverse:credential#subject"}},
						Object:    VarOrNodeOrLiteral{Node: &Node{NamedNode: &IRI{Full: order}}},
					}},
				},
				{
					Simple: SimpleWhereCondition{TriplePattern: TriplePattern{
						Subject:   VarOrNode{Variable: "credId"},
						Predicate: VarOrNamedNode{NamedNode: &IRI{Full: "dataverse:credential#type"}},
						Object:    VarOrNodeOrLiteral{Node: &Node{NamedNode: &IRI{Prefixed: "exec:DigitalServiceExecutionOrderCredential"}}},
					}},
				},
				{
					Simple: SimpleWhereCondition{TriplePattern: TriplePattern{
						Subject:   VarOrNode{Variable: "credId"},
						Predicate: VarOrNamedNode{NamedNode: &IRI{Full: "dataverse:credential#claim"}},
						Object:    VarOrNodeOrLiteral{Variable: "claim"},
					}},
				},
				{
					Simple: SimpleWhereCondition{TriplePattern: TriplePattern{
						Subject:   VarOrNode{Variable: "claim"},
						Predicate: VarOrNamedNode{NamedNode: &IRI{Prefixed: "exec:executedBy"}},
						Object:    VarOrNodeOrLiteral{Node: &Node{NamedNode: &IRI{Full: execSvc}}},
					}},
				},
				{
					Simple: SimpleWhereCondition{TriplePattern: TriplePattern{
						Subject:   VarOrNode{Variable: "claim"},
						Predicate: VarOrNamedNode{NamedNode: &IRI{Prefixed: "exec:inZone"}},
						Object:    VarOrNodeOrLiteral{Variable: "zone"},
					}},
				},
			},
			Limit: 1,
		},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Results.Bindings) != 1 {
		return "", fmt.Errorf("could not find order")
	}

	zoneBinding, ok := resp.Results.Bindings[0]["zone"]
	if !ok {
		return "", fmt.Errorf("could not find order zone")
	}
	if zoneBinding.Type != "uri" {
		return "", fmt.Errorf("could not find order zone")
	}

	zoneIRI, ok := zoneBinding.Value.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("could not decode order select")
	}

	var iri IRI
	if err := mapstructure.Decode(zoneIRI, &iri); err != nil {
		return "", fmt.Errorf("could not decode order select: %w", err)
	}

	return iri.Full, nil
}

func (c *Client) GetResourceGovCode(ctx context.Context, resource string) (string, error) {
	resp, err := c.queryCognitariumSelect(ctx, Select{
		Query: SelectQuery{
			Prefixes: []Prefix{{
				Prefix:    "gov",
				Namespace: "https://w3id.org/okp4/ontology/v3/schema/credential/governance/text/",
			}},
			Select: []SelectItem{{Variable: "code"}},
			Where: []WhereCondition{
				{
					Simple: SimpleWhereCondition{TriplePattern: TriplePattern{
						Subject:   VarOrNode{Variable: "credId"},
						Predicate: VarOrNamedNode{NamedNode: &IRI{Full: "dataverse:credential#subject"}},
						Object:    VarOrNodeOrLiteral{Node: &Node{NamedNode: &IRI{Full: resource}}},
					}},
				},
				{
					Simple: SimpleWhereCondition{TriplePattern: TriplePattern{
						Subject:   VarOrNode{Variable: "credId"},
						Predicate: VarOrNamedNode{NamedNode: &IRI{Full: "dataverse:credential#type"}},
						Object:    VarOrNodeOrLiteral{Node: &Node{NamedNode: &IRI{Prefixed: "gov:GovernanceTextCredential"}}},
					}},
				},
				{
					Simple: SimpleWhereCondition{TriplePattern: TriplePattern{
						Subject:   VarOrNode{Variable: "credId"},
						Predicate: VarOrNamedNode{NamedNode: &IRI{Full: "dataverse:credential#claim"}},
						Object:    VarOrNodeOrLiteral{Variable: "claim"},
					}},
				},
				{
					Simple: SimpleWhereCondition{TriplePattern: TriplePattern{
						Subject:   VarOrNode{Variable: "claim"},
						Predicate: VarOrNamedNode{NamedNode: &IRI{Prefixed: "gov:isGovernedBy"}},
						Object:    VarOrNodeOrLiteral{Variable: "gov"},
					}},
				},
				{
					Simple: SimpleWhereCondition{TriplePattern: TriplePattern{
						Subject:   VarOrNode{Variable: "gov"},
						Predicate: VarOrNamedNode{NamedNode: &IRI{Prefixed: "gov:fromGovernance"}},
						Object:    VarOrNodeOrLiteral{Variable: "code"},
					}},
				},
			},
			Limit: 1,
		},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Results.Bindings) != 1 {
		return "", fmt.Errorf("could not find governance code")
	}

	codeBinding, ok := resp.Results.Bindings[0]["code"]
	if !ok {
		return "", fmt.Errorf("could not find governance code")
	}
	if codeBinding.Type != "uri" {
		return "", fmt.Errorf("could not find governance code")
	}

	codeIRI, ok := codeBinding.Value.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("could not decode governance select")
	}

	var iri IRI
	if err := mapstructure.Decode(codeIRI, &iri); err != nil {
		return "", fmt.Errorf("could not decode governance select: %w", err)
	}

	return iri.Full, nil
}

func (c *Client) ExecGovernance(ctx context.Context, govCode, action, subject, zone string) (*struct {
	Result   string
	Evidence string
}, error,
) {
	program, err := makeGovCheckProgram(govCode, action, subject, zone)
	if err != nil {
		return nil, err
	}

	resp, err := c.logicClient.Ask(ctx, &logictypes.QueryServiceAskRequest{
		Program: program,
		Query:   "tell(Result, Evidence).",
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Answer.Results) != 1 {
		return nil, fmt.Errorf("no result")
	}

	resolveVar := func(name string) *string {
		for _, substitution := range resp.Answer.Results[0].Substitutions {
			if substitution.Variable == name {
				return &substitution.Expression
			}
		}
		return nil
	}

	result := resolveVar("Result")
	evidence := resolveVar("Evidence")
	if result == nil || evidence == nil {
		return nil, fmt.Errorf("couldn't resolve variables")
	}

	return &struct {
		Result   string
		Evidence string
	}{Result: *result, Evidence: *evidence}, nil
}

func queryCognitariumAddr(ctx context.Context, dataverseAddr string, wasmClient wasmtypes.QueryClient) (string, error) {
	query, err := json.Marshal(map[string]interface{}{
		"dataverse": struct{}{},
	})
	if err != nil {
		return "", err
	}

	resp, err := wasmClient.SmartContractState(ctx, &wasmtypes.QuerySmartContractStateRequest{
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

func (c *Client) queryCognitariumSelect(ctx context.Context, q Select) (*SelectResponse, error) {
	query, err := json.Marshal(map[string]interface{}{
		"select": q,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.wasmClient.SmartContractState(ctx, &wasmtypes.QuerySmartContractStateRequest{
		Address:   c.cognitariumAddr,
		QueryData: query,
	})
	if err != nil {
		return nil, err
	}

	var res SelectResponse
	if err := json.Unmarshal(resp.Data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
