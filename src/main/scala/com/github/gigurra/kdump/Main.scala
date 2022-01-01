package com.github.gigurra.kdump

import com.github.gigurra.kdump.config.AppConfig
import internal.util
import internal.util.kubectl
import internal.util.kubectl.K8sResource
import internal.util.async
import org.slf4j.LoggerFactory

import scala.concurrent.{ExecutionContext, Future}

private val log = LoggerFactory.getLogger("kdump")

@main def kdump(): Unit =
  val appConfig = AppConfig.default
  dumpCurrentContext(appConfig)

def dumpCurrentContext(appConfig: AppConfig): Unit =

  val outputDir = appConfig.outputDir

  util.file.delete(outputDir)
  util.file.mkDirs(outputDir)

  log.info(s"Downloading all resources from current context to dir '$outputDir'")

  // Do all of these in parallel
  val allResourceTypeNames = async(kubectl.resourceTypeNames().filter(appConfig.isResourceIncluded))
  val everything: Seq[K8sResource] = kubectl.downloadAllResources(allResourceTypeNames)

  for (namespaceOpt, resources) <- everything.groupBy(_.namespace) do
    dumpNamespacedResources(outputDir, namespaceOpt, resources)

  log.info(s"DONE!")

def dumpNamespacedResources(rootOutputDir: String,
                            namespace: Option[String],
                            resources: Seq[K8sResource]): Unit =

  log.info(s"processing ${namespace.map(n => s"namespace $n").getOrElse("global resources")}")

  val outputDir = namespace.fold(rootOutputDir)(n => s"$rootOutputDir/$n")

  util.file.mkDirs(outputDir)

  for (kind, resources) <- resources.groupBy(_.qualifiedKind) do
    for resource <- resources do
      val filepath = s"$outputDir/${util.file.sanitizeFileName(resource.qualifiedName)}.yaml"
      if resource.isSecret then
        println(s"ignoring secret, not yet implemented: $resource")
      else
        util.file.string2File(filepath, resource.sourceYaml)
