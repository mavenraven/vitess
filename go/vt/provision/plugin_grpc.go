package provision

import (
	"context"
	"flag"
	"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"time"
	"vitess.io/vitess/go/vt/proto/provision"
)

var (
	errNeedGrpcEndpoint = fmt.Errorf("need grpc endpoint to use grpc provisioning")
	//FIXME: underscores or dashes
	provisionGrpcEndpoint = flag.String("provision_grpc_endpoint", "", "Endpoint for gRPC server.")
)
type grpcProvisioner struct {}

func newGRPCProvisioner(config map[string]string) (Provisioner, error){
	//FIXME: skip this for now
	/*
	grpcEndpointConfig, ok := config[grpcEndpoint]
	if !ok {
		return nil, errNeedGrpcEndpoint
	}

	if grpcEndpointConfig == "" {
		return nil, errNeedGrpcEndpoint
	}
	 */

	return &grpcProvisioner{}, nil
}

func (p *grpcProvisioner) RequestCreateKeyspace(ctx context.Context, keyspace string) error {
	//FIXME: cli option for endpont
	dialTimeout, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	//FIXME: tls
	conn, err := grpc.DialContext(dialTimeout, *provisionGrpcEndpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		//FIXME: better error
		fmt.Errorf("dialing to provisioner timed out")
	}
	defer conn.Close()

	client := provision.NewProvisionClient(conn)
	req := &provision.RequestCreateKeyspaceRequest{
		Keyspace:             keyspace,
	}

	_, err = client.RequestCreateKeyspace(
		ctx,
		req,
		//FIXME: cli option
		grpc_retry.WithPerRetryTimeout(5 * time.Second),
		//FIXME: cli option
		grpc_retry.WithMax(3),
		grpc_retry.WithBackoff(
			grpc_retry.BackoffLinear(1 * time.Second),
		),
	)

	return err
}

func (p *grpcProvisioner) RequestDeleteKeyspace(ctx context.Context, keyspace string) error {
	//FIXME: cli option for endpont
	dialTimeout, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	//FIXME: tls
	conn, err := grpc.DialContext(dialTimeout, *provisionGrpcEndpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		//FIXME: better error
		fmt.Errorf("dialing to provisioner timed out")
	}
	defer conn.Close()

	client := provision.NewProvisionClient(conn)
	req := &provision.RequestDeleteKeyspaceRequest{
		Keyspace:             keyspace,
	}

	_, err = client.RequestDeleteKeyspace(
		ctx,
		req,
		//FIXME: cli option
		grpc_retry.WithPerRetryTimeout(5 * time.Second),
		//FIXME: cli option
		grpc_retry.WithMax(3),
		grpc_retry.WithBackoff(
			grpc_retry.BackoffLinear(1 * time.Second),
		),
	)

	return err
}

func init() {
	provisioners["grpc"] = newGRPCProvisioner
}