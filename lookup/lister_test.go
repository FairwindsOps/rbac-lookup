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
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"google.golang.org/api/cloudresourcemanager/v1"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestLoadRoleBindings(t *testing.T) {
	l := genLister()

	loadRoleBindings(t, l)

	assert.Len(t, l.rbacSubjectsByScope, 0, "Expected no rbac subjects initially")

	createRoleBindings(t, l)

	loadRoleBindings(t, l)

	assert.Len(t, l.rbacSubjectsByScope, 2, "Expected 2 rbac subjects")

	expectedRbacSubject := rbacSubject{
		Kind: "User",
		RolesByScope: map[string][]simpleRole{
			"foo": {{
				Kind: "Role",
				Name: "bar",
				Source: simpleRoleSource{
					Kind: "RoleBinding",
					Name: "testing",
				},
			}},
		},
	}

	assert.EqualValues(t, l.rbacSubjectsByScope["joe"], expectedRbacSubject)
	assert.EqualValues(t, l.rbacSubjectsByScope["sue"], expectedRbacSubject)
}

func TestLoadClusterRoleBindings(t *testing.T) {
	l := genLister()

	loadClusterRoleBindings(t, l)

	assert.Len(t, l.rbacSubjectsByScope, 0, "Expected no rbac subjects initially")

	createClusterRoleBindings(t, l)

	loadClusterRoleBindings(t, l)

	assert.Len(t, l.rbacSubjectsByScope, 2, "Expected 2 rbac subjects")

	expectedRbacSubject := rbacSubject{
		Kind: "User",
		RolesByScope: map[string][]simpleRole{
			"cluster-wide": {{
				Kind: "ClusterRole",
				Name: "bar",
				Source: simpleRoleSource{
					Kind: "ClusterRoleBinding",
					Name: "testing",
				},
			}},
		},
	}

	assert.EqualValues(t, l.rbacSubjectsByScope["joe"], expectedRbacSubject)
	assert.EqualValues(t, l.rbacSubjectsByScope["sue"], expectedRbacSubject)
}

func TestLoadAll(t *testing.T) {
	l := genLister()

	loadAll(t, l)

	assert.Len(t, l.rbacSubjectsByScope, 0, "Expected no rbac subjects initially")

	createRoleBindings(t, l)

	createClusterRoleBindings(t, l)

	loadAll(t, l)

	assert.Len(t, l.rbacSubjectsByScope, 2, "Expected 2 rbac subjects")

	expectedRbacSubject := rbacSubject{
		Kind: "User",
		RolesByScope: map[string][]simpleRole{
			"cluster-wide": {{
				Kind: "ClusterRole",
				Name: "bar",
				Source: simpleRoleSource{
					Kind: "ClusterRoleBinding",
					Name: "testing",
				},
			}},
			"foo": {{
				Kind: "Role",
				Name: "bar",
				Source: simpleRoleSource{
					Kind: "RoleBinding",
					Name: "testing",
				},
			}},
		},
	}

	assert.EqualValues(t, l.rbacSubjectsByScope["joe"], expectedRbacSubject)
	assert.EqualValues(t, l.rbacSubjectsByScope["sue"], expectedRbacSubject)
}

func TestLoadGke(t *testing.T) {
	policy := &cloudresourcemanager.Policy{
		Bindings: []*cloudresourcemanager.Binding{{
			Role:    "roles/container.admin",
			Members: []string{"user:jane@example.com", "user:joe@example.com"},
		}, {
			Role:    "roles/container.developer",
			Members: []string{"serviceAccount:ci@example.iam.gserviceaccount.com"},
		}, {
			Role:    "roles/viewer",
			Members: []string{"group:devs@example.com"},
		}, {
			Role:    "roles/owner",
			Members: []string{"user:jane@example.com"},
		}},
	}

	l := genLister()

	assert.Len(t, l.rbacSubjectsByScope, 0, "Expected no rbac subjects initially")

	l.loadGkeIamPolicy(policy)

	assert.Len(t, l.rbacSubjectsByScope, 4, "Expected 4 rbac subjects")

	assert.EqualValues(t, l.rbacSubjectsByScope["jane@example.com"], rbacSubject{
		Kind: "User",
		RolesByScope: map[string][]simpleRole{
			"project-wide": {{
				Kind: "IAM",
				Name: "gke-admin",
				Source: simpleRoleSource{
					Kind: "IAMRole",
					Name: "container.admin",
				},
			}, {
				Kind: "IAM",
				Name: "gcp-owner",
				Source: simpleRoleSource{
					Kind: "IAMRole",
					Name: "owner",
				},
			}},
		},
	})

	assert.EqualValues(t, l.rbacSubjectsByScope["joe@example.com"], rbacSubject{
		Kind: "User",
		RolesByScope: map[string][]simpleRole{
			"project-wide": {{
				Kind: "IAM",
				Name: "gke-admin",
				Source: simpleRoleSource{
					Kind: "IAMRole",
					Name: "container.admin",
				},
			}},
		},
	})

	assert.EqualValues(t, l.rbacSubjectsByScope["devs@example.com"], rbacSubject{
		Kind: "Group",
		RolesByScope: map[string][]simpleRole{
			"project-wide": {{
				Kind: "IAM",
				Name: "gcp-viewer",
				Source: simpleRoleSource{
					Kind: "IAMRole",
					Name: "viewer",
				},
			}},
		},
	})

	assert.EqualValues(t, l.rbacSubjectsByScope["ci@example.iam.gserviceaccount.com"], rbacSubject{
		Kind: "ServiceAccount",
		RolesByScope: map[string][]simpleRole{
			"project-wide": {{
				Kind: "IAM",
				Name: "gke-developer",
				Source: simpleRoleSource{
					Kind: "IAMRole",
					Name: "container.developer",
				},
			}},
		},
	})
}

