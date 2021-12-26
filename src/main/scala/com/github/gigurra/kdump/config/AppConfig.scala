package com.github.gigurra.kdump.config


case class AppConfig(outputDir: String,
                     appendContextToOutputDir: Boolean,
                     excludedResourceTypes: Set[String],
                     secretsHandling: AppConfig.SecretsHandling)

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