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

	assert.Len(t, l.RbacSubjectsByScope, 0, "Expected no rbac subjects initially")

	createRoleBindings(t, l)

	loadRoleBindings(t, l)

	assert.Len(t, l.RbacSubjectsByScope, 3, "Expected 3 rbac subjects")

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

	expectedRbacSubjectSA := rbacSubject{
		Kind: "ServiceAccount",
		RolesByScope: map[string][]simpleRole{
			"two": {{
				Kind: "ClusterRole",
				Name: "cluster-admin",
				Source: simpleRoleSource{
					Kind: "RoleBinding",
					Name: "testing-sa",
				},
			}},
			"three": {{
				Kind: "ClusterRole",
				Name: "cluster-admin",
				Source: simpleRoleSource{
					Kind: "RoleBinding",
					Name: "testing-sa",
				},
			}},
		},
	}

	assert.EqualValues(t, expectedRbacSubject, l.RbacSubjectsByScope["joe"])
	assert.EqualValues(t, expectedRbacSubject, l.RbacSubjectsByScope["sue"])
	assert.EqualValues(t, expectedRbacSubjectSA, l.RbacSubjectsByScope["circleci:circleci"])
}

func TestLoadClusterRoleBindings(t *testing.T) {
	l := genLister()

	loadClusterRoleBindings(t, l)

	assert.Len(t, l.RbacSubjectsByScope, 0, "Expected no rbac subjects initially")

	createClusterRoleBindings(t, l)

	loadClusterRoleBindings(t, l)

	assert.Len(t, l.RbacSubjectsByScope, 3, "Expected 3 rbac subjects")

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

	expectedRbacSubjectSA := rbacSubject{
		Kind: "ServiceAccount",
		RolesByScope: map[string][]simpleRole{
			"cluster-wide": {{
				Kind: "ClusterRole",
				Name: "cluster-admin",
				Source: simpleRoleSource{
					Kind: "ClusterRoleBinding",
					Name: "circleci-cluster-admin",
				},
			}},
		},
	}

	assert.EqualValues(t, expectedRbacSubject, l.RbacSubjectsByScope["joe"])
	assert.EqualValues(t, expectedRbacSubject, l.RbacSubjectsByScope["sue"])
	assert.EqualValues(t, expectedRbacSubjectSA, l.RbacSubjectsByScope["circleci:circleci"])
}

