/*
Copyright 2021 The KubeVela Authors.

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

package clients

import (
	"fmt"

	pkgmulticluster "github.com/kubevela/pkg/multicluster"
	"github.com/kubevela/workflow/api/v1alpha1"
	"github.com/kubevela/workflow/pkg/cue/packages"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/oam-dev/kubevela/pkg/auth"
	"github.com/oam-dev/kubevela/pkg/utils/common"

	apiConfig "github.com/kubevela/velaux/pkg/server/config"
)

var kubeClient client.Client
var kubeConfig *rest.Config

// SetKubeClient for test
func SetKubeClient(c client.Client) {
	kubeClient = c
}

func setKubeConfig(conf *rest.Config) (err error) {
	if conf == nil {
		conf, err = config.GetConfig()
		if err != nil {
			return err
		}
	}
	kubeConfig = conf
	kubeConfig.Wrap(auth.NewImpersonatingRoundTripper)
	return nil
}

// SetKubeConfig generate the kube config from the config of apiserver
func SetKubeConfig(c apiConfig.Config) error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}
	kubeConfig = conf
	kubeConfig.Burst = c.KubeBurst
	kubeConfig.QPS = float32(c.KubeQPS)
	return setKubeConfig(kubeConfig)
}

// GetKubeClient create and return kube runtime client
func GetKubeClient() (client.Client, error) {
	if kubeClient != nil {
		return kubeClient, nil
	}
	if kubeConfig == nil {
		return nil, fmt.Errorf("please call SetKubeConfig first")
	}
	err := v1alpha1.AddToScheme(common.Scheme)
	if err != nil {
		return nil, err
	}
	return pkgmulticluster.NewClient(kubeConfig, pkgmulticluster.ClientOptions{
		Options: client.Options{Scheme: common.Scheme},
	})
}

// GetKubeConfig create/get kube runtime config
func GetKubeConfig() (*rest.Config, error) {
	if kubeConfig == nil {
		return nil, fmt.Errorf("please call SetKubeConfig first")
	}
	return kubeConfig, nil
}

// GetPackageDiscover get package discover
func GetPackageDiscover() (*packages.PackageDiscover, error) {
	conf, err := GetKubeConfig()
	if err != nil {
		return nil, err
	}
	pd, err := packages.NewPackageDiscover(conf)
	if err != nil {
		if !packages.IsCUEParseErr(err) {
			return nil, err
		}
	}
	return pd, nil
}

// GetDiscoveryClient return a discovery client
func GetDiscoveryClient() (*discovery.DiscoveryClient, error) {
	conf, err := GetKubeConfig()
	if err != nil {
		return nil, err
	}
	dc, err := discovery.NewDiscoveryClientForConfig(conf)
	if err != nil {
		return nil, err
	}
	return dc, nil
}
