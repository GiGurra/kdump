package config

import (
	"github.com/gigurra/go-util/crypt"
	"github.com/gigurra/kdump/internal/k8s"
	"github.com/samber/lo"
	"log"
	"strings"
)

type AppConfig struct {
	OutputDir             string
	DeletePrevDir         bool
	ExcludedResourceTypes []string
	SecretsEncryptKey     string
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
		"extractionresults.vulnerabilities.protect.gke.io",
		"endpointslices",
		"ippools",
		"leases",
	}
}

func GetDefaultAppConfig() AppConfig {
	return AppConfig{
		OutputDir:             "test",
		ExcludedResourceTypes: getDefaultExcludedResourceTypes(),
		SecretsEncryptKey:     "",
	}
}

func (config *AppConfig) Validate() {

	validateNonEmpty := func(str string) string {
		if len(strings.TrimSpace(str)) == 0 {
			log.Fatal("empty string is not allowed")
		}
		return str
	}

	validateCryptoKey := func(key string) string {
		crypt.Encrypt("Hello Encrypt", key)
		return key
	}

	validateNonEmpty(config.OutputDir)
	if len(config.SecretsEncryptKey) > 0 {
		validateCryptoKey(config.SecretsEncryptKey)
	}
}

func (config *AppConfig) IncludeSecrets() bool {
	return len(config.SecretsEncryptKey) > 0
}

func (config *AppConfig) IsResourceIncluded(resourceType *k8s.ApiResourceType) bool {
	return !lo.Contains(config.ExcludedResourceTypes, resourceType.Name) &&
		!lo.Contains(config.ExcludedResourceTypes, resourceType.QualifiedName) &&
		(!resourceType.IsSecret() || config.IncludeSecrets())
}

func (config *AppConfig) FilterIncludedResources(resourceTypes []*k8s.ApiResourceType) []*k8s.ApiResourceType {
	return lo.Filter(resourceTypes, func(r *k8s.ApiResourceType, index int) bool {
		return config.IsResourceIncluded(r)
	})
}
