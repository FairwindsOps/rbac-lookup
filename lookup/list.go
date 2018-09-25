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
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// Required for GKE Auth
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type clusterInfo struct {
	ClusterName    string
	GkeZone        string
	GkeProjectName string
}

// List outputs rbac bindings where subject names match given string
func List(args []string, outputFormat string) {
	kubeconfig := getKubeConfig()
	clientset, err := getClientSet(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	ci := getClusterInfo(kubeconfig)

	filter := ""
	if len(args) > 0 {
		filter = args[0]
	}

	l := lister{
		filter:              filter,
		clientset:           clientset,
		gkeProjectName:      ci.GkeProjectName,
		rbacSubjectsByScope: make(map[string]rbacSubject),
	}

	loadErr := l.loadAll()
	if loadErr != nil {
		fmt.Printf("Error loading RBAC: %v\n", loadErr)
		os.Exit(1)
	}

	l.printRbacBindings(outputFormat)
}

func getKubeConfig() string {
	var kubeconfig string
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	} else if home := homeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		fmt.Println("Parsing kubeconfig failed, please set KUBECONFIG env var")
		os.Exit(1)
	}

	if _, err := os.Stat(kubeconfig); err != nil {
		// kubeconfig doesn't exist
		fmt.Printf("%s does not exist - please make sure you have a kubeconfig configured.\n", kubeconfig)
		panic(err.Error())
	}

	return kubeconfig
}

func getClusterInfo(kubeconfig string) *clusterInfo {
	c, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	s := strings.Split(c.CurrentContext, "_")
	if s[0] == "gke" {
		return &clusterInfo{
			ClusterName:    s[3],
			GkeZone:        s[2],
			GkeProjectName: s[1],
		}
	}
	return &clusterInfo{}
}

func getClientSet(kubeconfig string) (*kubernetes.Clientset, error) {
	flag.Parse()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
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
