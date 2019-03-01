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

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"

	// Required for GKE Auth
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type lister struct {
	clientset           kubernetes.Interface
	filter              string
	gkeProjectName      string
	subjectKind         string
	rbacSubjectsByScope map[string]rbacSubject
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

	if l.gkeProjectName != "" {
		gkeErr := l.loadGkeRoleBindings()

		if gkeErr != nil {
			return gkeErr
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

func (l *lister) loadRoleBindings() error {
	roleBindings, err := l.clientset.RbacV1().RoleBindings("").List(metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, roleBinding := range roleBindings.Items {
		for _, subject := range roleBinding.Subjects {
			if l.nameMatches(subject.Name) && l.kindMatches(subject.Kind) {
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
			if l.nameMatches(subject.Name) && l.kindMatches(subject.Kind) {
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

func (l *lister) loadGkeIamPolicy(policy *cloudresourcemanager.Policy) {
	for _, binding := range policy.Bindings {
		if sr, ok := gkeIamRoles[binding.Role]; ok {
			for _, member := range binding.Members {
				s := strings.Split(member, ":")
				memberKind := strings.Title(s[0])
				memberName := s[1]
				if l.nameMatches(memberName) && l.kindMatches(memberKind) {
					rbacSubj, exist := l.rbacSubjectsByScope[memberName]
					if !exist {
						rbacSubj = rbacSubject{
							Kind:         memberKind,
							RolesByScope: make(map[string][]simpleRole),
						}
					}

					rbacSubj.RolesByScope[gkeIamScope] = append(rbacSubj.RolesByScope[gkeIamScope], sr)
					l.rbacSubjectsByScope[memberName] = rbacSubj
				}
			}
		}
	}
}

func (l *lister) loadGkeRoleBindings() error {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, cloudresourcemanager.CloudPlatformScope)
	if err != nil {
		return err
	}

	crmService, err := cloudresourcemanager.New(c)
	if err != nil {
		return err
	}

	resource := l.gkeProjectName
	ipr := &cloudresourcemanager.GetIamPolicyRequest{}

	policy, err := crmService.Projects.GetIamPolicy(resource, ipr).Context(ctx).Do()
	if err != nil {
		return err
	}

	l.loadGkeIamPolicy(policy)

	return nil
}

func (l *lister) nameMatches(name string) bool {
	return l.filter == "" || strings.Contains(name, l.filter)
}

func (l *lister) kindMatches(kind string) bool {
	if l.subjectKind == "" {
		return true
	}

	lowerKind := strings.ToLower(kind)

	return lowerKind == l.subjectKind
}
