/*
Copyright 2018 Heptio Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/sonobuoy/pkg/buildinfo"
	"github.com/vmware-tanzu/sonobuoy/pkg/errlog"
)

type versionFlags struct {
	kubecfg Kubeconfig
	short   bool
}

var versionflags versionFlags

func NewCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print sonobuoy version",
		Run:   runVersion,
		Args:  cobra.ExactArgs(0),
	}

	AddKubeconfigFlag(&versionflags.kubecfg, cmd.Flags())
	AddShortFlag(&versionflags.short, cmd.Flags())

	return cmd
}

func runVersion(cmd *cobra.Command, args []string) {
	if versionflags.short {
		fmt.Println(buildinfo.Version)
		return
	}

	fmt.Printf("Sonobuoy Version: %s\n", buildinfo.Version)
	fmt.Printf("MinimumKubeVersion: %s\n", buildinfo.MinimumKubeVersion)
	fmt.Printf("MaximumKubeVersion: %s\n", buildinfo.MaximumKubeVersion)
	fmt.Printf("GitSHA: %s\n", buildinfo.GitSHA)

	// Get Kubernetes version, this is last so that the regular version information
	// will be shown even if the API server cannot be contacted and throws an error
	apiVersion, skipk8sCheck := getK8Sversion()
	if !skipk8sCheck {
		fmt.Println("API Version: ", apiVersion)
	} else {
		fmt.Println("API Version check skipped due to missing `--kubeconfig` or other error")
	}
}

func getK8Sversion() (string, bool) {

	if versionflags.kubecfg.String() != "" {
		sbc, err := getSonobuoyClientFromKubecfg(versionflags.kubecfg)
		if err != nil {
			errlog.LogError(errors.Wrap(err, "could not create sonobuoy client"))
			return "", true
		}

		client, err := sbc.Client()
		if err != nil {
			errlog.LogError(err)
			return "", true
		}

		apiVersion, err := client.Discovery().ServerVersion()
		if err != nil {
			errlog.LogError(err)
			return "", true
		}

		return apiVersion.GitVersion, false
	}

	return "", true
}
