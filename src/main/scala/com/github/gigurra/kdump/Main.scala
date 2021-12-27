package com.github.gigurra.kdump

import com.github.gigurra.kdump.config.AppConfig
import internal.util
import internal.util.kubectl
import internal.util.kubectl.K8sResource
import internal.util.async

import scala.concurrent.{ExecutionContext, Future}

@main def kdump(): Unit =
  val appConfig = AppConfig.default
  dumpCurrentContext(appConfig)

def dumpCurrentContext(appConfig: AppConfig): Unit =

  val currentK8sContext = kubectl.currentContext()
  val outputDir = appConfig.outDir(currentK8sContext)

  util.file.delete(outputDir)
  util.file.mkDirs(outputDir)

  println(s"Downloading all resources from current context '$currentK8sContext' to dir '$outputDir'")

  // Do all of these in parallel
  val namespaces = async.run(kubectl.namespaces())
  val allResourceTypeNames = async.run(kubectl.resourceTypeNames().filter(appConfig.isResourceIncluded))
  val globalResourceTypeNames = async.run(kubectl.globalResourceTypeNames().filter(appConfig.isResourceIncluded))
  val namespacedResourceTypeNames = async.run(kubectl.namespacedResourceTypeNames().filter(appConfig.isResourceIncluded))

  for namespace <- namespaces do
    dumpNamespacedResources(outputDir, namespace, namespacedResourceTypeNames)


def dumpNamespacedResources(outputDir: String,
                            namespace: String,
                            namespacedResourceTypeNames: List[String]): Unit =

  println(s"processing namespace $namespace")

  val namespaceDir = s"$outputDir/$namespace"

  util.file.mkDirs(namespaceDir)

  val parsedYaml = kubectl.downloadAllNamespacedResources(namespace, namespacedResourceTypeNames)

  for (kind, resources) <- parsedYaml.groupBy(_.qualifiedKind) do
    val resourceDir = s"$namespaceDir/${util.file.sanitizeFileName(kind)}"
    util.file.mkDirs(resourceDir)
    for resource <- resources do
      if resource.isSecret then
        println(s"ignoring secret, not yet implemented: $resource")
      else
        util.file.string2File(s"$resourceDir/${util.file.sanitizeFileName(resource.name)}.yaml", resource.sourceYaml)
