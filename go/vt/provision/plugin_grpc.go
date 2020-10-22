package provision

import (
	"context"
	"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"time"
	"vitess.io/vitess/go/vt/log"
	"vitess.io/vitess/go/vt/proto/provision"
)

var (
	errNeedGrpcEndpoint = fmt.Errorf("need grpc endpoint to use grpc provisioning")
	//FIXME: underscores or dashes
	grpcEndpoint = "provision_grpc_endpoint"
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

	log.Errorf("before dial")

	//FIXME: cli option for endpont
	dialTimeout, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	//FIXME: tls
	conn, err := grpc.DialContext(dialTimeout, "localhost:9696", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		//FIXME: better error
		fmt.Errorf("dialing to provisioner timed out")
	}
	defer conn.Close()

	client := provision.NewProvisionClient(conn)
	req := &provision.RequestCreateKeyspaceRequest{
		Keyspace:             keyspace,
	}

	pe, err := client.RequestCreateKeyspace(
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

	if err != nil {
		return err
	}

	switch pe.Code {
	case provision.Code_OK:
		return nil
	case provision.Code_UNKNOWN:
		//FIXME: better error
		return fmt.Errorf("unknown error")
	default:
		//FIXME: better error
		return fmt.Errorf("unhandled grpc case")
	}
}

func (p *grpcProvisioner) RequestDeleteKeyspace(ctx context.Context, keyspace string) error {
	return nil
}

func init() {
	provisioners["grpc"] = newGRPCProvisioner
}