package provision

import (
	"vitess.io/vitess/go/vt/log"
	vtrpcpb "vitess.io/vitess/go/vt/proto/vtrpc"
	"vitess.io/vitess/go/vt/vterrors"
)

var (
	provisioners = make(map[string]Provisioner)
)

func factory(provisionerType string) Provisioner {
	p, ok := provisioners[provisionerType]
	if !ok {
		log.Error(vterrors.Errorf(
			vtrpcpb.Code_INVALID_ARGUMENT,
			"failed to find %s provisioner, defaulting to noop",
			provisionerType,
			))
		return noopProvisioner{}
	}
	return p
}
