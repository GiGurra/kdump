package config

import (
	"github.com/gigurra/kdump/internal/k8s"
	"github.com/thoas/go-funk"
	"strings"
)

type Encryption interface {
	Encrypt(data string) []byte
}

type NoEncryption struct {
}

func (n NoEncryption) Encrypt(data string) []byte {
	return []byte(data)
}

type AppConfig struct {
	OutputDir             string
	DeletePrevDir         bool
	ExcludedResourceTypes []string
	IncludeSecrets        bool
	SecretsEncryption     Encryption
}

func getDefaultExcludedResourceTypes() []string {
	return []string{
		"limitranges",
		"podtemplates",
		"replicationcontrollers",
		"resourcequotas",
		"events",
		"jobs",
		"jobs.batch",
		"pods",
		"componentstatuses",
		"endpoints",
		"endpointslices.discovery.k8s.io",
		"replicasets.apps",
		"clusterauthtokens",
		"clusteruserattributes",
		"controllerrevisions.apps",
		"apiservices.apiregistration.k8s.io",
		"clusterinformations",
		"felixconfigurations",
		"ippools",
		"nodes",
		"csinodes.storage.k8s.io",
		"csidrivers.storage.k8s.io",
		"priorityclasses.scheduling.k8s.io",
		"ciliumendpoints.cilium.io",
		"ciliumlocalredirectpolicies.cilium.io",
		"ciliumnetworkpolicies.cilium.io",
		"ciliumclusterwidenetworkpolicies.cilium.io",
		"ciliumegressnatpolicies.cilium.io",
		"ciliumexternalworkloads.cilium.io",
		"ciliumidentities.cilium.io",
		"flowschemas.flowcontrol.apiserver.k8s.io",
		"prioritylevelconfigurations.flowcontrol.apiserver.k8s.io",
		"horizontalpodautoscalers.autoscaling",
		"runtimeclasses.node.k8s.io",
		"nodes.metrics.k8s.io",
		"ciliumnodes.cilium.io",
		"events.events.k8s.io",
		"leases.coordination.k8s.io",
		"certificaterequests.cert-manager.io",
		"orders.acme.cert-manager.io",
		"challenges.acme.cert-manager.io",
		"mutatingwebhookconfigurations.admissionregistration.k8s.io",
		"validatingwebhookconfigurations.admissionregistration.k8s.io",
		"certificatesigningrequests.certificates.k8s.io",
		"ingresses.extensions",
		"pods.metrics.k8s.io",
	}
}

func GetDefaultAppConfig() AppConfig {
	return AppConfig{
		OutputDir:             "test",
		ExcludedResourceTypes: getDefaultExcludedResourceTypes(),
		IncludeSecrets:        false,
		SecretsEncryption:     &NoEncryption{},
	}
}

func (config *AppConfig) IsResourceIncluded(resourceType *k8s.ApiResourceType) bool {
	return !funk.ContainsString(config.ExcludedResourceTypes, resourceType.Name) &&
		!funk.ContainsString(config.ExcludedResourceTypes, resourceType.QualifiedName) &&
		(strings.ToLower(resourceType.Name) != "secrets" || config.IncludeSecrets)
}

func (config *AppConfig) FilterIncludedResources(resourceTypes []*k8s.ApiResourceType) []*k8s.ApiResourceType {
	return funk.Filter(resourceTypes, func(r *k8s.ApiResourceType) bool {
		return config.IsResourceIncluded(r)
	}).([]*k8s.ApiResourceType)
}
