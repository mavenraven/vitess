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

// Package trace contains a helper interface that allows various tracing
// tools to be plugged in to components using this interface. If no plugin is
// registered, the default one makes all trace calls into no-ops.
package provision

import (
	"flag"
	"fmt"
	"vitess.io/vitess/go/vt/log"
	"vitess.io/vitess/go/vt/vterrors"
)

var (
	//FIXME: _ or -
	provisionerType = "provisioner_type"
	ErrInvalidProvisionerType = fmt.Errorf("provisionerType not found")
	flags                     = make (map[string]string)
)

type Provisioner interface {
	CreateKeyspace(keyspace string) error
	DeleteKeyspace(keyspace string) error
}

func CreateKeyspace(keyspace string) error {
	p, err := NewProvisioner(flags[provisionerType], flags)
	if err != nil {
		log.Error(vterrors.Wrapf(err, "failed to find %s provisioner, defaulting to noop", flags[provisionerType]))
		p = noopProvisioner{}
	}
	return p.CreateKeyspace(keyspace)
}

func DeleteKeyspace(keyspace string) error {
	p, err := NewProvisioner(flags[provisionerType], flags)
	if err != nil {
		log.Error(vterrors.Wrapf(err, "failed to find %s provisioner, defaulting to noop", flags[provisionerType]))
		p = noopProvisioner{}
	}
	return p.DeleteKeyspace(keyspace)
}

func init() {
	 flags[provisionerType] = *flag.String(provisionerType, "noop", "Provisioner type to use")
}

