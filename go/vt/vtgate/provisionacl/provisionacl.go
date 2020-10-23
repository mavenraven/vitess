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

package provisionacl

import (
	"flag"
	"strings"

	querypb "vitess.io/vitess/go/vt/proto/query"
)

var (
	// ProvisionAuthorizedUsers specifies the users that can do provisioning operations via DDL.
	ProvisionAuthorizedUsers = flag.String("provisioner_authorized_users", "", "List of users authorized to run provisioning operations via DDL, or '%' to allow all users.")

	// allowAll is true if the special value of "*" was specified
	allowAll bool

	// aCL contains a set of allowed usernames
	acl map[string]struct{}
)

// Init parses the users option and sets allowAll / acl accordingly
func Init() {
	acl = make(map[string]struct{})
	allowAll = false

	if *ProvisionAuthorizedUsers == "%" {
		allowAll = true
		return
	} else if *ProvisionAuthorizedUsers == "" {
		return
	}

	for _, user := range strings.Split(*ProvisionAuthorizedUsers, ",") {
		user = strings.TrimSpace(user)
		acl[user] = struct{}{}
	}
}

// Authorized returns true if the given caller is allowed to execute vschema operations
func Authorized(caller *querypb.VTGateCallerID) bool {
	if allowAll {
		return true
	}

	user := caller.GetUsername()
	_, ok := acl[user]
	return ok
}
