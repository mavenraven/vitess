package provision

import (
	"context"
	"flag"
	"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"time"
	"vitess.io/vitess/go/vt/proto/provision"
	vtrpcpb "vitess.io/vitess/go/vt/proto/vtrpc"
	"vitess.io/vitess/go/vt/vterrors"
)


var (
	ErrNeedGrpcEndpoint = vterrors.Errorf(
		vtrpcpb.Code_FAILED_PRECONDITION,
		"need grpc endpoint to use grpc provisioning",
	)

	provisionGrpcEndpoint = flag.String("provisioner_grpc_endpoint", "", "")
	provisionGrpcDialTimeout = flag.Duration("provisioner_grpc_dial_timeout", time.Duration(5 * time.Second), "")
	provisionGrpcRequestTimeout = flag.Duration("provisioner_grpc_per_retry_timeout", time.Duration(5 * time.Second), "")
	provisionGrpcMaxRetries = flag.Uint("provisioner_grpc_max_retries", 3, "")
)
type grpcProvisioner struct {}

func (p grpcProvisioner) RequestCreateKeyspace(ctx context.Context, keyspace string) error {
	if *provisionGrpcEndpoint == "" {
		return ErrNeedGrpcEndpoint
	}
	dialTimeout, cancel := context.WithTimeout(ctx, *provisionGrpcDialTimeout)
	defer cancel()

	//FIXME: tls
	conn, err := grpc.DialContext(dialTimeout, *provisionGrpcEndpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		vterrors.Wrapf(err, "dialing to provisioner timed out")
	}
	defer conn.Close()

	client := provision.NewProvisionClient(conn)
	req := &provision.RequestCreateKeyspaceRequest{
		Keyspace:             keyspace,
	}

	_, err = client.RequestCreateKeyspace(
		ctx,
		req,
		grpc_retry.WithPerRetryTimeout(*provisionGrpcRequestTimeout),
		grpc_retry.WithMax(*provisionGrpcMaxRetries),
		grpc_retry.WithBackoff(
			grpc_retry.BackoffLinear(1 * time.Second),
		),
	)

	return err
}

func (p grpcProvisioner) RequestDeleteKeyspace(ctx context.Context, keyspace string) error {
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
	provisioners["grpc"] = grpcProvisioner{}
}