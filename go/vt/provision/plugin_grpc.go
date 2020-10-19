package provision

import (
	"fmt"
)

var (
	errNeedGrpcEndpoint = fmt.Errorf("need grpc endpoint to use grpc provisioning")
	grpcEndpoint = "provisioner_grpc_endpoint"
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

func (p *grpcProvisioner) RequestCreateKeyspace(keyspace string) error {
	return nil
}

func (p *grpcProvisioner) RequestDeleteKeyspace(keyspace string) error {
	return nil
}

func init() {
	provisioners["grpc"] = newGRPCProvisioner
}