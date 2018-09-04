// Copyright 2018 ReactiveOps
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
		_, err := l.clientset.RbacV1().RoleBindings(roleBinding.Namespace).Create(&roleBinding)
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
		_, err := l.clientset.RbacV1().ClusterRoleBindings().Create(&clusterRoleBinding)
		assert.Nil(t, err, "Expected no error creating cluster role bindings")
	}
}
