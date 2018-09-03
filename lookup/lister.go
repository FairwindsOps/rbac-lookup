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
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
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

type lister struct {
	rbacSubjectsByScope map[string]rbacSubject
	clientset           kubernetes.Interface
	filter              string
}

func (l *lister) loadAll() error {
	rbErr := l.loadRoleBindings()

	if rbErr != nil {
		return rbErr
	}

	crbErr := l.loadClusterRoleBindings()

	if crbErr != nil {
		return crbErr
	}

	return nil
}

func (l *lister) loadRoleBindings() error {
	roleBindings, err := l.clientset.RbacV1().RoleBindings("").List(metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, roleBinding := range roleBindings.Items {
		for _, subject := range roleBinding.Subjects {
			if l.filter == "" || strings.Contains(subject.Name, l.filter) {
				if rbacSubj, exist := l.rbacSubjectsByScope[subject.Name]; exist {
					rbacSubj.addRoleBinding(&roleBinding)
				} else {
					rbacSubj := rbacSubject{
						Kind:         subject.Kind,
						RolesByScope: make(map[string][]simpleRole),
					}
					rbacSubj.addRoleBinding(&roleBinding)
					l.rbacSubjectsByScope[subject.Name] = rbacSubj
				}
			}
		}
	}

	return nil
}

func (l *lister) loadClusterRoleBindings() error {
	clusterRoleBindings, err := l.clientset.RbacV1().ClusterRoleBindings().List(metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		for _, subject := range clusterRoleBinding.Subjects {
			if l.filter == "" || strings.Contains(subject.Name, l.filter) {
				if rbacSubj, exist := l.rbacSubjectsByScope[subject.Name]; exist {
					rbacSubj.addClusterRoleBinding(&clusterRoleBinding)
				} else {
					rbacSubj := rbacSubject{
						Kind:         subject.Kind,
						RolesByScope: make(map[string][]simpleRole),
					}
					rbacSubj.addClusterRoleBinding(&clusterRoleBinding)
					l.rbacSubjectsByScope[subject.Name] = rbacSubj
				}
			}
		}
	}

	return nil
}

func (l *lister) printRbacBindings(outputFormat string) {
	if len(l.rbacSubjectsByScope) < 1 {
		fmt.Println("No RBAC Bindings found")
		return
	}

	names := make([]string, 0, len(l.rbacSubjectsByScope))
	for name := range l.rbacSubjectsByScope {
		names = append(names, name)
	}
	sort.Strings(names)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)

	if outputFormat == "wide" {
		fmt.Fprintln(w, "SUBJECT\t SCOPE\t ROLE\t SOURCE")
	} else {
		fmt.Fprintln(w, "SUBJECT\t SCOPE\t ROLE")
	}

	for _, subjectName := range names {
		rbacSubject := l.rbacSubjectsByScope[subjectName]
		for scope, simpleRoles := range rbacSubject.RolesByScope {
			for _, simpleRole := range simpleRoles {
				if outputFormat == "wide" {
					fmt.Fprintf(w, "%s/%s \t %s\t %s/%s\t %s/%s\n", rbacSubject.Kind, subjectName, scope, simpleRole.Kind, simpleRole.Name, simpleRole.Source.Kind, simpleRole.Source.Name)
				} else {
					fmt.Fprintf(w, "%s \t %s\t %s/%s\n", subjectName, scope, simpleRole.Kind, simpleRole.Name)
				}
			}
		}
	}
	w.Flush()
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

func (rbacSubj *rbacSubject) addRoleBinding(roleBinding *rbacv1.RoleBinding) {
	simpleRole := simpleRole{
		Name:   roleBinding.RoleRef.Name,
		Source: simpleRoleSource{Name: roleBinding.Name, Kind: "RoleBinding"},
	}

	simpleRole.Kind = roleBinding.RoleRef.Kind
	rbacSubj.RolesByScope[roleBinding.Namespace] = append(rbacSubj.RolesByScope[roleBinding.Namespace], simpleRole)
}
