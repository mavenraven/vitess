package provision

import (
	"flag"
	"fmt"
)

var (
	//FIXME: underscores or dashes?
	grpcEndpointKey = "provisioner_grpc_endpoint"
	grpcEndpointValue = flag.String(grpcEndpointKey, "", "Endpoint to send provisioning requests.")
	errNeedGrpcEndpoint = fmt.Errorf("need grpc endpoint to use grpc provisioning")
)


type grpcProvisioner struct {}

func newGRPCProvisioner(config map[string]string) (*grpcProvisioner, error){
	grpcEndpointConfig, ok := config[grpcEndpointKey]
	if !ok {
		grpcEndpointConfig = *grpcEndpointValue
	}

	if grpcEndpointConfig == "" {
		return nil, errNeedGrpcEndpoint
	}

	return &grpcProvisioner{}, nil
}

func (p *grpcProvisioner) CreateKeyspace(keyspace string) error {
	return nil
}

func (p *grpcProvisioner) DeleteKeyspace(keyspace string) error {
	return nil
}

func init() {
	flagsConfig[grpcEndpointKey] = *grpcEndpointValue
}