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
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"google.golang.org/api/cloudresourcemanager/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"

	// Required for different auth providers like GKE, OIDC
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type Lister struct {
	Clientset            kubernetes.Interface
	Filter               string
	GkeParsedProjectName string
	SubjectKind          string
	RbacSubjectsByScope  map[string]rbacSubject
}

func (l *Lister) loadAll() error {
	rbErr := l.loadRoleBindings()

	if rbErr != nil {
		return rbErr
	}

	crbErr := l.loadClusterRoleBindings()

	if crbErr != nil {
		return crbErr
	}

	if l.GkeParsedProjectName != "" {
		policy, gkeErr := loadGkeIAMPolicy(l.GkeParsedProjectName)

		if gkeErr != nil {
			return gkeErr
		}

		l.loadGkeIamPolicy(policy)
	}

	return nil
}

func (l *Lister) printRbacBindings(outputFormat string) {
	if len(l.RbacSubjectsByScope) < 1 {
		fmt.Println("No RBAC Bindings found")
		return
	}

	names := make([]string, 0, len(l.RbacSubjectsByScope))
	for name := range l.RbacSubjectsByScope {
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
		rbacSubject := l.RbacSubjectsByScope[subjectName]
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

func (l *Lister) loadRoleBindings() error {
	roleBindings, err := l.Clientset.RbacV1().RoleBindings("").List(context.Background(), metav1.ListOptions{})

	if err != nil {
		fmt.Println("Error loading role bindings")
		return err
	}

	for _, roleBinding := range roleBindings.Items {
		for _, subject := range roleBinding.Subjects {
			if l.nameMatches(subject.Name) && l.kindMatches(subject.Kind) {
				subjectKey := subject.Name
				if subject.Kind == "ServiceAccount" {
					subjectKey = fmt.Sprintf("%s:%s", subject.Namespace, subject.Name)
				}
				if rbacSubj, exist := l.RbacSubjectsByScope[subjectKey]; exist {
					rbacSubj.addRoleBinding(&roleBinding)
				} else {
					rbacSubj := rbacSubject{
						Kind:         subject.Kind,
						RolesByScope: make(map[string][]simpleRole),
					}
					rbacSubj.addRoleBinding(&roleBinding)

					l.RbacSubjectsByScope[subjectKey] = rbacSubj
				}
			}
		}
	}

	return nil
}

func (l *Lister) loadClusterRoleBindings() error {
	clusterRoleBindings, err := l.Clientset.RbacV1().ClusterRoleBindings().List(context.Background(), metav1.ListOptions{})

	if err != nil {
		fmt.Println("Error loading cluster role bindings")
		return err
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		for _, subject := range clusterRoleBinding.Subjects {
			if l.nameMatches(subject.Name) && l.kindMatches(subject.Kind) {
				subjectKey := subject.Name
				if subject.Kind == "ServiceAccount" {
					subjectKey = fmt.Sprintf("%s:%s", subject.Namespace, subject.Name)
				}
				if rbacSubj, exist := l.RbacSubjectsByScope[subjectKey]; exist {
					rbacSubj.addClusterRoleBinding(&clusterRoleBinding)
				} else {
					rbacSubj := rbacSubject{
						Kind:         subject.Kind,
						RolesByScope: make(map[string][]simpleRole),
					}
					rbacSubj.addClusterRoleBinding(&clusterRoleBinding)

					l.RbacSubjectsByScope[subjectKey] = rbacSubj
				}
			}
		}
	}

	return nil
}

func (l *Lister) loadGkeIamPolicy(policy *cloudresourcemanager.Policy) {
	for _, binding := range policy.Bindings {
		if sr, ok := gkeIamRoles[binding.Role]; ok {
			for _, member := range binding.Members {
				s := strings.Split(member, ":")
				memberKind := strings.Title(s[0])
				memberName := s[1]
				if l.nameMatches(memberName) && l.kindMatches(memberKind) {
					rbacSubj, exist := l.RbacSubjectsByScope[memberName]
					if !exist {
						rbacSubj = rbacSubject{
							Kind:         memberKind,
							RolesByScope: make(map[string][]simpleRole),
						}
					}

					rbacSubj.RolesByScope[gkeIamScope] = append(rbacSubj.RolesByScope[gkeIamScope], sr)
					l.RbacSubjectsByScope[memberName] = rbacSubj
				}
			}
		}
	}
}

func (l *Lister) nameMatches(name string) bool {
	return l.Filter == "" || strings.Contains(name, l.Filter)
}

func (l *Lister) kindMatches(kind string) bool {
	if l.SubjectKind == "" {
		return true
	}

	lowerKind := strings.ToLower(kind)

	return lowerKind == l.SubjectKind
}
