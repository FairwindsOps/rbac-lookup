package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type RbacSubject struct {
	Kind             string
	ClusterWideRoles []SimpleRole
	RolesByNamespace map[string][]SimpleRole
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

func main() {
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
	clientset, err := kubernetes.NewForConfig(config)
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
					Kind:             subject.Kind,
					RolesByNamespace: make(map[string][]SimpleRole),
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
					Kind:             subject.Kind,
					RolesByNamespace: make(map[string][]SimpleRole),
				}
				addSimpleRoleCRB(&rbacSubject, &clusterRoleBinding)
				rbacBindings[subject.Name] = rbacSubject
			}
		}
	}

	printRbacBindings(rbacBindings)
}

func printRbacBindings(rbacBindings map[string]RbacSubject) {
	names := make([]string, 0, len(rbacBindings))
	for name := range rbacBindings {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, subjectName := range names {
		rbacSubject := rbacBindings[subjectName]
		color.Green("%s (%s)", subjectName, rbacSubject.Kind)
		if len(rbacSubject.ClusterWideRoles) > 0 {
			color.Cyan("- cluster wide")
			for _, simpleRole := range rbacSubject.ClusterWideRoles {
				fmt.Printf("  - %s: %s\n", simpleRole.Kind, simpleRole.Name)
				fmt.Printf("    source: %s\n", simpleRole.Source)
			}
		}

		for namespace, simpleRoles := range rbacSubject.RolesByNamespace {
			color.Cyan("- %s", namespace)
			for _, simpleRole := range simpleRoles {
				fmt.Printf("  - %s: %s\n", simpleRole.Kind, simpleRole.Name)
				fmt.Printf("    source: %s\n", simpleRole.Source)
			}
		}
		fmt.Println("")
	}
}

func addSimpleRoleCRB(rbacSubject *RbacSubject, clusterRoleBinding *rbacv1.ClusterRoleBinding) {
	simpleRole := SimpleRole{
		Name:   clusterRoleBinding.RoleRef.Name,
		Source: SimpleRoleSource{Name: clusterRoleBinding.Name, Kind: "ClusterRoleBinding"},
	}

	if clusterRoleBinding.RoleRef.Kind == "ClusterRole" {
		simpleRole.Kind = "cluster role"
	} else {
		simpleRole.Kind = "role"
	}

	rbacSubject.ClusterWideRoles = append(rbacSubject.ClusterWideRoles, simpleRole)
}

func addSimpleRole(rbacSubject *RbacSubject, roleBinding *rbacv1.RoleBinding) {
	simpleRole := SimpleRole{
		Name:   roleBinding.RoleRef.Name,
		Source: SimpleRoleSource{Name: roleBinding.Name, Kind: "ClusterRoleBinding"},
	}

	if roleBinding.RoleRef.Kind == "ClusterRole" {
		simpleRole.Kind = "cluster role"
	} else {
		simpleRole.Kind = "role"
	}

	rbacSubject.RolesByNamespace[roleBinding.Namespace] = append(rbacSubject.RolesByNamespace[roleBinding.Namespace], simpleRole)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
