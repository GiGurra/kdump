package k8s

import (
	"github.com/samber/lo"
	"gopkg.in/yaml.v2"
	"strings"
)

type ApiVersion struct {
	Version string
	Name    string
}

type ApiResourceType struct {
	Name          string
	ShortNames    []string
	Namespaced    bool
	Kind          string
	Verbs         []string
	ApiVersion    ApiVersion
	QualifiedName string
}

func (r *ApiResourceType) IsSecret() bool {
	return strings.ToLower(r.Name) == "secret" || strings.ToLower(r.Name) == "secrets"
}

type Resource struct {
	Kind              string           `yaml:"kind"`
	ApiVersion        string           `yaml:"apiVersion"`
	MetaData          ResourceMetadata `yaml:"metadata"`
	SourceYaml        string
	QualifiedTypeName string
}

func (r *Resource) IsNamespaced() bool {
	return len(r.MetaData.Namespace) > 0
}

func (r *Resource) IsGlobal() bool {
	return !r.IsNamespaced()
}

func (r *Resource) IsSecret() bool {
	return strings.ToLower(r.Kind) == "secret" || strings.ToLower(r.Kind) == "secrets"
}

func ParseResources(in string) []*Resource {

	root := make(map[string]interface{})

	err := yaml.Unmarshal([]byte(in), &root)

	if err != nil {
		panic("Could not parse input yaml due to " + err.Error())
	}

	if _, ok := root["items"]; ok {
		return ParseResourceList(root)
	} else {
		return []*Resource{ParseSingleResource(root, 0)}
	}
}

func GroupByNamespace(all []*Resource) map[string][]*Resource {
	return lo.GroupBy(all, func(r *Resource) string {
		return r.MetaData.Namespace
	})
}

type ResourceMetadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

func ParseResourceList(in map[string]interface{}) []*Resource {
	items := in["items"].([]interface{})
	return lo.Map(items, ParseSingleResource)
}

func ParseSingleResource(item interface{}, _ int) *Resource {

	m := item.(map[interface{}]interface{})

	delete(m, "lastRefresh")
	delete(m, "status")

	yamlString, marshallError := yaml.Marshal(m)
	if marshallError != nil {
		panic("Failed marshalling object to yaml, due to: " + marshallError.Error())
	}
	var out Resource
	if yaml.Unmarshal(yamlString, &out) != nil {
		panic("Failed reading back yaml: " + string(yamlString))
	}

	out.SourceYaml = string(yamlString)

	parts := strings.Split(out.ApiVersion, "/")
	if len(parts) == 1 {
		out.QualifiedTypeName = strings.ToLower(out.Kind)
	} else if len(parts) == 2 {
		out.QualifiedTypeName = strings.ToLower(out.Kind) + "." + parts[0]
	} else {
		panic("Unable to parse QualifiedTypeName from string: " + out.ApiVersion)
	}

	return &out
}
