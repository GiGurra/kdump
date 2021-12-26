package config

import "github.com/thoas/go-funk"

type AppConfig struct {
	OutputDir                string
	AppendContextToOutputDir bool
	ExcludedResourceTypes    []string
	IncludeSecrets           bool
	EncryptSecrets           bool
	SecretsEncryptAlgo       string
	SecretsEncryptKey        string
}

func getDefaultExcludedResourceTypes() []string {
	return []string{
		"events",
		"jobs",
		"pods",
		"componentstatuses",
		"endpoints",
		"endpointslices",
		"replicasets",
		"clusterauthtokens",
		"clusteruserattributes",
		"controllerrevisions",
		"apiservices",
		"clusterinformations",
		//"customresourcedefinitions",
		"felixconfigurations",
		"ippools",
		"nodes",
		"priorityclasses",
		"ciliumendpoints",
		"leases",
	}
}

func GetDefaultAppConfig() AppConfig {
	return AppConfig{
		OutputDir:                "test",
		AppendContextToOutputDir: true,
		ExcludedResourceTypes:    getDefaultExcludedResourceTypes(),
		IncludeSecrets:           false,
		EncryptSecrets:           true,
		SecretsEncryptAlgo:       "",
		SecretsEncryptKey:        "",
	}
}

func (config *AppConfig) GetOutDir(currentContext string) string {
	if config.AppendContextToOutputDir {
		return config.OutputDir + "/" + currentContext
	} else {
		return config.OutputDir
	}
}

func (config *AppConfig) IncludeResource(resourceTypeName string) bool {
	return !funk.ContainsString(config.ExcludedResourceTypes, resourceTypeName) &&
		(resourceTypeName != "secrets" || config.IncludeSecrets)
}
