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

var gkeIamScope = "project-wide"
var gkeIamRoles = map[string]simpleRole{
	"roles/container.clusterAdmin": {
		Kind: "IAM",
		Name: "gke-cluster-admin",
		Source: simpleRoleSource{
			Kind: "IAMRole",
			Name: "container.clusterAdmin",
		},
	},
	"roles/container.admin": {
		Kind: "IAM",
		Name: "gke-admin",
		Source: simpleRoleSource{
			Kind: "IAMRole",
			Name: "container.admin",
		},
	},
	"roles/container.developer": {
		Kind: "IAM",
		Name: "gke-developer",
		Source: simpleRoleSource{
			Kind: "IAMRole",
			Name: "container.developer",
		},
	},
	"roles/container.viewer": {
		Kind: "IAM",
		Name: "gke-viewer",
		Source: simpleRoleSource{
			Kind: "IAMRole",
			Name: "container.viewer",
		},
	},
	"roles/owner": {
		Kind: "IAM",
		Name: "gcp-owner",
		Source: simpleRoleSource{
			Kind: "IAMRole",
			Name: "owner",
		},
	},
	"roles/admin": {
		Kind: "IAM",
		Name: "gcp-admin",
		Source: simpleRoleSource{
			Kind: "IAMRole",
			Name: "admin",
		},
	},
	"roles/editor": {
		Kind: "IAM",
		Name: "gcp-editor",
		Source: simpleRoleSource{
			Kind: "IAMRole",
			Name: "editor",
		},
	},
	"roles/viewer": {
		Kind: "IAM",
		Name: "gcp-viewer",
		Source: simpleRoleSource{
			Kind: "IAMRole",
			Name: "viewer",
		},
	},
}
