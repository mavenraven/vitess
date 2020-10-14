/*
Copyright 2020 The Vitess Authors.

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

package engine

import (
	"vitess.io/vitess/go/sqltypes"
	"vitess.io/vitess/go/vt/proto/query"
	querypb "vitess.io/vitess/go/vt/proto/query"
	vtrpcpb "vitess.io/vitess/go/vt/proto/vtrpc"
	"vitess.io/vitess/go/vt/schema"
	"vitess.io/vitess/go/vt/sqlparser"
	"vitess.io/vitess/go/vt/vterrors"
	"vitess.io/vitess/go/vt/vtgate/vindexes"
)

var _ Primitive = (*CreateDatabaseDDL)(nil)

//CreateDatabase DDL represents the instructions to create a new keyspace via the user issuing a "CREATE DATABASE" type
//statement. As the actual creation logic is outside of the scope of vitess, the request is submitted to a service.
type CreateDatabaseDDL struct {
	RequestedKeyspace *vindexes.Keyspace
	noTxNeeded

	noInputs
}

func (v *CreateDatabaseDDL) description() PrimitiveDescription {
	return PrimitiveDescription{
		OperatorType: "CreateDatabaseDDL",
		Keyspace:     nil,
		Other: map[string]interface{}{
			"query": sqlparser.String(v.DDL),
		},
	}
}

//RouteType implements the Primitive interface
func (v *CreateDatabaseDDL) RouteType() string {
	return "CreateDatabaseDDL"
}

//GetKeyspaceName implements the Primitive interface
func (v *CreateDatabaseDDL) GetKeyspaceName() string {
	return nil, vterrors.Errorf(vtrpcpb.Code_INTERNAL, "not reachable") // FIXME: david - don't this is reachable
}

//GetTableName implements the Primitive interface
func (v *CreateDatabaseDDL) GetTableName() string {
	return v.DDL.Table.Name.String()
}

//Execute implements the Primitive interface
func (v *CreateDatabaseDDL) Execute(vcursor VCursor, bindVars map[string]*query.BindVariable, wantfields bool) (result *sqltypes.Result, err error) {
	onlineDDL, err := schema.NewOnlineDDL(v.GetKeyspaceName(), v.GetTableName(), v.SQL, v.Strategy, v.Options)
	if err != nil {
		return result, err
	}
	err = vcursor.SubmitOnlineDDL(onlineDDL)
	if err != nil {
		return result, err
	}

	result = &sqltypes.Result{
		Fields: []*querypb.Field{
			{
				Name: "uuid",
				Type: sqltypes.VarChar,
			},
		},
		Rows: [][]sqltypes.Value{
			{
				sqltypes.NewVarChar(onlineDDL.UUID),
			},
		},
		RowsAffected: 1,
	}
	return result, err
}

//StreamExecute implements the Primitive interface
func (v *CreateDatabaseDDL) StreamExecute(vcursor VCursor, bindVars map[string]*query.BindVariable, wantields bool, callback func(*sqltypes.Result) error) error {
	return vterrors.Errorf(vtrpcpb.Code_INTERNAL, "not reachable") // FIXME: david - copied from online_ddl.go, also have no idea if this should work
}

//GetFields implements the Primitive interface
func (v *CreateDatabaseDDL) GetFields(vcursor VCursor, bindVars map[string]*query.BindVariable) (*sqltypes.Result, error) {
	return nil, vterrors.Errorf(vtrpcpb.Code_INTERNAL, "not reachable") // FIXME: david - copied from online_ddl.go, also have no idea if this should work
}
