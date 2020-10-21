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
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"vitess.io/vitess/go/mysql"
	"vitess.io/vitess/go/test/endtoend/cluster"
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
)

func TestMain(m *testing.M) {
	defer cluster.PanicHandler(nil)
	flag.Parse()

	exitCode := func() int {
		clusterForProvisionTest = cluster.NewCluster(cell, hostname)
		clusterForProvisionTest.VtGateExtraArgs = []string {
			"-provision_authorized_users",
			"%",
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

func TestCreateKeyspace(t *testing.T) {
	defer cluster.PanicHandler(t)

	ctx := context.Background()
	vtParams := mysql.ConnParams{
		Host: clusterForProvisionTest.Hostname,
		Port: clusterForProvisionTest.VtgateMySQLPort,
	}
	conn, err := mysql.Connect(ctx, &vtParams)
	require.Nil(t, err)

	qr, err := conn.ExecuteFetch("select * from information_schema.tables;", 1000, true)

	fmt.Printf("%+v\n", qr)

	//output, err := clusterForProvisionTest.VtctlclientProcess.ExecuteCommandWithOutput("GetSrvKeyspaceNames", cell)
	//require.Nil(t, err)
//	assert.Contains(t, strings.Split(output, "\n"), keyspaceUnshardedName)
//	assert.Contains(t, strings.Split(output, "\n"), keyspaceShardedName)
}

