package com.github.gigurra.kdump.config

import com.github.gigurra.kdump.config.AppConfig.SecretsHandling


case class AppConfig(outputDir: String,
                     appendContextToOutputDir: Boolean,
                     excludedResourceTypes: Set[String],
                     secretsHandling: AppConfig.SecretsHandling) {

  def outDir(currentContext: String): String =
    if appendContextToOutputDir then
      outputDir + "/" + currentContext
    else
      outputDir

  def isResourceIncluded(resourceTypeName: String): Boolean =
    !excludedResourceTypes.contains(resourceTypeName) &&
      (resourceTypeName != "secrets" || secretsHandling != SecretsHandling.DontStore)
}

object AppConfig {

  def default: AppConfig = AppConfig(
    outputDir = "test",
    appendContextToOutputDir = true,
    excludedResourceTypes = defaultExcludedResourceTypes,
    secretsHandling = SecretsHandling.DontStore,
  )

  def defaultExcludedResourceTypes: Set[String] = Set(
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
    "felixconfigurations",
    "ippools",
    "nodes",
    "priorityclasses",
    "ciliumendpoints",
    "leases",
    "certificaterequests",
    "orders",
  )

  sealed trait SecretsHandling

  object SecretsHandling {
    object DontStore extends SecretsHandling

    object StoreAsPlainText extends SecretsHandling

    case class StoreEncrypted(encryption: String /*todo: impl*/) extends SecretsHandling
  }
}