package provision

import (
	"flag"
	"fmt"
)

var (
	//FIXME: underscores or dashes?
	grpcEndpoint        = "provisioner_grpc_endpoint"
	errNeedGrpcEndpoint = fmt.Errorf("need grpc endpoint to use grpc provisioning")
)


type grpcProvisioner struct {}

func newGRPCProvisioner(config map[string]string) (Provisioner, error){
	grpcEndpointConfig, ok := config[grpcEndpoint]
	if !ok {
		return nil, errNeedGrpcEndpoint
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
	provisioners["grpc"] = newGRPCProvisioner
	flags[grpcEndpoint] = *flag.String(
		grpcEndpoint,
		"",
		"Endpoint to send provisioning requests.",
	)
}