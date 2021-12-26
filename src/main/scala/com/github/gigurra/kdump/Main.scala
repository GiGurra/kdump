package com.github.gigurra.kdump

import com.github.gigurra.kdump.config.AppConfig
import internal.util
import internal.util.kubectl
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

  val ops = for namespace <- namespaces yield
    dumpNamespacedResources(outputDir, namespace, namespacedResourceTypeNames)

  ops.foreach(_.join)


def dumpNamespacedResources(outputDir: String,
                            namespace: String,
                            namespacedResourceTypeNames: List[String]): async.asyncOp[Unit] =

  async.run {
    util.file.mkDirs(s"$outputDir/$namespace")
    for resourceTypeName <- namespacedResourceTypeNames do
      dumpNamespacedResources(outputDir, namespace, resourceTypeName)
  }


def dumpNamespacedResources(outputDir: String,
                            namespace: String,
                            resourceTypeName: String): Unit =
  val resourceNames = kubectl.listNamespacedResourcesOfType(namespace, resourceTypeName)
  if resourceNames.nonEmpty then
    if resourceTypeName == "secrets" then
      dumpSecrets(outputDir, namespace, resourceTypeName, resourceNames)
    else
      dumpRegularNamespacedResources(outputDir, namespace, resourceTypeName, resourceNames)

def dumpRegularNamespacedResources(outputDir: String,
                                   namespace: String,
                                   resourceTypeName: String,
                                   resourceNames: List[String]): Unit =
  val itemDir = s"$outputDir/$namespace/$resourceTypeName"
  util.file.mkDirs(itemDir)
  for item <- resourceNames do
    println("hej")

/*
for _
, item := range resourceNames {
resource := kubectl.DownloadNamespacedResource(namespace, resourceTypeName, item, "yaml")
log.Printf("Storing item %v in folder %v", item, itemDir)
fileutil.String2File(itemDir + "/" + fileutil.ReplaceInvalidChars(item) + ".yaml", resource)*/
//}


def dumpSecrets(outputDir: String,
                namespace: String,
                resourceTypeName: String,
                resourceNames: List[String]): Unit =
  println("ignoring storage of secrets. Not yet implemented")
