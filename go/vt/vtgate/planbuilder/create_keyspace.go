package planbuilder

import (
	"vitess.io/vitess/go/vt/vtgate/engine"
)

func buildCreateKeyspacePlan(keyspaceName string, ifNotExists bool) engine.Primitive {
	return &engine.CreateKeyspace{
		RequestedKeyspace: keyspaceName,
		IfNotExists: ifNotExists,
	}
}