func TestLoadGkeFilters(t *testing.T) {
	policy := &cloudresourcemanager.Policy{
		Bindings: []*cloudresourcemanager.Binding{{
			Role:    "roles/container.admin",
			Members: []string{"user:jane@example.com", "user:joe@example.com"},
		}, {
			Role:    "roles/container.developer",
			Members: []string{"serviceAccount:ci@example.iam.gserviceaccount.com"},
		}, {
			Role:    "roles/viewer",
			Members: []string{"group:devs@example.com"},
		}, {
			Role:    "roles/owner",
			Members: []string{"user:jane@example.com"},
		}},
	}

	l := genLister()
	l.filter = "example"
	l.subjectKind = "user"

	assert.Len(t, l.rbacSubjectsByScope, 0, "Expected no rbac subjects initially")

	l.loadGkeIamPolicy(policy)

	assert.Len(t, l.rbacSubjectsByScope, 2, "Expected 2 rbac subjects")

	assert.EqualValues(t, l.rbacSubjectsByScope["jane@example.com"], rbacSubject{
		Kind: "User",
		RolesByScope: map[string][]simpleRole{
			"project-wide": {{
				Kind: "IAM",
				Name: "gke-admin",
				Source: simpleRoleSource{
					Kind: "IAMRole",
					Name: "container.admin",
				},
			}, {
				Kind: "IAM",
				Name: "gcp-owner",
				Source: simpleRoleSource{
					Kind: "IAMRole",
					Name: "owner",
				},
			}},
		},
	})

	assert.EqualValues(t, l.rbacSubjectsByScope["joe@example.com"], rbacSubject{
		Kind: "User",
		RolesByScope: map[string][]simpleRole{
			"project-wide": {{
				Kind: "IAM",
				Name: "gke-admin",
				Source: simpleRoleSource{
					Kind: "IAMRole",
					Name: "container.admin",
				},
			}},
		},
	})
}

func genLister() lister {
	return lister{
		clientset:           testclient.NewSimpleClientset(),
		rbacSubjectsByScope: make(map[string]rbacSubject),
	}
}

func loadAll(t *testing.T, l lister) {
	err := l.loadAll()

	assert.Nil(t, err, "Expected no error loading all rbac Bindings")
}

func loadRoleBindings(t *testing.T, l lister) {
	err := l.loadRoleBindings()

	assert.Nil(t, err, "Expected no error loading role bindings")
}

func createRoleBindings(t *testing.T, l lister) {
	roleBindings := []rbacv1.RoleBinding{{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testing",
			Namespace: "foo",
		},
		Subjects: []rbacv1.Subject{{
			Name: "joe",
			Kind: "User",
		}, {
			Name: "sue",
			Kind: "User",
		}},
		RoleRef: rbacv1.RoleRef{
			Kind: "Role",
			Name: "bar",
		},
	}}

	for _, roleBinding := range roleBindings {
		_, err := l.clientset.RbacV1().RoleBindings(roleBinding.Namespace).Create(context.Background(), &roleBinding, metav1.CreateOptions{})
		assert.Nil(t, err, "Expected no error creating role bindings")
	}
}

func loadClusterRoleBindings(t *testing.T, l lister) {
	err := l.loadClusterRoleBindings()

	assert.Nil(t, err, "Expected no error loading cluster role bindings")
}

func createClusterRoleBindings(t *testing.T, l lister) {
	clusterRoleBindings := []rbacv1.ClusterRoleBinding{{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testing",
		},
		Subjects: []rbacv1.Subject{{
			Name: "joe",
			Kind: "User",
		}, {
			Name: "sue",
			Kind: "User",
		}},
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole",
			Name: "bar",
		},
	}}

	for _, clusterRoleBinding := range clusterRoleBindings {
		_, err := l.clientset.RbacV1().ClusterRoleBindings().Create(context.Background(), &clusterRoleBinding, metav1.CreateOptions{})
		assert.Nil(t, err, "Expected no error creating cluster role bindings")
	}
}
