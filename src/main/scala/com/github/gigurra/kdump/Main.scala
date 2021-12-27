package com.github.gigurra.kdump

import com.github.gigurra.kdump.config.AppConfig
import internal.util
import internal.util.kubectl
import internal.util.async
import internal.util.k8sYaml

import scala.concurrent.{ExecutionContext, Future}

@main def kdump(): Unit =
  val appConfig = AppConfig.default
  dumpCurrentContext(appConfig)

def dumpCurrentContext(appConfig: AppConfig): Unit =

  val currentK8sContext = kubectl.currentContext()
  val outputDir = appConfig.outDir(currentK8sContext)

  util.file.delete(outputDir)
  util.file.mkDirs(outputDir)

  println(s"Downloading all resources from current context to dir '$outputDir'")

  // Do all of these in parallel
  val namespacesOp = async.run(kubectl.namespaces())
  val allResourceTypeNamesOp = async.run(kubectl.resourceTypeNames().filter(appConfig.isResourceIncluded))
  val globalResourceTypeNamesOp = async.run(kubectl.globalResourceTypeNames().filter(appConfig.isResourceIncluded))
  val namespacedResourceTypeNamesOp = async.run(kubectl.namespacedResourceTypeNames().filter(appConfig.isResourceIncluded))

  // Gather the results of parallel computations
  lazy val namespaces = namespacesOp.join
  lazy val allResourceTypeNames = allResourceTypeNamesOp.join
  lazy val globalResourceTypeNames = globalResourceTypeNamesOp.join
  lazy val namespacedResourceTypeNames = namespacedResourceTypeNamesOp.join

  // We could do these in parallel as well, but doing it brute force in parallel overloads kubectl and the k8s api server :P
  for namespace <- namespaces do
    dumpNamespacedResources(outputDir, namespace, namespacedResourceTypeNames)


def dumpNamespacedResources(outputDir: String,
                            namespace: String,
                            namespacedResourceTypeNames: List[String]): Unit =

  println(s"processing namespace $namespace")

  val namespaceDir = s"$outputDir/$namespace"

  util.file.mkDirs(namespaceDir)

  val allYaml = kubectl.downloadAllResources(namespace, namespacedResourceTypeNames, "yaml")
  val parsedYaml: Seq[k8sYaml.K8sResource] = k8sYaml.parseList(allYaml)
  val resourcesByType: Map[String, Seq[k8sYaml.K8sResource]] = parsedYaml.groupBy(_.qualifiedKind)

  for (kind, resources) <- resourcesByType do
    val resourceDir = s"$namespaceDir/${util.file.sanitizeFileName(kind)}"
    util.file.mkDirs(resourceDir)
    for resource <- resources do
      if resource.isSecret then
        println(s"ignoring secret, not yet implemented: $resource")
      else
        util.file.string2File(s"$resourceDir/${util.file.sanitizeFileName(resource.name)}.yaml", resource.source)
