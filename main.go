package main

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/gigurra/kdump/config"
	"github.com/gigurra/kdump/internal/k8s"
	"github.com/gigurra/kdump/internal/util/util_crypt"
	"github.com/gigurra/kdump/internal/util/util_file"
	"github.com/gigurra/kdump/internal/util/util_log"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
)

type Params struct {
	OutputDir            boa.Required[string] `descr:"The output directory to create"`
	DeletePrevDir        boa.Required[bool]   `descr:"Delete previous output directory" default:"false"`
	SecretsEncryptionKey boa.Optional[string] `descr:"Symmetric secrets encryption hex key for aes GCM (lower case 64 chars)"`
	GcpLogFormat         boa.Required[bool]   `descr:"Use GCP log format" default:"false"`
}

func main() {

	f := Params{}

	boa.Wrap{
		Use:    "kdump -o <output-dir> [-d] [-e <encryption-key>]",
		Short:  "Dump all kubernetes resources as yaml files to a dir",
		Params: &f,
		ParamEnrich: boa.ParamEnricherCombine(
			boa.ParamEnricherName,
			boa.ParamEnricherShort,
			// ParamEnricherEnv, // don't want this
			boa.ParamEnricherBool,
		),
		Run: func(cmd *cobra.Command, args []string) {
			appConfig := config.GetDefaultAppConfig()
			appConfig.OutputDir = f.OutputDir.Value()
			appConfig.DeletePrevDir = f.DeletePrevDir.Value()
			appConfig.SecretsEncryptKey = f.SecretsEncryptionKey.GetOrElse("")
			appConfig.GcpLogFormat = f.GcpLogFormat.Value()
			appConfig.Validate()
			dumpCurrentContext(appConfig)
		},
	}.ToApp()
}

func dumpCurrentContext(appConfig config.AppConfig) {

	if appConfig.GcpLogFormat {
		util_log.ConfigureGcpCompatibleJsonDefaultSlog(util_log.LogCfg{}.Default())
	}

	slog.Info("Running kdump version " + Version)

	slog.Info("Checking that we have at least one kubectl context...")
	if len(k8s.ListAvailableContexts()) == 0 {
		slog.Error("Bailing! No kubectl contexts available!")
		os.Exit(1)
	}

	slog.Info("Checking that we have a current kubectl context...")
	if k8s.CurrentContext() == "" {
		slog.Error("Bailing! No current kubectl context available!")
		os.Exit(1)
	}

	slog.Info("Checking output dir..")
	rootOutputDir := ensureRootOutputDir(appConfig)

	slog.Info("Downloading all resources from current context")

	apiResourceTypes := k8s.ApiResourceTypes()
	resourcesToDownload := appConfig.FilterIncludedResources(apiResourceTypes.Accessible.All)

	slog.Info(fmt.Sprintf("Downloading all resources of %d types", len(resourcesToDownload)))
	everythingRaw := k8s.DownloadEverything(resourcesToDownload)

	slog.Info(fmt.Sprintf("Running kubectl neat on everything.. (this takes about 3-4x the download time)"))
	everything := k8s.PipeToCommand(everythingRaw, "kubectl", "neat")

	slog.Info(fmt.Sprintf("Parsing %d bytes...\n", len(everything)))

	k8sResources := k8s.ParseResources(everything)
	k8sResourcesByNamespace := k8s.GroupByNamespace(k8sResources)

	slog.Info(fmt.Sprintf("Storing %d resources in '%s'...\n", len(k8sResources), rootOutputDir))
	for namespace, resources := range k8sResourcesByNamespace {
		outDir := rootOutputDir
		if namespace != "" {
			outDir = rootOutputDir + "/" + namespace
			util_file.CreateFolderIfNotExists(outDir, "could not create output dir: "+outDir)
		}
		for _, resource := range resources {
			name := util_file.SanitizePath(resource.MetaData.Name)
			typ := util_file.SanitizePath(resource.QualifiedTypeName)
			filename := name + "." + typ + ".yaml"
			if resource.IsSecret() {
				filePath := outDir + "/" + filename + ".aes"
				util_file.String2File(filePath, util_crypt.Encrypt(resource.SourceYaml, appConfig.SecretsEncryptKey))
			} else {
				filePath := outDir + "/" + filename
				util_file.String2File(filePath, resource.SourceYaml)
			}
		}
	}
}

func ensureRootOutputDir(appConfig config.AppConfig) string {

	out := appConfig.OutputDir

	if appConfig.DeletePrevDir {
		util_file.DeleteIfExists(out, fmt.Sprintf("removal of outputdir '%s' failed", out))
	}

	if util_file.Exists(out, fmt.Sprintf("checking outputdir '%s' failed", out)) {
		log.Fatal("Bailing! output-dir already exists: " + out)
	}

	util_file.CreateFolderIfNotExists(out, fmt.Sprintf("could not create folder '%s'", out))

	return out
}
