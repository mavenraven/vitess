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

package main

// Imports and register the gRPC vtgateservice server

import (
	"context"
	"vitess.io/vitess/go/vt/vtgate/grpcvtgateconn"
	_ "vitess.io/vitess/go/vt/vtgate/grpcvtgateservice"
	"vitess.io/vitess/go/vt/vtgate/vtgateconn"
)

func init() {
	vtgateconn.RegisterDialer("grpc", func(ctx context.Context, addr string) (vtgateconn.Impl, error) {
		return grpcvtgateconn.DialWithOpts(ctx)(ctx, addr)
	})
	/*
	fmt.Printf("WOW init for pid %v", os.Getpid())
	vtgate.RegisterVTGates = append(vtgate.RegisterVTGates, func(vtGate vtgateservice.VTGateService) {
		if servenv.GRPCCheckServiceMap("vtgateservice") {
			fmt.Printf("WOW registering vtgateservice for pid %v", os.Getpid())
			vtgateservicepb.RegisterVitessServer(servenv.GRPCServer, grpcvtgateservice.NewVTGate(vtGate))
		}
	})
	 */

}
