package lookup

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type RbacSubject struct {
	Kind         string
	RolesByScope map[string][]SimpleRole
}

type SimpleRole struct {
	Kind   string
	Name   string
	Source SimpleRoleSource
}

type SimpleRoleSource struct {
	Kind string
	Name string
}

func ListAll() {
	clientset, err := getClientSet()
	if err != nil {
		panic(err.Error())
	}

	roleBindings, err := clientset.RbacV1().RoleBindings("").List(metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	rbacBindings := make(map[string]RbacSubject)

	for _, roleBinding := range roleBindings.Items {
		for _, subject := range roleBinding.Subjects {
			if rbacSubject, exist := rbacBindings[subject.Name]; exist {
				addSimpleRole(&rbacSubject, &roleBinding)
			} else {
				rbacSubject := RbacSubject{
					Kind:         subject.Kind,
					RolesByScope: make(map[string][]SimpleRole),
				}
				addSimpleRole(&rbacSubject, &roleBinding)
				rbacBindings[subject.Name] = rbacSubject
			}
		}
	}

	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		for _, subject := range clusterRoleBinding.Subjects {
			if rbacSubject, exist := rbacBindings[subject.Name]; exist {
				addSimpleRoleCRB(&rbacSubject, &clusterRoleBinding)
			} else {
				rbacSubject := RbacSubject{
					Kind:         subject.Kind,
					RolesByScope: make(map[string][]SimpleRole),
				}
				addSimpleRoleCRB(&rbacSubject, &clusterRoleBinding)
				rbacBindings[subject.Name] = rbacSubject
			}
		}
	}

	printRbacBindings(rbacBindings)
}

func printRbacBindings(rbacBindings map[string]RbacSubject) {
	if len(rbacBindings) < 1 {
		fmt.Println("No RBAC Bindings found")
		return
	}

	names := make([]string, 0, len(rbacBindings))
	for name := range rbacBindings {
		names = append(names, name)
	}
	sort.Strings(names)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)

	// fmt.Fprintln(w, "SUBJECT\t SCOPE\t ROLE\t SOURCE")
	fmt.Fprintln(w, "SUBJECT\t SCOPE\t ROLE")
	for _, subjectName := range names {
		rbacSubject := rbacBindings[subjectName]
		for scope, simpleRoles := range rbacSubject.RolesByScope {
			for _, simpleRole := range simpleRoles {
				fmt.Fprintf(w, "%s \t %s\t %s/%s\n", subjectName, scope, simpleRole.Kind, simpleRole.Name)
			}
		}
	}
	w.Flush()
}

func addSimpleRoleCRB(rbacSubject *RbacSubject, clusterRoleBinding *rbacv1.ClusterRoleBinding) {
	simpleRole := SimpleRole{
		Name:   clusterRoleBinding.RoleRef.Name,
		Source: SimpleRoleSource{Name: clusterRoleBinding.Name, Kind: "ClusterRoleBinding"},
	}

	simpleRole.Kind = clusterRoleBinding.RoleRef.Kind
	scope := "cluster-wide"
	rbacSubject.RolesByScope[scope] = append(rbacSubject.RolesByScope[scope], simpleRole)
}

func addSimpleRole(rbacSubject *RbacSubject, roleBinding *rbacv1.RoleBinding) {
	simpleRole := SimpleRole{
		Name:   roleBinding.RoleRef.Name,
		Source: SimpleRoleSource{Name: roleBinding.Name, Kind: "RoleBinding"},
	}

	simpleRole.Kind = roleBinding.RoleRef.Kind
	rbacSubject.RolesByScope[roleBinding.Namespace] = append(rbacSubject.RolesByScope[roleBinding.Namespace], simpleRole)
}

func getClientSet() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset

	return kubernetes.NewForConfig(config)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
