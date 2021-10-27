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
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// Required for GKE, OIDC, and more
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// List outputs rbac bindings where subject names match given string
func List(args []string, kubeConfig, kubeContext, outputFormat, subjectKind string, enableGke bool) {

	clientConfig := getClientConfig(kubeConfig, kubeContext)

	kubeconfig, err := clientConfig.ClientConfig()
	if err != nil {
		fmt.Printf("Error getting Kubernetes config: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		fmt.Printf("Error generating Kubernetes clientset from kubeconfig: %v\n", err)
		os.Exit(2)
	}

	filter := ""
	if len(args) > 0 {
		filter = args[0]
	}

	l := Lister{
		Filter:              filter,
		SubjectKind:         subjectKind,
		Clientset:           clientset,
		RbacSubjectsByScope: make(map[string]rbacSubject),
	}

	if enableGke {
		rawConfig, err := clientConfig.RawConfig()
		if err != nil {
			fmt.Printf("Error getting Kubernetes raw config: %v\n", err)
			os.Exit(3)
		}

		ci := getClusterInfo(&rawConfig, kubeContext)
		l.GkeParsedProjectName = ci.ParsedProjectName
	}

	loadErr := l.loadAll()
	if loadErr != nil {
		fmt.Printf("Error loading RBAC Bindings: %v\n", loadErr)
		os.Exit(4)
	}

	l.printRbacBindings(outputFormat)
}

func getClientConfig(kubeConfig, kubeContext string) clientcmd.ClientConfig {
	configRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configRules.ExplicitPath = kubeConfig
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		configRules,
		&clientcmd.ConfigOverrides{CurrentContext: kubeContext},
	)
}
