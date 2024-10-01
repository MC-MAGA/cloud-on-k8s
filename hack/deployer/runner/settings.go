Deployer
Deployer is the provisioning tool that aims to be the interface to multiple Kubernetes providers. Currently, it supports GKE and AKS.

Typical usage
Provision
GKE

Install Google Cloud SDK

Install Google Cloud SDK beta components by running gcloud components install beta

Make sure that container registry authentication is correctly configured as described here

Set GCLOUD_PROJECT to the name of the GCloud project you wish to use

(optional) Set CLOUDSDK_CONFIG to a directory which should be used for gcloud SDK if you don't want to have the default one overwritten.

Run from the project root:

make switch-gke bootstrap-cloud
AKS

Install Azure CLI

Set RESOURCE_GROUP to the name of the Resource Group you wish to deploy in

Run from the project root:

make switch-aks bootstrap-cloud
Kind

No need to install the Kind CLI. Deployer will do that for you and run Kind inside a Docker container without changing the host system.

Run from the project root:

make switch-kind bootstrap-cloud
This will give you a working Kind cluster based on default values. See Advanced usage on how to tweak these configuration defaults if the need arises. Relevant parameters for Kind are: client_version which is the version of Kind to use. Make sure to check the Kind release notes when changing the client version and make sure kubernetesVersion and client_version are compatible. kind.nodeImage allows you to use a specific Kind node image matching your chosen Kind version. Again, the Kind release notes list the compatible pre-built node images for each version. kind.ipFamily allows you to switch between either an IPv4 or IPv6 network setup.

Deprovision
make delete-cloud
Advanced usage
Deployer uses two config files:

config/plans.yml - to store defaults/baseline settings for different use cases (different providers, CI/dev)
config/deployer-config-*.yml - to "pick" on of the predefined configs from config/plans.yml and allow overriding settings.
You can adjust many parameters that clusters are deployed with. Exhaustive list is defined in settings.go.

Running make switch-* (eg. make-switch-gke) changes the current context. Running make create-default-config generates config/deployer-config-*.yml file for the respective provider using environment variables specific to that providers configuration needs. After the file is generated, you can make edit it to suit your needs and run make bootstrap-cloud to deploy. Currently chosen provider is stored in config/provider file.

You can run deployer directly (not via Makefile in repo root). For details run:

./deployer help// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package runner

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Plans encapsulates list of plans, expected to map to a file
type Plans struct {
	Plans []Plan
}

// Plan encapsulates information needed to provision a cluster
type Plan struct {
	Id                string `yaml:"id"` //nolint:revive
	Operation         string `yaml:"operation"`
	ClusterName       string `yaml:"clusterName"`
	ClientVersion     string `yaml:"clientVersion"`
	ClientBuildDefDir string `yaml:"clientBuildDefDir"`
	Provider          string `yaml:"provider"`
	KubernetesVersion string `yaml:"kubernetesVersion"`
	MachineType       string `yaml:"machineType"`
	// Abbreviations not all-caps to allow merging with mergo in  `merge` as mergo does not understand struct tags and
	// we use lowercase in the YAML
	Gke                     *GKESettings   `yaml:"gke,omitempty"`
	Aks                     *AKSSettings   `yaml:"aks,omitempty"`
	Ocp                     *OCPSettings   `yaml:"ocp,omitempty"`
	Eks                     *EKSSettings   `yaml:"eks,omitempty"`
	Kind                    *KindSettings  `yaml:"kind,omitempty"`
	Tanzu                   *TanzuSettings `yaml:"tanzu,omitempty"`
	ServiceAccount          bool           `yaml:"serviceAccount"`
	EnforceSecurityPolicies bool           `yaml:"enforceSecurityPolicies"`
	DiskSetup               string         `yaml:"diskSetup"`
}

// GKESettings encapsulates settings specific to GKE
type GKESettings struct {
	GCloudProject    string `yaml:"gCloudProject"`
	Region           string `yaml:"region"`
	LocalSsdCount    int    `yaml:"localSsdCount"`
	NodeCountPerZone int    `yaml:"nodeCountPerZone"`
	GcpScopes        string `yaml:"gcpScopes"`
	ClusterIPv4CIDR  string `yaml:"clusterIpv4Cidr"`
	ServicesIPv4CIDR string `yaml:"servicesIpv4Cidr"`
	Private          bool   `yaml:"private"`
	NetworkPolicy    bool   `yaml:"networkPolicy"`
	Autopilot        bool   `yaml:"autopilot"`
}

// AKSSettings encapsulates settings specific to AKS
type AKSSettings struct {
	ResourceGroup string `yaml:"resourceGroup"`
	Location      string `yaml:"location"`
	Zones         string `yaml:"zones"`
	NodeCount     int    `yaml:"nodeCount"`
}

// OCPSettings encapsulates settings specific to OCP on GCloud
type OCPSettings struct {
	BaseDomain    string `yaml:"baseDomain"`
	GCloudProject string `yaml:"gCloudProject"`
	Region        string `yaml:"region"`
	AdminUsername string `yaml:"adminUsername"`
	WorkDir       string `yaml:"workDir"`
	StickyWorkDir bool   `yaml:"stickyWorkDir"`
	PullSecret    string `yaml:"pullSecret"`
	LocalSsdCount int    `yaml:"localSsdCount"`
	NodeCount     int    `yaml:"nodeCount"`
}

// EKSSettings are specific to Amazon EKS.
type EKSSettings struct {
	NodeAMI   string `yaml:"nodeAMI"`
	NodeCount int    `yaml:"nodeCount"`
	Region    string `yaml:"region"`
	WorkDir   string `yaml:"workDir"`
}

type KindSettings struct {
	NodeCount int    `yaml:"nodeCount"`
	NodeImage string `yaml:"nodeImage"`
	IPFamily  string `yaml:"ipFamily"`
}

type TanzuSettings struct {
	AKSSettings    `yaml:",inline"`
	InstallerImage string `yaml:"installerImage"`
	WorkDir        string `yaml:"workDir"`
	SSHPubKey      string `yaml:"sshPubKey"`
}

// RunConfig encapsulates Id used to choose a plan and a map of overrides to apply to the plan, expected to map to a file
type RunConfig struct {
	Id        string                 `yaml:"id"` //nolint:revive
	Overrides map[string]interface{} `yaml:"overrides"`
}

func ParseFiles(plansFile, runConfigFile string) (Plans, RunConfig, error) {
	yml, err := os.ReadFile(plansFile)
	if err != nil {
		return Plans{}, RunConfig{}, err
	}

	var plans Plans
	err = yaml.Unmarshal(yml, &plans)
	if err != nil {
		return Plans{}, RunConfig{}, err
	}

	yml, err = os.ReadFile(runConfigFile)
	if err != nil {
		return Plans{}, RunConfig{}, err
	}

	var runConfig RunConfig
	err = yaml.Unmarshal(yml, &runConfig)
	if err != nil {
		return Plans{}, RunConfig{}, err
	}

	return plans, runConfig, nil
}
