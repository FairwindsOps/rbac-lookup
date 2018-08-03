package lookup

import (
	"fmt"
	"os"
	"sort"
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
	clientset           kubernetes.Clientset
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

	return nil
}

func (l *lister) loadClusterRoleBindings() error {
	clusterRoleBindings, err := l.clientset.RbacV1().ClusterRoleBindings().List(metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		for _, subject := range clusterRoleBinding.Subjects {
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

	return nil
}

func (l *lister) printRbacBindings() {
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

	// fmt.Fprintln(w, "SUBJECT\t SCOPE\t ROLE\t SOURCE")
	fmt.Fprintln(w, "SUBJECT\t SCOPE\t ROLE")
	for _, subjectName := range names {
		rbacSubject := l.rbacSubjectsByScope[subjectName]
		for scope, simpleRoles := range rbacSubject.RolesByScope {
			for _, simpleRole := range simpleRoles {
				fmt.Fprintf(w, "%s \t %s\t %s/%s\n", subjectName, scope, simpleRole.Kind, simpleRole.Name)
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
