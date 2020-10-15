package planbuilder

import (
	"vitess.io/vitess/go/vt/vtgate/engine"
)

func buildCreateKeyspacePlan(keyspaceName string, ifExists bool) engine.Primitive {
	return &engine.CreateKeyspace{
		RequestedKeyspace: keyspaceName,
		IfExists: ifExists,
	}
}

