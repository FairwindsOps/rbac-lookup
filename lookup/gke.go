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
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	// Required for GKE Auth
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type gkeClusterInfo struct {
	ClusterName       string
	Region            string
	ParsedProjectName string
}

func getClusterInfo(c *clientcmdapi.Config, kubeContext string) *gkeClusterInfo {
	context := c.Contexts[c.CurrentContext]
	if kubeContext != "" {
		context = c.Contexts[kubeContext]
	}

	ci := gkeClusterInfo{}

	if context != nil && context.Cluster != "" {
		s := strings.Split(context.Cluster, "_")
		if s[0] == "gke" {
			ci.ClusterName = s[3]
			ci.Region = s[2]
			ci.ParsedProjectName = s[1]
		}
	}

	return &ci
}

func loadGkeIAMPolicy(parsedProjectName string) (*cloudresourcemanager.Policy, error) {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, cloudresourcemanager.CloudPlatformReadOnlyScope)
	if err != nil {
		fmt.Println("Error initializing Google API client")
		return nil, err
	}

	crmService, err := cloudresourcemanager.New(c)
	if err != nil {
		fmt.Println("Error initializing Google Cloud Resource Manager")
		return nil, err
	}

	ipr := &cloudresourcemanager.GetIamPolicyRequest{}

	var policy *cloudresourcemanager.Policy
	var err1, err2, err3 error

	policy, err1 = crmService.Projects.GetIamPolicy(parsedProjectName, ipr).Context(ctx).Do()
	if err1 != nil {
		fmt.Printf("Could not load IAM policy for %s project from parsed kubeconfig\n", parsedProjectName)

		var credentials *google.Credentials
		credentials, err2 = google.FindDefaultCredentials(ctx, cloudresourcemanager.CloudPlatformReadOnlyScope)

		if err2 != nil {
			return nil, err2
		}

		if credentials.ProjectID == "" {
			fmt.Println("No project ID found in default GCP credentials")
			return getPolicyFromEnvVar(crmService, ipr)
		}

		policy, err3 = crmService.Projects.GetIamPolicy(credentials.ProjectID, ipr).Context(ctx).Do()

		if err3 != nil {
			fmt.Printf("Could not load IAM policy for %s project from default GCP credentials\n", credentials.ProjectID)
			return getPolicyFromEnvVar(crmService, ipr)
		}

		return policy, nil
	}

	return policy, nil
}

func getPolicyFromEnvVar(crmService *cloudresourcemanager.Service, ipr *cloudresourcemanager.GetIamPolicyRequest) (*cloudresourcemanager.Policy, error) {
	envVar := os.Getenv("CLOUDSDK_CORE_PROJECT")
	if envVar == "" {
		return nil, errors.New("Error loading IAM policies for GKE, try setting CLOUDSDK_CORE_PROJECT environment variable")
	}

	policy, err := crmService.Projects.GetIamPolicy(envVar, ipr).Context(context.Background()).Do()

	if err != nil {
		fmt.Printf("Could not load IAM policy for %s project from CLOUDSDK_CORE_PROJECT environment variable\n", envVar)
		return nil, err
	}

	fmt.Printf("GCP IAM policy loaded for %s project from CLOUDSDK_CORE_PROJECT environment variable\n\n", envVar)
	return policy, nil
}
