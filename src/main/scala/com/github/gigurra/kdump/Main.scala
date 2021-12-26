package com.github.gigurra.kdump

import com.github.gigurra.kdump.config.AppConfig
import internal.util.*

import scala.concurrent.{ExecutionContext, Future}

@main def kdump(): Unit =
  val appConfig = AppConfig.default
  dumpCurrentContext(appConfig)

def dumpCurrentContext(appConfig: AppConfig): Unit =

  given ec: ExecutionContext = scala.concurrent.ExecutionContext.global

  val currentK8sContext = kubectl.currentContext()
  val outputDir = appConfig.outDir(currentK8sContext)

  file.delete(outputDir)
  file.mkDirs(outputDir)

  println(s"Downloading all resources from current context to dir '$outputDir'")

  val namespacesOp = async.run(kubectl.namespaces())
  val allResourceTypeNamesOp = async.run(kubectl.resourceTypeNames().filter(appConfig.isResourceIncluded))
  val globalResourceTypeNamesOp = async.run(kubectl.globalResourceTypeNames().filter(appConfig.isResourceIncluded))
  val namespacedResourceTypeNamesOp = async.run(kubectl.namespacedResourceTypeNames().filter(appConfig.isResourceIncluded))

  lazy val namespaces = async.await(namespacesOp)
  lazy val allResourceTypeNames = async.await(allResourceTypeNamesOp)
  lazy val globalResourceTypeNames = async.await(globalResourceTypeNamesOp)
  lazy val namespacedResourceTypeNames = async.await(namespacedResourceTypeNamesOp)

  println(namespaces)
  println()

  println(s"global resource type names:")
  globalResourceTypeNames.foreach(println)
  println()

  println(s"namespaced resource type names:")
  namespacedResourceTypeNames.foreach(println)
  println()