func TestLoadAll(t *testing.T) {
	l := genLister()

	loadAll(t, l)

	assert.Len(t, l.RbacSubjectsByScope, 0, "Expected no rbac subjects initially")

	createRoleBindings(t, l)

	createClusterRoleBindings(t, l)

	loadAll(t, l)

	assert.Len(t, l.RbacSubjectsByScope, 3, "Expected 3 rbac subjects")

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

	expectedRbacSubjectSA := rbacSubject{
		Kind: "ServiceAccount",
		RolesByScope: map[string][]simpleRole{
			"cluster-wide": {{
				Kind: "ClusterRole",
				Name: "cluster-admin",
				Source: simpleRoleSource{
					Kind: "ClusterRoleBinding",
					Name: "circleci-cluster-admin",
				},
			}},
			"two": {{
				Kind: "ClusterRole",
				Name: "cluster-admin",
				Source: simpleRoleSource{
					Kind: "RoleBinding",
					Name: "testing-sa",
				},
			}},
			"three": {{
				Kind: "ClusterRole",
				Name: "cluster-admin",
				Source: simpleRoleSource{
					Kind: "RoleBinding",
					Name: "testing-sa",
				},
			}},
		},
	}

	assert.EqualValues(t, expectedRbacSubject, l.RbacSubjectsByScope["joe"])
	assert.EqualValues(t, expectedRbacSubject, l.RbacSubjectsByScope["sue"])
	assert.EqualValues(t, expectedRbacSubjectSA, l.RbacSubjectsByScope["circleci:circleci"])
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

	assert.Len(t, l.RbacSubjectsByScope, 0, "Expected no rbac subjects initially")

	l.loadGkeIamPolicy(policy)

	assert.Len(t, l.RbacSubjectsByScope, 4, "Expected 4 rbac subjects")

	assert.EqualValues(t, l.RbacSubjectsByScope["jane@example.com"], rbacSubject{
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

	assert.EqualValues(t, l.RbacSubjectsByScope["joe@example.com"], rbacSubject{
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

	assert.EqualValues(t, l.RbacSubjectsByScope["devs@example.com"], rbacSubject{
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

	assert.EqualValues(t, l.RbacSubjectsByScope["ci@example.iam.gserviceaccount.com"], rbacSubject{
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
	l.Filter = "example"
	l.SubjectKind = "user"

	assert.Len(t, l.RbacSubjectsByScope, 0, "Expected no rbac subjects initially")

	l.loadGkeIamPolicy(policy)

	assert.Len(t, l.RbacSubjectsByScope, 2, "Expected 2 rbac subjects")

	assert.EqualValues(t, rbacSubject{
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
	}, l.RbacSubjectsByScope["jane@example.com"])

	assert.EqualValues(t, rbacSubject{
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
	}, l.RbacSubjectsByScope["joe@example.com"])
}

func genLister() Lister {
	return Lister{
		Clientset:           testclient.NewSimpleClientset(),
		RbacSubjectsByScope: make(map[string]rbacSubject),
	}
}

func loadAll(t *testing.T, l Lister) {
	err := l.loadAll()

	assert.Nil(t, err, "Expected no error loading all rbac Bindings")
}

func loadRoleBindings(t *testing.T, l Lister) {
	err := l.loadRoleBindings()

	assert.Nil(t, err, "Expected no error loading role bindings")
}

func createRoleBindings(t *testing.T, l Lister) {
	roleBindings := []rbacv1.RoleBinding{
		{
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
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testing-sa",
				Namespace: "two",
			},
			Subjects: []rbacv1.Subject{{
				Name:      "circleci",
				Kind:      "ServiceAccount",
				Namespace: "circleci",
			}},
			RoleRef: rbacv1.RoleRef{
				Kind: "ClusterRole",
				Name: "cluster-admin",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testing-sa",
				Namespace: "three",
			},
			Subjects: []rbacv1.Subject{{
				Name:      "circleci",
				Kind:      "ServiceAccount",
				Namespace: "circleci",
			}},
			RoleRef: rbacv1.RoleRef{
				Kind: "ClusterRole",
				Name: "cluster-admin",
			},
		},
	}

	for _, roleBinding := range roleBindings {
		_, err := l.Clientset.RbacV1().RoleBindings(roleBinding.Namespace).Create(context.Background(), &roleBinding, metav1.CreateOptions{})
		assert.Nil(t, err, "Expected no error creating role bindings")
	}
}

func loadClusterRoleBindings(t *testing.T, l Lister) {
	err := l.loadClusterRoleBindings()

	assert.Nil(t, err, "Expected no error loading cluster role bindings")
}

func createClusterRoleBindings(t *testing.T, l Lister) {
	clusterRoleBindings := []rbacv1.ClusterRoleBinding{
		{
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
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "circleci-cluster-admin",
			},
			Subjects: []rbacv1.Subject{{
				Name:      "circleci",
				Kind:      "ServiceAccount",
				Namespace: "circleci",
			}},
			RoleRef: rbacv1.RoleRef{
				Kind: "ClusterRole",
				Name: "cluster-admin",
			},
		},
	}

	for _, clusterRoleBinding := range clusterRoleBindings {
		_, err := l.Clientset.RbacV1().ClusterRoleBindings().Create(context.Background(), &clusterRoleBinding, metav1.CreateOptions{})
		assert.Nil(t, err, "Expected no error creating cluster role bindings")
	}
}
