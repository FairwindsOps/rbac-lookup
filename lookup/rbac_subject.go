// Copyright 2018 FairwindsOps Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lookup

import (
	rbacv1 "k8s.io/api/rbac/v1"
)

type rbacSubject struct {
	Kind         string
	RolesByScope map[string][]simpleRole
}

type simpleRole struct {
	Kind   string
	Name   string
	Source simpleRoleSource
}

type simpleRoleSource struct {
	Kind string
	Name string
}

func (rbacSubj *rbacSubject) addRoleBinding(roleBinding *rbacv1.RoleBinding) {
	simpleRole := simpleRole{
		Name: roleBinding.RoleRef.Name,
		Source: simpleRoleSource{
			Name: roleBinding.Name,
			Kind: "RoleBinding",
		},
	}

	simpleRole.Kind = roleBinding.RoleRef.Kind
	rbacSubj.RolesByScope[roleBinding.Namespace] = append(rbacSubj.RolesByScope[roleBinding.Namespace], simpleRole)
}

func (rbacSubj *rbacSubject) addClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) {
	simpleRole := simpleRole{
		Name:   clusterRoleBinding.RoleRef.Name,
		Source: simpleRoleSource{Name: clusterRoleBinding.Name, Kind: "ClusterRoleBinding"},
	}

	simpleRole.Kind = clusterRoleBinding.RoleRef.Kind
	scope := "cluster-wide"
	rbacSubj.RolesByScope[scope] = append(rbacSubj.RolesByScope[scope], simpleRole)
}
