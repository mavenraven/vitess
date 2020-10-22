/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sequence

import (
	"context"
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"net"
	"os"
	"testing"
	"time"
	"vitess.io/vitess/go/mysql"
	"vitess.io/vitess/go/test/endtoend/cluster"
	"vitess.io/vitess/go/vt/log"
	"vitess.io/vitess/go/vt/proto/provision"
	"vitess.io/vitess/go/vt/proto/topodata"
)

var (
	clusterForProvisionTest *cluster.LocalProcessCluster
	keyspaceShardedName     = "test_ks_sharded"
	keyspaceUnshardedName   = "test_ks_unsharded"
	cell                    = "zone1"
	cell2                   = "zone2"
	hostname                = "localhost"
	servedTypes             = map[topodata.TabletType]bool{topodata.TabletType_MASTER: true, topodata.TabletType_REPLICA: true, topodata.TabletType_RDONLY: true}
	sqlSchema               = `create table vt_insert_test (
								id bigint auto_increment,
								msg varchar(64),
								keyspace_id bigint(20) unsigned NOT NULL,
								primary key (id)
								) Engine=InnoDB`
	grpcServerAddress net.Addr
)

func TestMain(m *testing.M) {
	defer cluster.PanicHandler(nil)
	flag.Parse()

	exitCode := func() int {

		addrChan := make(chan net.Addr)
		defer close(addrChan)

		errorChan := make(chan error)
		defer close(errorChan)

		go func() {
			err := startGrpcServer(context.TODO(), addrChan)
			if err != nil {
				errorChan <- err
			}
		}()

		select {
		case err := <-errorChan:
			log.Error(err)
			return 1
		case grpcServerAddress = <-addrChan:
		}

		clusterForProvisionTest = cluster.NewCluster(cell, hostname)
		clusterForProvisionTest.VtGateExtraArgs = []string {
			"-provision_authorized_users",
			"%",
			"-provision_type",
			"grpc",
			"-provision_grpc_endpoint",
			grpcServerAddress.String(),

		}

		defer clusterForProvisionTest.Teardown()

		// Start topo server
		if err := clusterForProvisionTest.StartTopo(); err != nil {
			return 1
		}

		if err := clusterForProvisionTest.TopoProcess.ManageTopoDir("mkdir", "/vitess/"+cell2); err != nil {
			return 1
		}

		if err := clusterForProvisionTest.VtctlProcess.AddCellInfo(cell2); err != nil {
			return 1
		}


		keyspaceUnsharded := &cluster.Keyspace{
			Name:      keyspaceUnshardedName,
			SchemaSQL: sqlSchema,
		}
		if err := clusterForProvisionTest.StartKeyspace(*keyspaceUnsharded, []string{keyspaceUnshardedName}, 1, false); err != nil {
			return 1
		}
		if err := clusterForProvisionTest.VtctlclientProcess.ExecuteCommand("SetKeyspaceShardingInfo", "-force", keyspaceUnshardedName, "keyspace_id", "uint64"); err != nil {
			return 1
		}
		if err := clusterForProvisionTest.VtctlclientProcess.ExecuteCommand("RebuildKeyspaceGraph", keyspaceUnshardedName); err != nil {
			return 1
		}

		// Start vtgate
		if err := clusterForProvisionTest.StartVtgate(); err != nil {
			return 1
		}


		return m.Run()
	}()
	os.Exit(exitCode)
}

func TestProvisionKeyspace(t *testing.T) {
	fmt.Println(grpcServerAddress)
	defer cluster.PanicHandler(t)

	ctx := context.Background()
	vtParams := mysql.ConnParams{
		Host: clusterForProvisionTest.Hostname,
		Port: clusterForProvisionTest.VtgateMySQLPort,
		ConnectTimeoutMs: 1000,
	}
	conn, err := mysql.Connect(ctx, &vtParams)
	require.Nil(t, err)

	log.Info(clusterForProvisionTest.Keyspaces)
	qr, err := conn.ExecuteFetch("CREATE DATABASE my_keyspace;", 10, true)
	require.Nil(t, err)

	assert.Equal(t, uint64(1), qr.RowsAffected, "got the following back from vtgate instead: %v", qr.Rows)

	_, err = clusterForProvisionTest.VtctlclientProcess.ExecuteCommandWithOutput("GetKeyspace", "my_keyspace")
	//If GetKeyspace doesn't return an error, the keyspace exists.
	require.Nil(t, err)
}

type testGrpcServer struct {}

func (_ testGrpcServer)RequestCreateKeyspace(ctx context.Context, rckr *provision.RequestCreateKeyspaceRequest) (*provision.ProvisionError, error) {
	log.Info("got request for keyspace: " + rckr.Keyspace)
	//We're doing this in a go routine to simulate the fact that RequestCreateKeyspace does not block while the
	//the keyspace is being created.
	go func() {
		<- time.After(10 * time.Second)
		err := clusterForProvisionTest.VtctlProcess.CreateKeyspace(rckr.Keyspace)
		if err != nil {
			log.Error(err)
		}
	}()
	return &provision.ProvisionError{Code: provision.Code_OK}, nil
}


func startGrpcServer(ctx context.Context, addr chan net.Addr) error {
	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", "localhost:")
	if err != nil {
		return err
	}

	defer listener.Close()

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	defer grpcServer.Stop()

	provision.RegisterProvisionServer(grpcServer, testGrpcServer{})
	addr <- listener.Addr()
	return grpcServer.Serve(listener)
}