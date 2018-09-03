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
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// Required for GKE Auth
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// List outputs rbac bindings where subject names match given string
func List(args []string, outputFormat string) {
	clientset, err := getClientSet()
	if err != nil {
		panic(err.Error())
	}

	filter := ""
	if len(args) > 0 {
		filter = args[0]
	}

	l := lister{
		filter:              filter,
		clientset:           clientset,
		rbacSubjectsByScope: make(map[string]rbacSubject),
	}

	loadErr := l.loadAll()
	if loadErr != nil {
		fmt.Printf("Error loading RBAC: %v\n", loadErr)
		os.Exit(1)
	}

	l.printRbacBindings(outputFormat)
}

func getClientSet() (*kubernetes.Clientset, error) {
	var kubeconfig string
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	} else if home := homeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		fmt.Println("Parsing kubeconfig failed, please set KUBECONFIG env var")
		os.Exit(1)
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
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
