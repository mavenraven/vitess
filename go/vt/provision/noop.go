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

package provision

import (
	"fmt"
)

type noopProvisioner struct{}

func newNoopProvisioner(config map[string]string) (Provisioner, error){
	return noopProvisioner{}, nil
}

func (noopProvisioner) RequestCreateKeyspace(keyspace string) error {
	//FIXME: better error
	return fmt.Errorf("not implemented")
}

func (noopProvisioner) RequestDeleteKeyspace(keyspace string) error {
	//FIXME: better error
	return fmt.Errorf("not implemented")
}

func init() {
	provisioners["noop"] = newNoopProvisioner
}

